package routing

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/eventing"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/request"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/routing/request/http"
	Route         = "app"
	DefaultRoute  = "core:routing/default"
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
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(eventing.Handler)
	})
}

func newAgent(handler eventing.Agent) *agentT {
	a := new(agentT)
	a.log = true
	a.defaultRoute.Name = DefaultRoute
	a.defaultRoute.Ex = httpx.Do
	a.router = rest.NewRouter()
	a.routerModify("", httpx.Do)

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
	if m.Name() == messaging.ConfigEvent {
		a.configure(m)
	}
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool              { return a.log }
func (a *agentT) Route() string          { return Route }
func (a *agentT) Timeout() time.Duration { return a.timeout }
func (a *agentT) Do() rest.Exchange {
	if rt, ok := a.router.Lookup(DefaultRoute); ok {
		return rt.Ex
	}
	return httpx.Do
}

// Exchange - implementation for rest.Exchangeable interface
func (a *agentT) Exchange(r *http.Request) (resp *http.Response, err error) {
	rt := a.routerLookup()
	if rt != nil && rt.Uri == "" {
		status := messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("host configuration is empty")).WithLocation(a.Name())
		a.handler.Notify(status)
		return serverErrorResponse, status.Err
	}
	var status *messaging.Status

	url := uri.BuildURL(rt.Uri, r.URL.Path, r.URL.Query())
	// TODO : need to check and remove Caching header.
	resp, status = request.Do(a, r.Method, url, httpx.CloneHeaderWithEncoding(r), r.Body)
	if status.Err != nil {
		a.handler.Notify(status.WithLocation(a.Name()))
	}
	if resp.StatusCode == http.StatusGatewayTimeout {
		resp.Header.Add(access2.XTimeout, fmt.Sprintf("%v", a.timeout))
	}
	return resp, status.Err
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case httpx.ContentTypeExchange:
		if ex, ok := httpx.ConfigExchangeContent(m); ok {
			a.routerModify("", ex)
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
		a.routerModify(hostName, nil)
		if a.timeout, ok = config.Timeout(a, m); !ok {
			return
		}
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}

func (a *agentT) routerModify(uri string, ex rest.Exchange) {
	a.router.Modify(DefaultRoute, uri, ex)
}

func (a *agentT) routerLookup() (r *rest.Route) {
	r, _ = a.router.Lookup(DefaultRoute)
	return
}
