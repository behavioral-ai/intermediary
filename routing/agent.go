package routing

import (
	"errors"
	"fmt"
	center "github.com/behavioral-ai/center/messaging"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/access2"
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
	comms  *center.Communication

	review *messaging.Review
}

// init - register an agent constructor
func init() {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(representation1.Initialize(nil), nil, center.Comms)
	})
}

func ConstructorOverride(m map[string]string, ex rest.Exchange, comms *center.Communication) {
	repository.RegisterConstructor(NamespaceName, func() messaging.Agent {
		return newAgent(representation1.Initialize(m), ex, comms)
	})
}

func newAgent(state *representation1.Routing, ex rest.Exchange, comms *center.Communication) *agentT {
	a := new(agentT)
	a.state = state
	a.comms = comms
	if ex == nil {
		ex = httpx.Do
	}
	a.router = rest.NewRouter()
	a.router.Modify(defaultRoute, a.state.AppHost, ex)
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
	if m.Name == messaging.ConfigEvent {
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
		a.comms.Notify(messaging.NewStatusMessage(status, a.Name()))
		return serverErrorResponse, status.Err
	}
	var status *messaging.Status

	url := uri.BuildURL(rt.Uri, r.URL.Path, r.URL.Query())
	// TODO : need to check and remove Caching header.
	resp, status = request.Do(a, r.Method, url, httpx.CloneHeaderWithEncoding(r), r.Body)
	if status.Err != nil {
		a.comms.Notify(messaging.NewStatusMessage(status.WithLocation(a.Name()), a.Name()))
	}
	if resp.StatusCode == http.StatusGatewayTimeout {
		resp.Header.Add(access2.XTimeout, fmt.Sprintf("%v", a.state.Timeout))
	}
	return resp, status.Err
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
	a.comms.Trace(a.Name(), task, observation, action)
}

func (a *agentT) configure(m *messaging.Message) {
	switch m.ContentType() {
	case messaging.ContentTypeMap:
		cfg := messaging.ConfigMapContent(m)
		if cfg == nil {
			messaging.Reply(m, messaging.ConfigEmptyMapError(a.Name()), a.Name())
			return
		}
		a.state.Update(cfg)
		a.router.Modify(defaultRoute, a.state.AppHost, nil)
	case messaging.ContentTypeReview:
		r := messaging.ReviewContent(m)
		if r == nil {
			messaging.Reply(m, messaging.ConfigEmptyReviewError(a.Name()), a.Name())
			return
		}
		a.review = r
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
