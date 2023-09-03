package dto

//一些需要用到的结构

type Json map[string]interface{}

type Sort struct {
	Field     string
	Ascending bool
}

type Page struct {
	PageNum  int
	PageSize int
}

type EsConnectID struct {
	EsConnectID int `json:"es_connect"`
}
