package proto

import (
	"context"
	"net/http"
	"time"
)

type CatIndicesRequest struct {
	Index []string

	Bytes                   string
	ExpandWildcards         string
	Format                  string
	H                       []string
	Health                  string
	Help                    *bool
	IncludeUnloadedSegments *bool
	Local                   *bool
	MasterTimeout           time.Duration
	Pri                     *bool
	S                       []string
	Time                    string
	V                       *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

type CatHealthRequest struct {
	Format string
	H      []string
	Help   *bool
	S      []string
	Time   string
	Ts     *bool
	V      *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

type CatShardsRequest struct {
	Index []string

	Bytes         string
	Format        string
	H             []string
	Help          *bool
	Local         *bool
	MasterTimeout time.Duration
	S             []string
	Time          string
	V             *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

type CatCountRequest struct {
	Index []string

	Format string
	H      []string
	Help   *bool
	S      []string
	V      *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

type CatAllocationRequest struct {
	NodeID []string

	Bytes         string
	Format        string
	H             []string
	Help          *bool
	Local         *bool
	MasterTimeout time.Duration
	S             []string
	V             *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}

type CatAliasesRequest struct {
	Name []string

	ExpandWildcards string
	Format          string
	H               []string
	Help            *bool
	Local           *bool
	S               []string
	V               *bool

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}
