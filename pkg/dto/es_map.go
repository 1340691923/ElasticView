package dto

type EsMapGetProperties struct {
	EsConnectID int    `json:"es_connect"`
	IndexName   string `json:"index_name"`
}

type UpdateMapping struct {
	EsConnect  int    `json:"es_connect"`
	IndexName  string `json:"index_name"`
	TypeName   string `json:"type_name"`
	Properties Json   `json:"properties"`
}

type EsMappingInfo struct {
	IndexNameList []string `json:"index_name_list"`
	EsConnect     int      `json:"es_connect"`
	Mappings      Json     `json:"mappings"`
	IndexName     string   `json:"index_name"`
}
