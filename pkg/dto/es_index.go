package dto

type EsIndexInfo struct {
	EsConnect int    `json:"es_connect"`
	Settings  Json   `json:"settings"`
	IndexName string `json:"index_name"`
	Types     string `json:"types"`
}
