package proto

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type Response struct {
	statusCode int
	header     http.Header
	resByte    []byte
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) ResByte() []byte {
	return r.resByte
}

func (r *Response) JsonRawMessage() json.RawMessage {
	res := json.RawMessage{}
	json.Unmarshal(r.resByte, &res)
	return res
}

func NewResponse(statusCode int, header http.Header, readCloser io.ReadCloser) (res *Response, err error) {
	res = new(Response)
	defer readCloser.Close()
	var resByte []byte
	resByte, err = io.ReadAll(readCloser)
	if err != nil {
		return
	}

	res.resByte = resByte
	res.header = header
	res.statusCode = statusCode
	return
}

func (r *Response) StatusErr() (err error) {
	if gjson.GetBytes(r.resByte, "status").Int() > 201 {
		err = errors.New(string(r.resByte))
		return
	}
	return
}
