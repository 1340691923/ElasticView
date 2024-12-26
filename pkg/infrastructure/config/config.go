package config

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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
	HomePath    string
	ConfigFile  string
	Profile     bool
	ProfileAddr string
	ProfilePort uint64
	Tracing     bool
	TracingFile string
	CmdName     string
	AdminPwd    string
	BuildVue    bool
}

type ConfigLog struct {
	StorageDays  int    `json:"storageDays"`
	LogDir       string `json:"logDir"`
	PluginLogDir string `json:"pluginLogDir"`
}

type ConfigSqlite struct {
	DbPath string `json:"dbPath"`
	DbName string `json:"dbName"`
}

type ConfigMysql struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	DbName   string `json:"dbName"`
}

type Config struct {
	HomePath              string
	CmdName               string
	StoreFileDir          string       `json:"storeFileDir"`
	EnableLogEs           bool         `json:"enableLogEs"`
	EnableLogEsRes        bool         `json:"enableLogEsRes"`
	RootUrl               string       `json:"rootUrl"`
	Log                   ConfigLog    `json:"log"`
	CheckForevUpdates     bool         `json:"checkForevUpdates"`
	CheckForPluginUpdates bool         `json:"checkForPluginUpdates"`
	Port                  int          `json:"port"`
	PluginRpcPort         int          `json:"pluginRpcPort"`
	DbType                string       `json:"dbType"`
	Sqlite                ConfigSqlite `json:"sqlite"`
	Mysql                 ConfigMysql  `json:"mysql"`
	Version               string       `json:"version"`
	DeBug                 bool         `json:"deBug"`
	EsPwdSecret           string       `json:"esPwdSecret"`
	EvKey                 string       `json:"evKey"`
	Plugin                Plugin       `json:"plugin"`
	WatermarkContent      string       `json:"watermarkContent"`
	Translation           Translation  `json:"translation"`
	OAuth                 OAuth        `json:"oAuth"`
	Ai                    AI           `json:"ai"`
}

const (
	EV_ROOT_URL          = "EV_ROOT_URL"
	EV_STORE_FILE_DIR    = "EV_STORE_FILE_DIR"
	EV_LOG_STORAGE_DAYS  = "EV_LOG_STORAGE_DAYS"
	EV_LOG_DIR           = "EV_LOG_DIR"
	EV_PLUGIN_LOG_DIR    = "EV_PLUGIN_LOG_DIR"
	EV_DB_TYPE           = "EV_DB_TYPE"
	EV_SQLITE_DB_PATH    = "EV_SQLITE_DB_PATH"
	EV_SQLITE_DB_NAME    = "EV_SQLITE_DB_NAME"
	EV_MYSQL_USERNAME    = "EV_MYSQL_USERNAME"
	EV_MYSQL_PWD         = "EV_MYSQL_PWD"
	EV_MYSQL_IP          = "EV_MYSQL_IP"
	EV_MYSQL_PORT        = "EV_MYSQL_PORT"
	EV_MYSQL_DBNAME      = "EV_MYSQL_DBNAME"
	EV_CONN_PWDSECRET    = "EV_CONN_PWDSECRET"
	EV_KEY               = "EV_KEY"
	EV_PLUGIN_LOAD_PATH  = "EV_PLUGIN_LOAD_PATH"
	EV_PLUGIN_STORE_PATH = "EV_PLUGIN_STORE_PATH"
	EV_WATERMARK_CONTENT = "EV_WATERMARK_CONTENT"
	EV_TRANSLATION_LANG  = "EV_TRANSLATION_LANG"
	EV_BIG_MODE_KEY      = "EV_BIG_MODE_KEY"
)

func (this *Config) LoadEnv() *Config {
	if os.Getenv(EV_ROOT_URL) != "" {
		this.RootUrl = os.Getenv(EV_ROOT_URL)
	}

	if os.Getenv(EV_STORE_FILE_DIR) != "" {
		this.StoreFileDir = os.Getenv(EV_STORE_FILE_DIR)
	}
	if os.Getenv(EV_LOG_STORAGE_DAYS) != "" {
		this.Log.StorageDays = cast.ToInt(os.Getenv(EV_LOG_STORAGE_DAYS))
	}
	if os.Getenv(EV_LOG_DIR) != "" {
		this.Log.LogDir = cast.ToString(os.Getenv(EV_LOG_DIR))
	}
	if os.Getenv(EV_PLUGIN_LOG_DIR) != "" {
		this.Log.PluginLogDir = cast.ToString(os.Getenv(EV_PLUGIN_LOG_DIR))
	}
	if os.Getenv(EV_DB_TYPE) != "" {
		this.DbType = cast.ToString(os.Getenv(EV_DB_TYPE))
	}
	if os.Getenv(EV_SQLITE_DB_PATH) != "" {
		this.Sqlite.DbPath = cast.ToString(os.Getenv(EV_SQLITE_DB_PATH))
	}
	if os.Getenv(EV_SQLITE_DB_NAME) != "" {
		this.Sqlite.DbName = cast.ToString(os.Getenv(EV_SQLITE_DB_NAME))
	}
	if os.Getenv(EV_MYSQL_USERNAME) != "" {
		this.Mysql.Username = cast.ToString(os.Getenv(EV_MYSQL_USERNAME))
	}
	if os.Getenv(EV_MYSQL_PWD) != "" {
		this.Mysql.Pwd = cast.ToString(os.Getenv(EV_MYSQL_PWD))
	}
	if os.Getenv(EV_MYSQL_IP) != "" {
		this.Mysql.IP = cast.ToString(os.Getenv(EV_MYSQL_IP))
	}
	if os.Getenv(EV_MYSQL_PORT) != "" {
		this.Mysql.Port = cast.ToString(os.Getenv(EV_MYSQL_PORT))
	}
	if os.Getenv(EV_MYSQL_DBNAME) != "" {
		this.Mysql.DbName = cast.ToString(os.Getenv(EV_MYSQL_DBNAME))
	}
	if os.Getenv(EV_CONN_PWDSECRET) != "" {
		this.EsPwdSecret = cast.ToString(os.Getenv(EV_CONN_PWDSECRET))
	}
	/*if os.Getenv(EV_KEY) != "" {
		this.EvKey = cast.ToString(os.Getenv(EV_KEY))
	}*/
	if os.Getenv(EV_PLUGIN_LOAD_PATH) != "" {
		this.Plugin.LoadPath = cast.ToString(os.Getenv(EV_PLUGIN_LOAD_PATH))
	}
	if os.Getenv(EV_PLUGIN_STORE_PATH) != "" {
		this.Plugin.StorePath = cast.ToString(os.Getenv(EV_PLUGIN_STORE_PATH))
	}
	if os.Getenv(EV_WATERMARK_CONTENT) != "" {
		this.WatermarkContent = cast.ToString(os.Getenv(EV_WATERMARK_CONTENT))
	}
	if os.Getenv(EV_TRANSLATION_LANG) != "" {
		this.Translation.Lang = cast.ToString(os.Getenv(EV_TRANSLATION_LANG))
	}
	if os.Getenv(EV_BIG_MODE_KEY) != "" {
		this.Ai.BigModeKey = cast.ToString(os.Getenv(EV_BIG_MODE_KEY))
	}

	log.Println("读取环境变量完毕", this)
	return this
}

type OAuth struct {
	WorkWechat WorkWechat `json:"workWechat"`
}

type AI struct {
	BigModeKey string `json:"bigModeKey"`
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

func (this *Config) MoveSqliteData(dir string) string {

	if dir == "" {

		dataDir := filepath.Join(util.GetCurrentDirectory(), "data")
		newPath := filepath.Join(this.GetStoreDir(), "data")

		if util.CheckFileIsExist(dataDir) && !util.CheckFileIsExist(newPath) {
			err := util.MoveDir(dataDir, newPath)
			if err == nil {
				log.Println("MoveDir success", dataDir, newPath)
			} else {
				log.Println("MoveDir err", dataDir, newPath, err)
			}
		}

		return newPath
	}
	newPath := filepath.Join(this.GetStoreDir(), dir)
	if filepath.IsAbs(dir) {
		return dir
	}

	if util.CheckFileIsExist(dir) && !util.CheckFileIsExist(newPath) {

		err := util.MoveDir(dir, newPath)
		if err == nil {
			log.Println("MoveDir success", dir, newPath)
		} else {
			log.Println("MoveDir err", dir, newPath, err)
		}
	}
	return newPath
}

func (cfg *Config) CreateDbDSN() string {
	var DSN string
	if cfg.DbType == SqliteDbTyp {
		dataDir := filepath.Join(cfg.GetStoreDir(), "data")

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

const StoreDir = "ev_store"

func (this *Config) GetStoreDir() string {
	return filepath.Join(this.HomePath, StoreDir)
}

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

	cfg = cfg.LoadEnv()

	os.MkdirAll(cfg.GetStoreDir(), os.ModePerm)

	cfg.Log.LogDir = cfg.NewStorePath(cfg.Log.LogDir)
	cfg.Log.PluginLogDir = cfg.NewStorePath(cfg.Log.PluginLogDir)
	cfg.Plugin.LoadPath = cfg.NewStorePath(cfg.Plugin.LoadPath)
	cfg.Plugin.StorePath = cfg.NewStorePath(cfg.Plugin.StorePath)
	cfg.StoreFileDir = cfg.NewStorePath(cfg.StoreFileDir)
	cfg.Sqlite.DbPath = cfg.MoveSqliteData(cfg.Sqlite.DbPath)

	log.Println("配置文件加载成功", cf)
	return cfg, nil
}

func (this *Config) NewStorePath(dir string) string {
	if dir == "" {
		return ""
	}
	newPath := filepath.Join(this.GetStoreDir(), dir)
	if filepath.IsAbs(dir) {
		return dir
	}
	if util.CheckFileIsExist(dir) && !util.CheckFileIsExist(newPath) {
		log.Println("开始迁移文件夹", dir, newPath)
		err := util.MoveDir(dir, newPath)
		if err == nil {
			log.Println("MoveDir success", dir, newPath)
		} else {
			log.Println("MoveDir err", dir, newPath, err)
		}
	}
	return newPath
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

func GetVersion() string {
	return strings.ReplaceAll(Version, "v", "")
}
