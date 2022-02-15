//应用启动引擎层
package application

var GlobConfig Config

//全局配置结构体
type Config struct {
	Log struct {
		StorageDays int    `json:"storageDays"`
		LogDir      string `json:"logDir"`
	} `json:"log"`
	Port  int `json:"port"`
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
	Esconfig struct{
		Addresses []string `json:"addresses"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"esconfig"`
	RedisPool struct{
		Addr string `json:"addr"`
		Passwd string `json:"passwd"`
		Db int `json:"db"`
		MaxIdle int `json:"maxIdle"`
		MaxActive int `json:"maxActive"`
	} `json:"redisPool"`
	DeBug bool `json:"deBug"`
}

