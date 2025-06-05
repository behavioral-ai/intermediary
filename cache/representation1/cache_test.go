package representation1

import (
	"fmt"
	"github.com/behavioral-ai/collective/resource"
	"time"
)

const (
	NamespaceName = "test:resiliency:agent/cache/request/http"
)

var (
	m = map[string]string{
		HostKey:         "www.google.com",
		CacheControlKey: "no-store, no-cache, max-age=0",
		TimeoutKey:      "750ms",
		IntervalKey:     "4m",
		SundayKey:       "13-15",
		MondayKey:       "8-16",
		TuesdayKey:      "6-10",
		WednesdayKey:    "12-12",
		ThursdayKey:     "0-23",
		FridayKey:       "22-23",
		SaturdayKey:     "3-8",
	}

	m2 = map[string]string{
		"host":          "www.google.com",
		"cache-control": "no-store, no-cache, max-age=0",
		"sun":           "13-15",
		"mon":           "8-16",
		"tue":           "14-21",
		"wed":           "12-12",
		"thu":           "19-23",
		"fri":           "22-23",
		"sat":           "3-8",
	}
)

func ExampleParseCache() {
	var cache Cache
	parseCache(&cache, m)

	fmt.Printf("test: parseCache() -> %v\n", cache)

	//Output:
	//test: parseCache() -> {false <nil> 750ms 4m0s www.google.com map[Cache-Control:[no-store, no-cache, max-age=0]] map[fri:{22 23} mon:{8 16} sat:{3 8} sun:{13 15} thu:{0 23} tue:{6 10} wed:{12 12}]}

}

func _ExampleNewCache() {
	resource.NewAgent()

	status := resource.Resolver.AddRepresentation(NamespaceName, Fragment, "author", m)
	fmt.Printf("test: AddRepresentation() -> [status:%v]\n", status)

	ct, status2 := resource.Resolver.Representation(NamespaceName)
	fmt.Printf("test: Representation() -> [ct:%v] [status:%v]\n", ct, status2)

	if buf, ok := ct.Value.([]byte); ok {
		fmt.Printf("test: Representation() -> [value:%v] [status:%v]\n", len(buf), status2)
	}

	//l := NewCache(NamespaceName)
	//fmt.Printf("test: NewCache() -> %v\n", l)

	//Output:
	//test: AddRepresentation() -> [status:OK]
	//test: Representation() -> [ct:fragment: v1 type: application/json value: true] [status:OK]
	//test: Representation() -> [value:200] [status:OK]
	//test: NewCache() -> &{false 0xc00010e138 www.google.com 750ms 4m0s map[Cache-Control:[no-store, no-cache, max-age=0]] map[fri:{22 23} mon:{8 16} sat:{3 8} sun:{13 15} thu:{0 23} tue:{6 10} wed:{12 12}]}

}

func _ExampleRange() {
	s := "3-15"
	r := NewRange(s)
	fmt.Printf("test: NewRange(\"%v\") -> %v\n", s, r)

	s = " 3-23 "
	r = NewRange(s)
	fmt.Printf("test: NewRange(\"%v\") -> %v\n", s, r)

	ts := time.Now().UTC()
	ok := r.In(ts)
	fmt.Printf("test: NewRange(\"%v\") -> [hour:%v] [in:%v]\n", s, ts.Hour(), ok)

	s = "3-17"
	r = NewRange(s)
	fmt.Printf("test: NewRange(\"%v\") -> %v\n", s, r)

	ts = time.Now().UTC()
	ok = r.In(ts)
	fmt.Printf("test: NewRange(\"%v\") -> [hour:%v] [in:%v]\n", s, ts.Hour(), ok)

	s = "19-23"
	r = NewRange(s)
	fmt.Printf("test: NewRange(\"%v\") -> %v\n", s, r)

	ts = time.Now().UTC()
	ok = r.In(ts)
	fmt.Printf("test: NewRange(\"%v\") -> [hour:%v] [in:%v]\n", s, ts.Hour(), ok)

	//Output:
	//test: NewRange("3-15") -> {3 15}
	//test: NewRange(" 3-23 ") -> {3 23}
	//test: NewRange(" 3-23 ") -> [hour:18] [in:true]
	//test: NewRange("3-17") -> {3 17}
	//test: NewRange("3-17") -> [hour:18] [in:false]
	//test: NewRange("19-23") -> {19 23}
	//test: NewRange("19-23") -> [hour:18] [in:false]

}
