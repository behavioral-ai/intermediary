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
	"github.com/behavioral-ai/intermediary/request"
	"github.com/behavioral-ai/intermediary/routing/representation1"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/routing/request/http"
	defaultRoute  = "core:routing/default"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	state  *representation1.Routing
	router *rest.Router

	handler eventing.Agent
}

// init - register an agent constructor
func init() {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(eventing.Handler, representation1.NewRouting(NamespaceName), nil)
	})
}

func ConstructorOverride(m map[string]string, ex rest.Exchange) {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		c := representation1.Initialize()
		c.Update(m)
		return newAgent(eventing.Handler, c, ex)
	})
}

func newAgent(handler eventing.Agent, state *representation1.Routing, ex rest.Exchange) *agentT {
	a := new(agentT)
	if state == nil {
		a.state = representation1.Initialize()
	} else {
		a.state = state
	}
	if ex == nil {
		ex = httpx.Do
	}
	a.router = rest.NewRouter()
	a.router.Modify(defaultRoute, a.state.AppHost, ex)
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
func (a *agentT) Log() bool              { return a.state.Log }
func (a *agentT) Route() string          { return a.state.LogRouteName }
func (a *agentT) Timeout() time.Duration { return a.state.Timeout }
func (a *agentT) Do() rest.Exchange {
	if rt, ok := a.router.Lookup(defaultRoute); ok {
		return rt.Ex
	}
	return httpx.Do
}

// Exchange - implementation for rest.Exchangeable interface
func (a *agentT) Exchange(r *http.Request) (resp *http.Response, err error) {
	rt, ok := a.router.Lookup(defaultRoute)
	if !ok || rt != nil && rt.Uri == "" {
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
		resp.Header.Add(access2.XTimeout, fmt.Sprintf("%v", a.state.Timeout))
	}
	return resp, status.Err
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
		a.router.Modify(defaultRoute, a.state.AppHost, nil)
	}
	messaging.Reply(m, messaging.StatusOK(), a.Name())
}

/*
func (a *agentT) routerModify(uri string, ex rest.Exchange) {
	a.router.Modify(defaultRoute, uri, ex)
}

func (a *agentT) routerLookup() (r *rest.Route) {
	r, _ = a.router.Lookup(defaultRoute)
	return
}


*/
