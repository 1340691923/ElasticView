package proto

import (
	"context"
	"io"
	"net/http"
	"time"
)

type CreateRequest struct {
	Index        string
	DocumentType string
	DocumentID   string

	Body io.Reader

	Pipeline            string
	Refresh             string
	Routing             string
	Timeout             time.Duration
	Version             *int
	VersionType         string
	WaitForActiveShards string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}
