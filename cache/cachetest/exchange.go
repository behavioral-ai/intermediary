package cachetest

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
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

func NextExchange(next httpx.Exchange) httpx.Exchange {
	return func(r *http.Request) (resp *http.Response, err error) {
		h := make(http.Header)
		h.Add(iox.AcceptEncoding, iox.GzipEncoding)
		req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
		req.Header = h
		resp, err = httpx.Do(req)
		if err != nil {
			fmt.Printf("test: httx.Do() -> [err:%v]\n", err)
		}
		return
	}
}
