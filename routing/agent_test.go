package routing

import (
	"fmt"
	"github.com/behavioral-ai/core/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/routing/representation1"
	"net/http"
	"time"
)

func ExampleNew() {
	a := newAgent(eventtest.New(), nil)

	fmt.Printf("test: newAgent() -> %v\n", a.Name())

	m := make(map[string]string)
	m[representation1.AppHostKey] = "google.com"
	a.Message(messaging.NewConfigMapMessage(m))
	time.Sleep(time.Second * 2)
	rt, ok := a.router.Lookup(defaultRoute)
	fmt.Printf("test: Message() -> [%v] [uri:%v] [ok:%v]\n", rt.Name, rt.Uri, ok)

	//Output:
	//test: newAgent() -> resiliency:agent/routing/request/http
	//test: Message() -> [core:routing/default] [uri:google.com] [ok:true]

}

func ExampleExchange() {
	url := "http://localhost:8080/search?q=golang"
	a := newAgent(eventtest.New(), nil)
	ex := a.Exchange

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err := ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	rt, _ := a.router.Lookup(defaultRoute)
	rt.Uri = "www.google.com"
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err = ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//notify-> 2025-03-25T14:44:49.521Z [resiliency:agent/routing/request/http [core:messaging.status] [] [Invalid Argument] [host configuration is empty]
	//test: Exchange() -> [resp:500] [err:host configuration is empty]
	//test: Exchange() -> [resp:200] [err:<nil>]

}
