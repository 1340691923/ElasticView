package dto

import "github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"

type ExecReq struct {
	PluginId string        `json:"plugin_id"`
	Sql      string        `json:"sql"`
	Args     []interface{} `json:"args"`
}

type ExecMoreReq struct {
	PluginId string    `json:"plugin_id"`
	Sqls     []ExecSql `json:"sqls"`
}

type ExecSql struct {
	Sql  string        `json:"sql"`
	Args []interface{} `json:"args"`
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

type SaveDb struct {
	PluginId  string                 `json:"plugin_id"`
	TableName string                 `json:"table"` // 目标表名
	Data      map[string]interface{} `json:"data"`  // 要插入或更新的数据
}

type UpdateDb struct {
	PluginId   string                 `json:"plugin_id"`
	TableName  string                 `json:"table"` // 目标表名
	UpdateSql  string                 `json:"update_sql"`
	UpdateArgs []interface{}          `json:"update_args"`
	Data       map[string]interface{} `json:"data"` // 要插入或更新的数据
}

type DeleteDb struct {
	PluginId  string        `json:"plugin_id"`
	TableName string        `json:"table"` // 目标表名
	WhereSql  string        `json:"where_sql"`
	WhereArgs []interface{} `json:"where_args"`
}

type InsertOrUpdateDb struct {
	PluginId   string                 `json:"plugin_id"`
	TableName  string                 `json:"table"` // 目标表名
	UpsertData map[string]interface{} // 没有则新增，有则更新
	UniqueKeys []string               // 冲突检查的唯一键
}
