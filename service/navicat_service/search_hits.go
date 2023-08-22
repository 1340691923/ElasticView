package navicat_service

import "encoding/json"

type SearchHit struct {
	Score          *float64               `json:"_score,omitempty"`   // computed score
	Index          string                 `json:"_index,omitempty"`   // index name
	Type           string                 `json:"_type,omitempty"`    // type meta field
	Id             string                 `json:"_id,omitempty"`      // external or internal
	Uid            string                 `json:"_uid,omitempty"`     // uid meta field (see MapperService.java for all meta fields)
	Routing        string                 `json:"_routing,omitempty"` // routing meta field
	Parent         string                 `json:"_parent,omitempty"`  // parent meta field
	Version        *int64                 `json:"_version,omitempty"` // version number, when Version is set to true in SearchService
	SeqNo          *int64                 `json:"_seq_no"`
	PrimaryTerm    *int64                 `json:"_primary_term"`
	Sort           []interface{}          `json:"sort,omitempty"`            // sort information
	Source         *json.RawMessage       `json:"_source,omitempty"`         // stored document source
	Fields         map[string]interface{} `json:"fields,omitempty"`          // returned (stored) fields
	MatchedQueries []string               `json:"matched_queries,omitempty"` // matched queries
	Shard          string                 `json:"_shard,omitempty"`          // used e.g. in Search Explain
	Node           string                 `json:"_node,omitempty"`           // used e.g. in Search Explain
}
