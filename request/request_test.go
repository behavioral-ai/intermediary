package request

import (
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"net/http"
	"time"
)

type agentT struct {
	timeout  time.Duration
	exchange rest.Exchange
}

func (a *agentT) String() string               { return a.Uri() }
func (a *agentT) Uri() string                  { return "agent:request" }
func (a *agentT) Message(m *messaging.Message) {}
func (a *agentT) Route() string                { return "cache" }
func (a *agentT) Log() bool                    { return true }
func (a *agentT) Timeout() time.Duration       { return a.timeout }
func (a *agentT) Do() rest.Exchange            { return a.exchange }

func ExampleDo_Get() {
	url := "https://www.google.com/search?q=golang"
	a := new(agentT)
	a.exchange = httpx.Do

	h := make(http.Header)
	h.Add(iox.AcceptEncoding, iox.GzipEncoding)
	h.Add(httpx.XRequestId, "1234-request-id")
	resp, status := Do(a, http.MethodGet, url, h, nil)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp.StatusCode, status)

	if resp.StatusCode == http.StatusOK {
		buf, err1 := iox.ReadAll(resp.Body, resp.Header)
		fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err1)
	}

	//Output:
	//test: do() -> [resp:200] [status:<nil>]
	//test: iox.ReadAll() -> [buf:82676] [err:<nil>]

}

func ExampleDo_Get_Timeout() {
	url := "https://www.google.com/search?q=golang"

	a := new(agentT)
	a.exchange = httpx.Do
	a.timeout = time.Second * 10

	h := make(http.Header)
	h.Add(iox.AcceptEncoding, "gzip")
	h.Add(httpx.XRequestId, "1234-request-id")
	resp, status := Do(a, http.MethodGet, url, h, nil)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp.StatusCode, status)

	if resp.StatusCode == http.StatusOK {
		buf, err1 := iox.ReadAll(resp.Body, resp.Header)
		fmt.Printf("test: iox.ReadAll() -> [buf:%v] [err:%v]\n", len(buf), err1)
	}

	a.timeout = time.Millisecond * 10
	resp, status = Do(a, http.MethodGet, url, h, nil)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp.StatusCode, status)

	//Output:
	//test: Do() -> [resp:200] [status:OK]
	//test: iox.ReadAll() -> [buf:82131] [err:<nil>]
	//test: Do() -> [resp:504] [status:Timeout [err:Get "https://www.google.com/search?q=golang": context deadline exceeded]]

}
