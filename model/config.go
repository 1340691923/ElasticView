//应用启动引擎层
package model

var GlobConfig Config

//全局配置结构体
type Config struct {
	Log struct {
		StorageDays int    `json:"storageDays"`
		LogDir      string `json:"logDir"`
	} `json:"log"`
	Port   int    `json:"port"`
	DbType string `json:"dbType"`
	Sqlite struct {
		DbPath string `json:"dbPath"`
	} `json:"sqlite"`
	Mysql struct {
		Username     string `json:"username"`
		Pwd          string `json:"pwd"`
		IP           string `json:"ip"`
		Port         string `json:"port"`
		DbName       string `json:"dbName"`
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns"`
	} `json:"mysql"`
	AppSecret string `json:"appSecret"`
	Version string `json:"version"`
	DeBug     bool   `json:"deBug"`
	EsPwdSecret string `json:"esPwdSecret"`
}

const (
	MysqlDbTyp  = "mysql"
	SqliteDbTyp = "sqlite3"
)

func (this *Config) GetDbType() string {
	if this.DbType == "mysql" {
		return MysqlDbTyp
	}
	return SqliteDbTyp
}
