package cache

import (
	"bytes"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/access2"
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
	NamespaceName = "test:resiliency:agent/cache/request/http"
	Route         = "cache"
)

var (
	noContentResponse   = httpx.NewResponse(http.StatusNoContent, nil, nil)
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	state    *representation1.Cache
	exchange rest.Exchange
	service  *operations.Service

	review   *messaging.Review
	ticker   *messaging.Ticker
	emissary *messaging.Channel
}

// init - register an agent constructor
func init() {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(representation1.Initialize(nil), nil, operations.Serve)
	})
}

func ConstructorOverride(m map[string]string, ex rest.Exchange, service *operations.Service) {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(representation1.Initialize(m), ex, service)
	})
}

func newAgent(state *representation1.Cache, ex rest.Exchange, service *operations.Service) *agentT {
	a := new(agentT)
	a.state = state
	a.service = service
	if ex == nil {
		a.exchange = httpx.Do
	} else {
		a.exchange = ex
	}
	a.ticker = messaging.NewTicker(messaging.ChannelEmissary, a.state.Interval)
	a.emissary = messaging.NewEmissaryChannel()
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
		if m.Name == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Name == messaging.StartupEvent {
			a.run()
			a.state.Running = true
			return
		}
		return
	}
	if m.Name == messaging.ShutdownEvent {
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
			a.service.Message(messaging.NewStatusMessage(status.WithLocation(a.Name()), a.Name()))
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

func (a *agentT) trace(task, observation, action string) {
	if a.review == nil {
		return
	}
	if !a.review.Started() {
		a.review.Start()
	}
	if a.review.Expired() {
		return
	}
	a.service.Trace(a.Name(), task, observation, action)
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeMap:
		cfg, status := messaging.MapContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		a.state.Update(cfg)
	case messaging.ContentTypeReview:
		r, status := messaging.ReviewContent(m)
		if !status.OK() {
			messaging.Reply(m, status, a.Name())
			return
		}
		a.review = r
	}
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
		a.service.Message(messaging.NewStatusMessage(status, a.Name()))
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
			a.service.Message(messaging.NewStatusMessage(status.WithLocation(a.Name()), a.Name()))
		}
	}()
	return nil
}
