package model

type DataxLinkInfoModel struct {
	Id       int    `db:"id" json:"id"`
	Ip       string `db:"ip" json:"ip"`
	Port     int    `db:"port" json:"port"`
	DbName   string `db:"db_name" json:"db_name"`
	Username string `db:"username" json:"username"`
	Pwd      string `db:"pwd" json:"pwd"`
	Remark   string `db:"remark" json:"remark"`
	Typ      string `db:"typ" json:"typ"`
	Updated  string `db:"updated" json:"updated"`
	Created  string `db:"created" json:"created"`
}
