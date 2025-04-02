package routingtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
	"strings"
)

const (
	googlePath = "/google/search"
	yahooPath  = "/yahoo/search"
)

// SearchExchange - HTTP exchange function
func searchExchange(r *http.Request) (resp *http.Response, err error) {
	ctx := context.Background()
	uri := ""
	values := r.URL.Query()
	q := values.Encode()
	if strings.HasPrefix(r.URL.Path, googlePath) {
		uri = "https://www.google.com/search?" + q
	} else {
		if strings.HasPrefix(r.URL.Path, yahooPath) {
			uri = "https://search.yahoo.com/search?" + q
		} else {
			return httpx.NewResponse(http.StatusBadRequest, nil, nil), err
		}
	}
	h := make(http.Header)
	h.Add(iox.AcceptEncoding, iox.GzipEncoding)
	if r.Context() != nil {
		ctx = r.Context()
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	req.Header = h
	resp, err = httpx.Do(req)
	if err != nil {
		fmt.Printf("test: httx.Do() -> [err:%v]\n", err)
	}
	return
}

type echo struct {
	Method string      `json:"method"`
	Host   string      `json:"host"`
	Url    string      `json:"url"`
	Header http.Header `json:"header"`
}

// EchoExchange - HTTP exchange function
func EchoExchange(r *http.Request) (resp *http.Response, err error) {
	e := echo{
		Method: r.Method,
		Host:   r.Host,
		Url:    r.URL.String(),
		Header: r.Header,
	}
	buf, err1 := json.Marshal(e)
	if err1 != nil {
		return httpx.NewResponse(http.StatusBadRequest, nil, nil), err1
	}
	buf2, err2 := indent(buf)
	if err2 != nil {
		return httpx.NewResponse(http.StatusInternalServerError, nil, nil), err2
	}
	buf3, encoding, err3 := iox.EncodeContent(r.Header, buf2)
	if err3 != nil {
		return httpx.NewResponse(http.StatusInternalServerError, nil, nil), err3
	}
	if encoding != "" {
		resp = httpx.NewResponse(http.StatusOK, nil, buf3)
		if resp.Header == nil {
			resp.Header = make(http.Header)
		}
		resp.Header.Add(iox.ContentEncoding, encoding)
	} else {
		resp = httpx.NewResponse(http.StatusOK, nil, buf2)
	}
	return
}

func indent(src []byte) ([]byte, error) {
	var buf bytes.Buffer

	err := json.Indent(&buf, src, "", "  ")
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

/*
func encode(in []byte, h http.Header) (resp *http.Response, err error) {
	var w iox.EncodingWriter

	buf := new(bytes.Buffer)
	w, err = iox.NewEncodingWriter(buf, h)
	if err != nil {
		return httpx.NewResponse(http.StatusInternalServerError, nil, buf), err
	}
	cnt, err1 := w.Write(in)
	if err1 != nil || cnt != len(in) {
		return httpx.NewResponse(http.StatusInternalServerError, nil, buf), err1
	}
	rh := make(http.Header)
	rh.Add(iox.ContentEncoding, iox.GzipEncoding)
	return httpx.NewResponse(http.StatusOK, rh, buf.Bytes()), nil
}


*/
