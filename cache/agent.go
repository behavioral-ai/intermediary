package cache

import (
	"bytes"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/eventing"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/intermediary/cache/representation1"
	"github.com/behavioral-ai/intermediary/request"
	"io"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/cache/request/http"
	Route         = "cache"
)

var (
	noContentResponse   = httpx.NewResponse(http.StatusNoContent, nil, nil)
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	state    *representation1.Cache
	exchange rest.Exchange

	ticker   *messaging.Ticker
	emissary *messaging.Channel
	handler  eventing.Agent
}

// init - register an agent constructor
func init() {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(eventing.Handler, representation1.NewCache(NamespaceName), nil)
	})
}

func ConstructorOverride(m map[string]string, ex rest.Exchange) {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		c := representation1.Initialize()
		c.Update(m)
		return newAgent(eventing.Handler, c, ex)
	})
}

func newAgent(handler eventing.Agent, state *representation1.Cache, ex rest.Exchange) *agentT {
	a := new(agentT)
	if state == nil {
		a.state = representation1.Initialize()
	} else {
		a.state = state
	}
	if ex == nil {
		a.exchange = httpx.Do
	} else {
		a.exchange = ex
	}
	a.ticker = messaging.NewTicker(messaging.ChannelEmissary, a.state.Interval)
	a.emissary = messaging.NewEmissaryChannel()
	a.handler = handler
	return a
}

// String - identity
func (a *agentT) String() string { return a.Name() }

// Name - agent identifier
func (a *agentT) Name() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if !a.state.Running {
		if m.Name() == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Name() == messaging.StartupEvent {
			a.run()
			a.state.Running = true
			return
		}
		return
	}
	if m.Name() == messaging.ShutdownEvent {
		a.state.Running = false
	}
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) run() {
	go emissaryAttend(a)
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool              { return true }
func (a *agentT) Route() string          { return Route }
func (a *agentT) Timeout() time.Duration { return a.state.Timeout }
func (a *agentT) Do() rest.Exchange      { return a.exchange }

// Link - chainable exchange
func (a *agentT) Link(next rest.Exchange) rest.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if !a.cacheable(r) {
			return next(r)
		}
		var (
			url    string
			status *messaging.Status
		)
		// cache lookup
		url = uri.BuildURL(a.state.Host, r.URL.Path, r.URL.Query())
		h := make(http.Header)
		h.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
		resp, status = request.Do(a, http.MethodGet, url, h, nil)
		if resp.StatusCode == http.StatusOK {
			resp.Header.Add(access2.XCached, "true")
			return resp, nil
		}
		resp.Header.Add(access2.XCached, "false")
		if status.Err != nil {
			a.handler.Notify(status.WithLocation(a.Name()))
		}
		// cache miss, call next exchange
		resp, err = next(r)
		if resp.StatusCode == http.StatusOK {
			// cache update
			err = a.cacheUpdate(url, r, resp)
			if err != nil {
				return serverErrorResponse, err
			}
		}
		return
	}
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeMap:
		cfg := messaging.ConfigMapContent(m)
		if cfg == nil {
			messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Name())
			return
		}
		a.state.Update(cfg)
	}
	/*
		case httpx.ContentTypeExchange:
			if ex, ok := httpx.ConfigExchangeContent(m); ok {
				a.exchange = ex
			}
		}
	*/

	messaging.Reply(m, messaging.StatusOK(), a.Name())
}

func (a *agentT) cacheable(r *http.Request) bool {
	if a.state.Host == "" || r.Method != http.MethodGet || httpx.CacheControlNoCache(r.Header) {
		return false
	}
	return a.state.Enabled.Load()
}

func (a *agentT) emissaryShutdown() {
	a.emissary.Close()
	a.ticker.Stop()
}

func (a *agentT) cacheUpdate(url string, r *http.Request, resp *http.Response) error {
	var (
		buf    []byte
		err    error
		status *messaging.Status
	)
	// TODO: Need to reset the body in the response after reading it.
	buf, err = io.ReadAll(resp.Body)
	if err != nil {
		status = messaging.NewStatus(messaging.StatusIOError, err).WithLocation(a.Name())
		a.handler.Notify(status)
		return err
	}
	resp.ContentLength = int64(len(buf))
	resp.Body = io.NopCloser(bytes.NewReader(buf))

	// cache update
	go func() {
		h2 := httpx.CloneHeader(resp.Header)
		h2.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
		_, status = request.Do(a, http.MethodPut, url, h2, io.NopCloser(bytes.NewReader(buf)))
		if status.Err != nil {
			a.handler.Notify(status.WithLocation(a.Name()))
		}
	}()
	return nil
}
