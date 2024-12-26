package dto

import "github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"

type ExecReq struct {
	PluginId string        `json:"plugin_id"`
	Sql      string        `json:"sql"`
	Args     []interface{} `json:"args"`
}

type SelectReq struct {
	PluginId string        `json:"plugin_id"`
	Sql      string        `json:"sql"`
	Args     []interface{} `json:"args"`
}

type LoadDebugPlugin struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
	Pid  int    `json:"pid"`
}
type StopDebugPlugin struct {
	ID string `json:"id"`
}

type PluginRunDsl struct {
	EsConnectData dto.EsConnectData `json:"es_connect_data"`
	HttpMethod    string            `json:"http_method"`
	Path          string            `json:"path"`
	Dsl           string            `json:"dsl"`
}
