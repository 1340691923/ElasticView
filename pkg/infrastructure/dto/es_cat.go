package dto

type EsCat struct {
	EsConnect        int    `json:"es_connect"`
	Cat              string `json:"cat"`
	IndexBytesFormat string `json:"index_bytes_format"`
}
