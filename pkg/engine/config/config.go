package config

import (
	"flag"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/spf13/viper"
)

var (
	GlobConfig     Config
	appName        string
	configFileDir  string
	configFileName string
	configFileExt  string
)

const (
	MysqlDbTyp  = "mysql"
	SqliteDbTyp = "sqlite3"
)

func init() {
	flag.StringVar(&appName, "appName", "ElasticView", "应用名")
	flag.StringVar(&configFileDir, "configFileDir",  "resources/config","配置文件夹名")
	flag.StringVar(&configFileName, "configFileName", "config", "配置文件名")
	flag.StringVar(&configFileExt, "configFileExt", "yml", "配置文件后缀")
	flag.Parse()
	core.Register(core.MaxLevel, "读取配置文件", initConfig)
}

// 初始化配置
func initConfig() (deferFn func(), err error) {
	deferFn = func() {}
	config := viper.New()
	config.AddConfigPath(configFileDir)
	config.SetConfigName(configFileName)
	config.SetConfigType(configFileExt)
	if err := config.ReadInConfig(); err != nil {
		return deferFn, err
	}
	if err := config.Unmarshal(&GlobConfig); err != nil {
		return deferFn, err
	}
	return deferFn, nil
}

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
	AppSecret   string `json:"appSecret"`
	Version     string `json:"version"`
	DeBug       bool   `json:"deBug"`
	EsPwdSecret string `json:"esPwdSecret"`
}

func (this *Config) GetDbType() string {
	if this.DbType == "mysql" {
		return MysqlDbTyp
	}
	return SqliteDbTyp
}
