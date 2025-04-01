package cachetest

import (
	"github.com/behavioral-ai/core/httpx"
	"net/http"
)

// ExchangeWriter - HTTP exchange writer function
func ExchangeWriter(w http.ResponseWriter, r *http.Request) {
	var (
		resp *http.Response
	)

	switch r.Method {
	case http.MethodGet:
		resp, _ = Exchange(r)
		if resp.StatusCode == http.StatusNotFound {
			w.WriteHeader(resp.StatusCode)
			return
		}
		httpx.WriteResponse(w, resp.Header, resp.StatusCode, resp.Body, nil)
	case http.MethodPut:
		resp, _ = Exchange(r)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
