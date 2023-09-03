package proto

import (
	"context"
	"net/http"
	"time"
)

type DeleteRequest struct {
	Index        string
	DocumentType string
	DocumentID   string

	IfPrimaryTerm       *int
	IfSeqNo             *int
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
