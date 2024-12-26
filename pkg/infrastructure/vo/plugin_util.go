package vo

type ExecSqlRes struct {
	RowsAffected int64 `json:"rows_affected"`
}

type SelectSqlRes struct {
	Result []map[string]interface{} `json:"result"`
}

type FirstSqlRes struct {
	Result map[string]interface{} `json:"result"`
}
