package proto

import (
	"context"
	"io"
	"net/http"
	"time"
)

type UpdateRequest struct {
	Index        string
	DocumentType string
	DocumentID   string

	Body io.Reader

	IfPrimaryTerm       *int
	IfSeqNo             *int
	Lang                string
	Refresh             string
	RequireAlias        *bool
	RetryOnConflict     *int
	Routing             string
	Source              []string
	SourceExcludes      []string
	SourceIncludes      []string
	Timeout             time.Duration
	WaitForActiveShards string

	Pretty     bool
	Human      bool
	ErrorTrace bool
	FilterPath []string

	Header http.Header

	ctx context.Context
}
