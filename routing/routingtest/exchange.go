package routingtest

import (
	"bytes"
	"encoding/json"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/iox"
	"net/http"
)

type echo struct {
	Method string      `json:"method"`
	Host   string      `json:"host"`
	Url    string      `json:"url"`
	Header http.Header `json:"header"`
}

// Exchange - HTTP exchange function
func Exchange(r *http.Request) (resp *http.Response, err error) {
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
