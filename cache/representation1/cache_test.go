package representation1

import (
	"fmt"
	"time"
)

var (
	m = map[string]string{
		"host":          "www.google.com",
		"cache-control": "no-store, no-cache, max-age=0",
		"sun":           "13-15",
		"mon":           "8-16",
		"tue":           "6-10",
		"wed":           "12-12",
		"thu":           "0-23",
		"fri":           "22-23",
		"sat":           "3-8",
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

func Example_ParseCache() {
	c := Initialize()
	parseCache(c, m)

	fmt.Printf("test: parseCache() -> %v\n", c)

	//Output:
	//test: parseCache() -> &{false www.google.com map[Cache-Control:[no-store, no-cache, max-age=0]] map[fri:{22 23} mon:{8 16} sat:{3 8} sun:{13 15} thu:{0 23} tue:{6 10} wed:{12 12}]}

}

func Example_NewCache() {
	c := newCache("", m)
	fmt.Printf("test: newCache() -> [enabled:%v]\n", c.Host != "")

	fmt.Printf("test: newCache() -> [now:%v]\n", c.Now())

	c = newCache("", m2)
	fmt.Printf("test: newCache() -> [now:%v]\n", c.Now())

	//Output:
	//test: newCache() -> [enabled:true]
	//test: newCache() -> [now:false]
	//test: newCache() -> [now:true]

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
