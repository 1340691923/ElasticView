package dto

import "github.com/1340691923/ElasticView/pkg/infrastructure/dto/common"

type EsDocUpdateByID struct {
	EsConnect int         `json:"es_connect"`
	ID        string      `json:"id"`
	JSON      common.Json `json:"json"`
	Type      string      `json:"type_name"`
	Index     string      `json:"index"`
}

type EsDocDeleteRowByID struct {
	EsConnect int    `json:"es_connect"`
	ID        string `json:"id"`
	IndexName string `json:"index_name"`
	Type      string `json:"type"`
}
