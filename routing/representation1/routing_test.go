package representation1

import (
	"fmt"
	"github.com/behavioral-ai/collective/resource"
)

const (
	NamespaceName = "test:resiliency:agent/routing/request/http"
)

var (
	m = map[string]string{
		AppHostKey:  "www.google.com",
		LogKey:      "true",
		TimeoutKey:  "750ms",
		LogRouteKey: "app2",
	}
)

func ExampleParseRouting() {
	var routing Routing
	parseRouting(&routing, m)

	fmt.Printf("test: parseRouting() -> %v\n", routing)

	//Output:
	//test: parseRouting() -> {true www.google.com app2 750ms}

}

func _ExampleNewRouting() {
	resource.NewAgent()

	status := resource.Resolver.AddRepresentation(NamespaceName, Fragment, "author", m)
	fmt.Printf("test: AddRepresentation() -> [status:%v]\n", status)

	ct, status2 := resource.Resolver.Representation(NamespaceName)
	fmt.Printf("test: Representation() -> [ct:%v] [status:%v]\n", ct, status2)

	if buf, ok := ct.Value.([]byte); ok {
		fmt.Printf("test: Representation() -> [value:%v] [status:%v]\n", len(buf), status2)
	}

	//l := NewRouting(NamespaceName)
	//fmt.Printf("test: NewRouting() -> %v\n", l)

	//Output:
	//test: AddRepresentation() -> [status:OK]
	//test: Representation() -> [ct:fragment: v1 type: application/json value: true] [status:OK]
	//test: Representation() -> [value:80] [status:OK]
	//test: NewRouting() -> &{true www.google.com app2 750ms}

}
