package main

import (
	"flag"
	"github.com/1340691923/ElasticView/application"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/router"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"go.uber.org/zap"
)

var (
	appName        string
	configFileDir  string
	configFileName string
	configFileExt  string
)

func init() {
	flag.StringVar(&appName, "appName", "ElasticView", "应用名")
	flag.StringVar(&configFileDir, "configFileDir", "config", "配置文件夹名")
	flag.StringVar(&configFileName, "configFileName", "config", "配置文件名")
	flag.StringVar(&configFileExt, "configFileExt", "json", "配置文件后缀")
	flag.Parse()
}

// By 肖文龙
func main() {
	app := application.NewApp(
		application.WithAppName(appName),
		application.WithConfigFileDir(configFileDir),
		application.WithConfigFileName(configFileName),
		application.WithConfigFileExt(configFileExt),
		application.RegisterInitFnObserver(application.InitLogs),
		application.RegisterInitFnObserver(application.InitMysql),
		application.RegisterInitFnObserver(application.InitTask),
		application.RegisterInitFnObserver(application.InitRbac),
		application.RegisterInitFnObserver(application.InitOpenWinBrowser),
	)

	err := app.InitConfig().NotifyInitFnObservers().Error()

	if err != nil {
		logs.Logger.Error("ElasticView 初始化失败", zap.String("err.Error()", err.Error()))
		panic(err)
	}

	port := ":" + strconv.Itoa(application.GlobConfig.Port)
	appServer := router.Init()

	go func() {
		if err := appServer.Listen(port); err != nil {
			logs.Logger.Error("ElasticView http服务启动失败:", zap.String("err.Error()", err.Error()))
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	logs.Logger.Info("ElasticView http服务停止中...")
	// 这里进行任务释放操作
	logs.Logger.Info("ElasticView http服务停止成功...")
}
