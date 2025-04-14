package routingtest

import (
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/routing"
	"net/http"
	"net/http/httptest"
)

func _ExampleSearchExchange() {
	agent := exchange.Agent(routing.NamespaceName)

	// configure exchange and host name
	agent.Message(httpx.NewConfigExchangeMessage(searchExchange))
	cfg := make(map[string]string)
	cfg[config.AppHostKey] = "localhost:8080"
	agent.Message(messaging.NewConfigMapMessage(cfg))

	// create request
	url := "https://localhost:8081/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)

	// create endpoint and server http
	e := host.NewEndpoint(agent)
	r := httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()

	// decoding when read all
	buf, err := iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: RoutingAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	r = httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()

	// not decoding when read all
	buf, err = iox.ReadAll(r.Result().Body, nil)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: RoutingAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	//Output:
	//test: RoutingAgent [status:200 ] [encoding:] [buff:82980]
	//test: RoutingAgent [status:200 ] [encoding:gzip] [buff:41075]

}

func ExampleEchoExchange() {
	agent := exchange.Agent(routing.NamespaceName)

	// configure exchange and host name
	agent.Message(httpx.NewConfigExchangeMessage(EchoExchange))
	cfg := make(map[string]string)
	cfg[config.AppHostKey] = "localhost:8080"
	agent.Message(messaging.NewConfigMapMessage(cfg))

	// create request
	url := "https://localhost:8081/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)

	// create endpoint and run
	e := host.NewEndpoint(agent)
	r := httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()

	// decoding when read all
	buf, err := iox.ReadAll(r.Result().Body, r.Result().Header)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [content-type:%v] [err:%v]\n", len(buf), http.DetectContentType(buf), err)
	fmt.Printf("test: RoutingAgent [status:%v ] [encoding:%v] [%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), string(buf))

	r = httptest.NewRecorder()
	e.ServeHTTP(r, req)
	r.Flush()

	// not decoding when read all
	buf, err = iox.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: iox.ReadAll() -> [buf:%v] [content-type:%v] [err:%v]\n", len(buf), http.DetectContentType(buf), err)
	fmt.Printf("test: RoutingAgent [status:%v ] [encoding:%v] [%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	//Output:
	//test: RoutingAgent [status:200 ] [encoding:] [buff:82980]
	//test: RoutingAgent [status:200 ] [encoding:gzip] [buff:41075]

}
