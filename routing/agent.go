package routing

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/uri"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/request"
	"net/http"
	"time"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/intermediary/routing"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type agentT struct {
	log      bool
	hostName string
	timeout  time.Duration

	exchange httpx.Exchange
	handler  messaging.Agent
}

// New - create a new cache agent
func New(handler messaging.Agent) messaging.Agent {
	return newAgent(handler)
}

func newAgent(handler messaging.Agent) *agentT {
	a := new(agentT)
	a.log = true
	a.exchange = httpx.Do
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
		return
	}
}

func (a *agentT) configure(m *messaging.Message) {
	var (
		ok bool
		ex httpx.Exchange
	)

	if ex, ok = httpx.ConfigExchangeContent(m); ok {
		a.exchange = ex
	}
	if a.hostName, ok = config.AppHostName(a, m); !ok {
		return
	}
	if a.timeout, ok = config.Timeout(a, m); !ok {
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}

// Log - implementation for Requester interface
func (a *agentT) Log() bool                { return a.log }
func (a *agentT) Timeout() time.Duration   { return a.timeout }
func (a *agentT) Exchange() httpx.Exchange { return a.exchange }

// Link - implementation for httpx.Chainable interface
func (a *agentT) Link(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		if a.hostName == "" {
			status := messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("host configuration is empty"), a.Uri())
			a.handler.Message(eventing.NewNotifyMessage(status))
			return serverErrorResponse, status.Err
		}
		var status *messaging.Status

		url := uri.BuildURL(a.hostName, r.URL.Path, r.URL.Query())
		resp, status = request.Do(a, r.Method, url, httpx.CloneHeaderWithEncoding(r), r.Body)
		if status.Err != nil {
			a.handler.Message(eventing.NewNotifyMessage(status.WithAgent(a.Uri())))
		}
		if resp.StatusCode == http.StatusGatewayTimeout {
			resp.Header.Add(access.XTimeout, fmt.Sprintf("%v", a.timeout))
		}
		return resp, status.Err
	}
}
