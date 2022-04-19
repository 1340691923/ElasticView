package model

type DataxListModel struct {
	Id         int    `json:"id" db:"id"`
	FormData   string `json:"form_data" db:"form_data"`
	Remark     string `json:"remark" db:"remark"`
	Table_name string `json:"table_name" db:"table_name"`
	Index_name string `json:"index_name" db:"index_name"`
	Error_msg  string `json:"error_msg" db:"error_msg"`
	Status     string `json:"status" db:"status"`
	Updated    string `json:"updated" db:"updated"`
	Created    string `json:"created" db:"created"`
}
