package cache

import (
	"bytes"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/profile"
	"github.com/behavioral-ai/intermediary/request"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	NamespaceName  = "resiliency:agent/behavioral-ai/intermediary/cache"
	Route          = "cache"
	defaultTimeout = time.Millisecond * 3000
)

var (
	noContentResponse   = httpx.NewResponse(http.StatusNoContent, nil, nil)
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
	maxDuration         = time.Minute * 30
)

type agentT struct {
	running  bool
	enabled  *atomic.Bool
	hostName string
	timeout  time.Duration
	profile  profile.Cache

	exchange httpx.Exchange
	ticker   *messaging.Ticker
	emissary *messaging.Channel
	handler  eventing.Agent
}

// New - create a new cache agent
func init() {
	a := newAgent(eventing.Handler)
	exchange.Register(a)
}

func newAgent(handler eventing.Agent) *agentT {
	a := new(agentT)
	a.enabled = new(atomic.Bool)
	a.enabled.Store(true)
	a.timeout = defaultTimeout

	a.exchange = httpx.Do
	a.ticker = messaging.NewTicker(messaging.ChannelEmissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
	a.handler = handler
	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if !a.running {
		if m.Event() == messaging.ConfigEvent {
			a.configure(m)
			return
		}
		if m.Event() == messaging.StartupEvent {
			a.run()
			a.running = true
			return
		}
		return
	}
	if m.Event() == messaging.ShutdownEvent {
		a.running = false
	}
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) run() {
	go emissaryAttend(a)
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool                { return true }
func (a *agentT) Route() string            { return Route }
func (a *agentT) Timeout() time.Duration   { return a.timeout }
func (a *agentT) Exchange() httpx.Exchange { return a.exchange }

// Link - chainable exchange
func (a *agentT) Link(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		var (
			cacheable = a.cacheable(r)
			url       string
			status    *messaging.Status
		)
		if cacheable {
			url = uri.BuildURL(a.hostName, r.URL.Path, r.URL.Query())
			h := make(http.Header)
			h.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
			resp, status = request.Do(a, http.MethodGet, url, h, nil)
			if resp.StatusCode == http.StatusOK {
				// Need for analytics
				//resp.Header.Add(access.XCached, "true")
				return resp, nil
			}
			if status.Err != nil {
				a.handler.Notify(status.WithAgent(a.Uri()))
			}
		}
		if next == nil {
			return noContentResponse, nil
		}
		resp, err = next(r)
		if cacheable && resp.StatusCode == http.StatusOK {
			var buf []byte
			buf, err = io.ReadAll(resp.Body)
			if err != nil {
				status = messaging.NewStatusError(messaging.StatusIOError, err, a.Uri())
				a.handler.Notify(status)
				return serverErrorResponse, err
			}
			resp.ContentLength = int64(len(buf))
			resp.Body = io.NopCloser(bytes.NewReader(buf))
			go func() {
				h := httpx.CloneHeader(resp.Header)
				h.Add(httpx.XRequestId, r.Header.Get(httpx.XRequestId))
				_, status = request.Do(a, http.MethodPut, url, h, io.NopCloser(bytes.NewReader(buf)))
				if status.Err != nil {
					a.handler.Notify(status.WithAgent(a.Uri()))
				}
			}()
		}
		return
	}
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeMap:
		var ok bool
		if a.hostName, ok = config.CacheHostName(a, m); !ok {
			return
		}
	case httpx.ContentTypeExchange:
		if ex, ok := httpx.ConfigExchangeContent(m); ok {
			a.exchange = ex
		}
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}

func (a *agentT) cacheable(r *http.Request) bool {
	if a.hostName == "" || r.Method != http.MethodGet || httpx.CacheControlNoCache(r.Header) {
		return false
	}
	return a.enabled.Load()
}

func (a *agentT) emissaryShutdown() {
	a.emissary.Close()
	a.ticker.Stop()
}
