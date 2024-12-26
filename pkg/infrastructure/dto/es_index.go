package dto

import "github.com/1340691923/ElasticView/pkg/infrastructure/dto/common"

type EsIndexInfo struct {
	EsConnect int         `json:"es_connect"`
	Settings  common.Json `json:"settings"`
	IndexName string      `json:"index_name"`
	Types     string      `json:"types"`
}
