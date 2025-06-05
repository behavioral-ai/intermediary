package routingtest

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/routing"
	_ "github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/intermediary/routing/representation1"
	"net/http"
	"net/http/httptest"
)

func ExampleExchange_Override() {
	cfg := make(map[string]string)
	cfg[representation1.AppHostKey] = "localhost:8080"

	routing.ConstructorOverride(cfg, Exchange, operationstest.NewService())
	agent := repository.Agent(routing.NamespaceName)

	// configure exchange and host name
	//agent.Message(httpx.NewConfigExchangeMessage(Exchange))
	//agent.Message(messaging.NewConfigMapMessage(cfg))

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

func _ExampleExchange() {
	cfg := make(map[string]string)
	cfg[representation1.AppHostKey] = "localhost:8080"

	agent := repository.Agent(routing.NamespaceName)
	// configure exchange and host name
	agent.Message(rest.NewExchangeMessage(Exchange))
	agent.Message(messaging.NewMapMessage(cfg))

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
