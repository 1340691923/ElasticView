package dto

import "github.com/1340691923/ElasticView/pkg/infrastructure/dto/common"

type EsMapGetProperties struct {
	EsConnectID int    `json:"es_connect"`
	IndexName   string `json:"index_name"`
}

type UpdateMapping struct {
	EsConnect  int         `json:"es_connect"`
	IndexName  string      `json:"index_name"`
	TypeName   string      `json:"type_name"`
	Properties common.Json `json:"properties"`
}

type EsMappingInfo struct {
	IndexNameList []string    `json:"index_name_list"`
	EsConnect     int         `json:"es_connect"`
	Mappings      common.Json `json:"mappings"`
	IndexName     string      `json:"index_name"`
}
