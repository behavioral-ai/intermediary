package cachetest

import (
	"github.com/behavioral-ai/core/httpx"
	"net/http"
)

var (
	respCache = httpx.NewResponseCache()
)

// Exchange - HTTP exchange function
func Exchange(r *http.Request) (resp *http.Response, err error) {
	switch r.Method {
	case http.MethodGet:
		resp = respCache.Get(r.URL.String())
	case http.MethodPut:
		respCache.Put(r.URL.String(), httpx.CreateResponse(r))
		resp = httpx.NewResponse(http.StatusOK, nil, nil)
	default:
		resp = httpx.NewResponse(http.StatusMethodNotAllowed, nil, nil)
	}
	return
}
