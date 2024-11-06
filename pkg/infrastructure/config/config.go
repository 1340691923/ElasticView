package config

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var TranslationCfg = map[string]map[string]interface{}{}

const (
	MysqlDbTyp  = "mysql"
	SqliteDbTyp = "sqlite3"
)

type CommandLineArgs struct {
	IconData    []byte
	HomePath    string
	ConfigFile  string
	Profile     bool
	ProfileAddr string
	ProfilePort uint64
	Tracing     bool
	TracingFile string
	CmdName     string
	AdminPwd    string
}

type Config struct {
	IconData       []byte
	HomePath       string
	CmdName        string
	StoreFileDir   string `json:"storeFileDir"`
	EnableLogEs    bool   `json:"enableLogEs"`
	EnableLogEsRes bool   `json:"enableLogEsRes"`
	RootUrl        string `json:"rootUrl"`
	Log            struct {
		StorageDays  int    `json:"storageDays"`
		LogDir       string `json:"logDir"`
		PluginLogDir string `json:"pluginLogDir"`
	} `json:"log"`
	CheckForevUpdates     bool   `json:"checkForevUpdates"`
	CheckForPluginUpdates bool   `json:"checkForPluginUpdates"`
	Port                  int    `json:"port"`
	PluginRpcPort         int    `json:"pluginRpcPort"`
	DbType                string `json:"dbType"`
	Sqlite                struct {
		DbPath string `json:"dbPath"`
		DbName string `json:"dbName"`
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
	AppSecret        string      `json:"appSecret"`
	Version          string      `json:"version"`
	DeBug            bool        `json:"deBug"`
	EsPwdSecret      string      `json:"esPwdSecret"`
	EvKey            string      `json:"evKey"` //易帷key
	Plugin           Plugin      `json:"plugin"`
	WatermarkContent string      `json:"watermarkContent"`
	Translation      Translation `json:"translation"`
	OAuth            OAuth       `json:"oAuth"`
}

type OAuth struct {
	WorkWechat WorkWechat `json:"workWechat"`
}

type WorkWechat struct {
	Corpid  string `json:"corpid"`
	AgentId string `json:"agentId"`
	Secert  string `json:"secert"`
	Enable  bool   `json:"enable"`
}

type Translation struct {
	Lang   string `json:"lang"`
	CfgDir string `json:"cfgDir"`
}

type Plugin struct {
	LoadPath  string `json:"loadPath"`
	StorePath string `json:"storePath"`
}

func (this *Plugin) Error() error {
	if this.LoadPath == "" {
		return errors.New("配置文件中插件配置文件夹没有添加")
	}
	if this.StorePath == "" {
		return errors.New("配置文件中插件数据存储目录配置没有添加")
	}
	if !util.CheckFileIsExist(this.LoadPath) {
		os.MkdirAll(this.LoadPath, os.ModePerm)
	}
	if !util.CheckFileIsExist(this.StorePath) {
		os.MkdirAll(this.StorePath, os.ModePerm)
	}

	return nil
}

func (this *Config) GetRootUrl() string {
	if this.RootUrl == "" {
		return "http://localhost:8090/"
	}

	if this.RootUrl[len(this.RootUrl)-1] != '/' {
		this.RootUrl += "/"
	}

	return this.RootUrl
}

func (this *Config) ParseAppUrlAndSubUrl() (string, string, error) {
	appUrl := this.GetRootUrl()

	if appUrl[len(appUrl)-1] != '/' {
		appUrl += "/"
	}

	url, err := url.Parse(appUrl)
	if err != nil {
		log.Println("err", err, appUrl)
		return "", "", err
	}

	appSubUrl := strings.TrimSuffix(url.Path, "/")
	return appUrl, appSubUrl, nil
}

func (this *Config) GetDbType() string {
	if this.DbType == "mysql" {
		return MysqlDbTyp
	}
	return SqliteDbTyp
}

func (cfg *Config) CreateDbDSN() string {
	var DSN string
	if cfg.DbType == SqliteDbTyp {
		dataDir := filepath.Join(util.GetCurrentDirectory(), "data")

		if cfg.Sqlite.DbPath != "" {
			dataDir = cfg.Sqlite.DbPath
		}

		if !util.CheckFileIsExist(dataDir) {
			os.MkdirAll(dataDir, os.ModePerm)
		}

		DSN = filepath.Join(dataDir, cfg.Sqlite.DbName) + "?_pragma=charset(utf8)&_pragma=parse_time(true)&_pragma=_busy_timeout(9999999)&mode=wal"
	} else {
		DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.Mysql.Username,
			cfg.Mysql.Pwd,
			cfg.Mysql.IP,
			cfg.Mysql.Port,
			cfg.Mysql.DbName,
		)
	}

	return DSN
}

var (
	vip *viper.Viper
)

func InitConfig(opt *CommandLineArgs) (cfg *Config, err error) {
	cfg = new(Config)
	cf := filepath.Join(opt.HomePath, opt.ConfigFile)
	if filepath.IsAbs(opt.ConfigFile) {
		cf = opt.ConfigFile
	}

	vip = viper.New()
	vip.SetConfigFile(cf)
	if err := vip.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "ReadInConfig")
	}
	if err := vip.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "config Unmarshal err")
	}
	cfg.HomePath = opt.HomePath
	cfg.IconData = opt.IconData
	cfg.CmdName = opt.CmdName

	if cfg.Translation.CfgDir != "" {

		err = filepath.Walk(cfg.Translation.CfgDir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return errors.WithStack(err)
			}

			if info.IsDir() {
				return nil
			}

			fileName := info.Name()

			ext := filepath.Ext(fileName)
			nameWithoutExt := strings.TrimSuffix(fileName, ext)

			if ext != ".json" {
				return nil
			}
			filePath := filepath.Join(cfg.Translation.CfgDir, fileName)

			fileBytes, err := ioutil.ReadFile(filePath)

			if err != nil {
				return errors.WithStack(err)
			}

			var data map[string]interface{}

			err = json.Unmarshal(fileBytes, &data)
			if err != nil {
				return errors.WithStack(err)
			}
			TranslationCfg[nameWithoutExt] = data

			return nil
		})

		if err != nil {
			return nil, errors.Wrap(err, "load i18n cfg err")
		}

	}

	log.Println("配置文件加载成功", cf)
	return cfg, nil
}

func (this *Config) SetEvKey(evKey string) *Config {
	this.EvKey = evKey
	vip.Set("evKey", evKey)
	return this
}

func (this *Config) SetWorkWechatSecert(secert string) *Config {
	this.OAuth.WorkWechat.Secert = secert
	vip.Set("oauth.workWechat.secert", secert)
	return this
}

func (this *Config) SetWorkWechatCorpid(corpid string) *Config {
	this.OAuth.WorkWechat.Corpid = corpid
	vip.Set("oauth.workWechat.corpid", corpid)
	return this
}

func (this *Config) SetWorkWechatAgentId(agentId string) *Config {
	this.OAuth.WorkWechat.AgentId = agentId
	vip.Set("oauth.workWechat.agentId", agentId)
	return this
}

func (this *Config) SetWorkWechatEnable(enable bool) *Config {
	this.OAuth.WorkWechat.Enable = enable
	vip.Set("oauth.workWechat.enable", enable)
	return this
}

func (this *Config) GetViperInstance() *viper.Viper {
	return vip
}

func (this *Config) GetLang() string {
	if this.Translation.Lang == "" {
		this.Translation.Lang = "zh-cn"
	}
	return this.Translation.Lang
}

func (this *Config) GetStorePath(tag string) string {
	dir := filepath.Join(this.StoreFileDir, tag)

	if !util.CheckFileIsExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
	return dir
}
