package dto

type EsRest struct {
	EsConnect int    `json:"es_connect"`
	Method    string `json:"method"`
	Body      string `json:"body"`
	Path      string `json:"path"`
}
