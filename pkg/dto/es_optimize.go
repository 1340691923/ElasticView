package dto

type EsOptimize struct {
	EsConnect int    `json:"es_connect"`
	IndexName string `json:"index_name"`
	Command   string `json:"command"`
}
