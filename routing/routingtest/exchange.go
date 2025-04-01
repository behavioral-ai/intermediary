package routingtest

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"strings"
)

const (
	googlePath = "/google/search"
	yahooPath  = "/yahoo/search"
)

// Exchange - HTTP exchange function
func Exchange(r *http.Request) (resp *http.Response, err error) {
	uri := ""
	values := r.URL.Query()
	q := values.Encode()
	if strings.HasPrefix(r.URL.Path, googlePath) {
		uri = "https://www.google.com/search?" + q
	} else {
		if strings.HasPrefix(r.URL.Path, yahooPath) {
			uri = "https://search.yahoo.com/search?" + q
		} else {
			return httpx.NewResponse(http.StatusBadRequest, nil, nil), err
		}
	}
	h := make(http.Header)
	h.Add(iox.AcceptEncoding, iox.GzipEncoding)
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header = h
	resp, err = httpx.Do(req)
	if err != nil {
		fmt.Printf("test: httx.Do() -> [err:%v]\n", err)
	}
	return
}
