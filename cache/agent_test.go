package cache

import (
	"fmt"
	"github.com/behavioral-ai/collective/operations/operationstest"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/cache/representation1"
	"net/http"
)

func ExampleNew() {
	//url := "https://www.google.com/search"
	a := newAgent(representation1.Initialize(nil), nil, operationstest.NewService())

	fmt.Printf("test: newAgent() -> %v\n", a.Name())
	m := make(map[string]string)
	m[representation1.HostKey] = "google.com"
	a.Message(messaging.NewConfigMapMessage(m))
	fmt.Printf("test: Message() -> %v\n", a.state.Host)

	//Output:
	//test: newAgent() -> test:resiliency:agent/cache/request/http
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
