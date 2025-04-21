package cache

import (
	"fmt"
	"github.com/behavioral-ai/collective/eventing/eventtest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/config"
	"net/http"
)

func ExampleNew() {
	//url := "https://www.google.com/search"
	a := newAgent(eventtest.New())

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())
	m := make(map[string]string)
	m[config.CacheHostKey] = "google.com"
	a.Message(messaging.NewConfigMapMessage(m))
	fmt.Printf("test: Message() -> %v\n", a.hostName)

	//Output:
	//test: newAgent() -> behavioral-ai.github.com:resiliency:agent/intermediary/cache
	//test: Message() -> google.com

}

func routingExchange(next rest.Exchange) rest.Exchange {
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
