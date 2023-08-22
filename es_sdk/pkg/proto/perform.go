package proto

import (
	"net/http"
	"net/url"
)

type PerformRequestOptions struct {
	Method          string
	Path            string
	Params          url.Values
	Body            interface{}
	ContentType     string
	IgnoreErrors    []int
	Headers         http.Header
	MaxResponseSize int64
	Stream          bool
}
