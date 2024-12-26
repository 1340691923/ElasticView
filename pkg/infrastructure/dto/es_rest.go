package dto

type EsRest struct {
	EsConnect int    `json:"es_connect"`
	Body      string `json:"body"`
	Path      string `json:"path"`
}
