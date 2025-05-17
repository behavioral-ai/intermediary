package request

import (
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"io"
	"net/http"
	"time"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type Requester interface {
	Route() string
	Log() bool
	Timeout() time.Duration
	Do() rest.Exchange
}

func Do(agent Requester, method string, url string, h http.Header, r io.ReadCloser) (resp *http.Response, status *messaging.Status) {
	start := time.Now().UTC()
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return serverErrorResponse, messaging.NewStatus(messaging.StatusInvalidArgument, err)
	}
	req.Header = h
	resp, err = httpx.ExchangeWithTimeout(agent.Timeout(), agent.Do())(req)
	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	if err != nil {
		status = messaging.NewStatus(resp.StatusCode, err)
		return
	}
	status = messaging.StatusOK()
	if agent.Log() {
		access.Log(access.EgressTraffic, start, time.Since(start), agent.Route(), req, resp, access.Threshold{Timeout: agent.Timeout()})
	}
	return
}
