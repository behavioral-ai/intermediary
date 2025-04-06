package routing

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/config"
	"net/http"
	"time"
)

func ExampleNew() {
	a := newAgent(eventtest.New())

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())

	m := make(map[string]string)
	m[config.AppHostKey] = "google.com"
	a.Message(messaging.NewConfigMapMessage(m))
	time.Sleep(time.Second * 2)
	fmt.Printf("test: Message() -> %v\n", a.hostName)

	//Output:
	//test: newAgent() -> resiliency:agent/behavioral-ai/intermediary/routing
	//test: Message() -> google.com

}

func ExampleExchange() {
	url := "http://localhost:8080/search?q=golang"
	a := newAgent(eventtest.New())
	ex := a.Link(nil)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err := ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	a.hostName = "www.google.com"
	req, _ = http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add(httpx.XRequestId, "1234-request-id")
	resp, err = ex(req)
	fmt.Printf("test: Exchange() -> [resp:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//notify-> 2025-03-25T14:44:49.521Z [resiliency:agent/behavioral-ai/resiliency/routing] [core:messaging.status] [] [Invalid Argument] [host configuration is empty]
	//test: Exchange() -> [resp:500] [err:host configuration is empty]
	//test: Exchange() -> [resp:200] [err:<nil>]

}
