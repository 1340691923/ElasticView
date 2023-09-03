package config

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

const (
	MysqlDbTyp  = "mysql"
	SqliteDbTyp = "sqlite3"
)

type CommandLineArgs struct {
	HomePath    string
	ConfigFile  string
	Profile     bool
	ProfileAddr string
	ProfilePort uint64
	Tracing     bool
	TracingFile string
	CmdName     string
}

type Config struct {
	HomePath string
	CmdName  string

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

func (cfg *Config) CreateDbSource() string {
	var dbSource string
	if cfg.DbType == SqliteDbTyp {
		dataDir := filepath.Join(util.GetCurrentDirectory(), "data")
		if !util.CheckFileIsExist(dataDir) {
			os.MkdirAll(dataDir, os.ModePerm)
		}
		dbSource = filepath.Join(dataDir, cfg.Sqlite.DbPath) + "?_loc=Local&_busy_timeout=9999999"
	} else {
		dbSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			cfg.Mysql.Username,
			cfg.Mysql.Pwd,
			cfg.Mysql.IP,
			cfg.Mysql.Port,
			cfg.Mysql.DbName,
		)
	}

	return dbSource
}

func InitConfig(opt *CommandLineArgs) (cfg *Config, err error) {
	cfg = new(Config)
	cf := filepath.Join(opt.HomePath, opt.ConfigFile)
	if filepath.IsAbs(opt.ConfigFile) {
		cf = opt.ConfigFile
	}

	vip := viper.New()
	vip.SetConfigFile(cf)
	if err := vip.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "ReadInConfig")
	}
	if err := vip.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "config Unmarshal err")
	}
	cfg.HomePath = opt.HomePath
	cfg.CmdName = opt.CmdName
	log.Println("配置文件加载成功", cf)
	return cfg, nil
}
