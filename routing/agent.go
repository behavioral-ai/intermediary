package routing

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/request"
	"github.com/behavioral-ai/intermediary/urn"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/intermediary/routing"
	Route         = "app"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	log          bool
	timeout      time.Duration
	router       *rest.Router
	defaultRoute rest.Route

	handler eventing.Agent
}

// New - create a new cache agent
func init() {
	a := newAgent(eventing.Handler)
	exchange.Register(a)
}

func newAgent(handler eventing.Agent) *agentT {
	a := new(agentT)
	a.log = true
	a.defaultRoute.Name = urn.DefaultRoute
	a.defaultRoute.Ex = httpx.Do
	a.router = rest.NewRouter()
	a.router.Modify(urn.DefaultRoute, "", httpx.Do)

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
	if m.Event() == messaging.ConfigEvent {
		a.configure(m)
	}
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool              { return a.log }
func (a *agentT) Route() string          { return Route }
func (a *agentT) Timeout() time.Duration { return a.timeout }
func (a *agentT) Do() rest.Exchange {
	if rt, ok := a.router.Lookup(urn.DefaultRoute); ok {
		return rt.Ex
	}
	return httpx.Do
}

// Exchange - implementation for rest.Exchangeable interface
func (a *agentT) Exchange(r *http.Request) (resp *http.Response, err error) {
	rt, ok := a.router.Lookup(urn.DefaultRoute)
	if ok && rt.Uri == "" {
		status := messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("host configuration is empty"), a.Uri())
		a.handler.Notify(status)
		return serverErrorResponse, status.Err
	}
	var status *messaging.Status

	url := uri.BuildURL(rt.Uri, r.URL.Path, r.URL.Query())
	// TODO : need to check and remove Caching header.
	resp, status = request.Do(a, r.Method, url, httpx.CloneHeaderWithEncoding(r), r.Body)
	if status.Err != nil {
		a.handler.Notify(status.WithAgent(a.Uri()))
	}
	if resp.StatusCode == http.StatusGatewayTimeout {
		if resp.Header == nil {
			resp.Header = make(http.Header)
		}
		resp.Header.Add(access.XTimeout, fmt.Sprintf("%v", a.timeout))
	}
	return resp, status.Err
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case httpx.ContentTypeExchange:
		if ex, ok := httpx.ConfigExchangeContent(m); ok {
			if rt, ok1 := a.router.Lookup(urn.DefaultRoute); ok1 {
				a.router.Modify(rt.Name, "", ex)
			}
		}
	case messaging.ContentTypeMap:
		var (
			ok       bool
			hostName string
		)
		hostName, ok = config.AppHostName(a, m)
		if !ok {
			return
		}
		if rt, ok1 := a.router.Lookup(urn.DefaultRoute); ok1 {
			a.router.Modify(rt.Name, hostName, nil)
			return
		}
		if a.timeout, ok = config.Timeout(a, m); !ok {
			return
		}
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
