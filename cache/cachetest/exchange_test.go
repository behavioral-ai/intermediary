package cachetest

import (
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/config"
	"net/http"
	"net/http/httptest"
)

func nextExchange(r *http.Request) (resp *http.Response, err error) {
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

func ExampleExchange() {
	agent := exchange.Agent(cache.NamespaceName)

	// configure exchange and host name
	agent.Message(httpx.NewConfigExchangeMessage(Exchange))
	cfg := make(map[string]string)
	cfg[config.CacheHostKey] = "localhost:8082"
	agent.Message(messaging.NewConfigMapMessage(cfg))

	// create request
	url := "https://localhost:8081/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)
	httpx.AddRequestId(req)

	// create endpoint and server Http
	e := host.NewEndpoint(agent, nextExchange)
	r := httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()
	buf, err := iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	r = httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()
	buf, err = iox.ReadAll(r.Result().Body, nil)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	//Output:
	//test: CacheAgent [status:200 ] [encoding:] [buff:82654]
	//test: CacheAgent [status:200 ] [encoding:gzip] [buff:41182]

}
