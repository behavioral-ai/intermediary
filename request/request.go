package request

import (
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"io"
	"net/http"
	"time"
)

var (
	serverErrorResponse = httpx.NewResponse(http.StatusInternalServerError, nil, nil)
)

type Requester interface {
	Log() bool
	Timeout() time.Duration
	Exchange() httpx.Exchange
}

func Do(agent Requester, method string, url string, h http.Header, r io.ReadCloser) (resp *http.Response, status *messaging.Status) {
	start := time.Now().UTC()
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return serverErrorResponse, messaging.NewStatusError(messaging.StatusInvalidArgument, err, "")
	}
	req.Header = h
	resp, err = httpx.ExchangeWithTimeout(agent.Timeout(), agent.Exchange())(req)
	if err != nil {
		status = messaging.NewStatusError(resp.StatusCode, err, "")
		return
	}
	status = messaging.StatusOK()
	if agent.Log() {
		access.Log(access.EgressTraffic, start, time.Since(start), "", req, resp, access.Threshold{Timeout: agent.Timeout()})
	}
	return
}
