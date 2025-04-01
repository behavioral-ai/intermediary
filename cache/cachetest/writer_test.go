package cachetest

/*
func ExampleExchange() {
	agent := cache.New(eventtest.New(nil))
	agent.exchange = cachingExchange
	agent.hostName = "localhost:8082"

	url := "https://localhost:8081/search/google?q=golang"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = make(http.Header)
	httpx.AddRequestId(req)

	chain := httpx.BuildChain(agent, routingExchange)
	r := httptest.NewRecorder()
	host.Exchange(r, req, chain)
	r.Flush()
	buf, err := iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	r = httptest.NewRecorder()
	host.Exchange(r, req, chain)
	r.Flush()
	buf, err = iox.ReadAll(r.Result().Body, r.Result().Header)
	if err != nil {
		fmt.Printf("test: iox.RedAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: CacheAgent [status:%v ] [encoding:%v] [buff:%v]\n", r.Result().StatusCode, r.Result().Header.Get(iox.ContentEncoding), len(buf))

	//Output:
	//fail

}


*/
