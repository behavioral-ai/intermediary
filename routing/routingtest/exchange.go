package routingtest

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
)

// Exchange - HTTP exchange function
func Exchange(r *http.Request) (resp *http.Response, err error) {
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
