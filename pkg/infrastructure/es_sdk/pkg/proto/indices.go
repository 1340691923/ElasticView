package proto

import (
	"io"
	"net/http"
	"time"
)

type IndicesPutSettingsRequest struct {
	Index []string

	Body io.Reader

	AllowNoIndices    *bool
	ExpandWildcards   string
	FlatSettings      *bool
	IgnoreUnavailable *bool
	MasterTimeout     time.Duration
	PreserveExisting  *bool
	Timeout           time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}

type IndicesCreateRequest struct {
	Index string

	Body io.Reader

	IncludeTypeName     *bool
	MasterTimeout       time.Duration
	Timeout             time.Duration
	WaitForActiveShards string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}

type IndicesDeleteRequest struct {
	Index []string

	AllowNoIndices    *bool
	ExpandWildcards   string
	IgnoreUnavailable *bool
	MasterTimeout     time.Duration
	Timeout           time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}

type ReindexRequest struct {
	Body io.Reader

	MaxDocs             *int
	Refresh             *bool
	RequestsPerSecond   *int
	Scroll              time.Duration
	Slices              interface{}
	Timeout             time.Duration
	WaitForActiveShards string
	WaitForCompletion   *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}

type IndicesGetSettingsRequest struct {
	Index []string

	Name []string

	AllowNoIndices    *bool
	ExpandWildcards   string
	FlatSettings      *bool
	IgnoreUnavailable *bool
	IncludeDefaults   *bool
	Local             *bool
	MasterTimeout     time.Duration

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}

type IndicesPutMappingRequest struct {
	Index        []string
	DocumentType string

	Body io.Reader

	AllowNoIndices    *bool
	ExpandWildcards   string
	IgnoreUnavailable *bool
	IncludeTypeName   *bool
	MasterTimeout     time.Duration
	Timeout           time.Duration
	WriteIndexOnly    *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header
}
