// 路由层
package router

import (
	"fmt"
	. "github.com/1340691923/ElasticView/api"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/1340691923/ElasticView/pkg/engine/config"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	. "github.com/1340691923/ElasticView/pkg/middleware"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/resources/views"
	. "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"log"
	"strconv"
)

func init() {
	core.Register(core.LastLevel, "启动gofiber服务", run)
}

func run() (fn func(), err error) {
	fn = func() {}
	port := ":" + strconv.Itoa(config.GlobConfig.Port)
	appServer := Init()
	go func() {
		if err := appServer.Listen(port); err != nil {
			logs.Logger.Error("ElasticView http服务启动失败:", zap.String("err.Error()", err.Error()))
			log.Panic(err)
		}
		InitOpenWinBrowser()
	}()

	return fn, err
}

func InitOpenWinBrowser() {
	if !config.GlobConfig.DeBug {
		port := ":" + strconv.Itoa(config.GlobConfig.Port)
		uri := fmt.Sprintf("%s%s", "http://127.0.0.1", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v", uri))
	}
}

func Init() *App {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	app := New(Config{
		AppName:     "ElasticView",
		JSONDecoder: json.Unmarshal,
		JSONEncoder: json.Marshal,
	})
	app.Use(cors.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use("/", filesystem.New(filesystem.Config{
		Root: views.GetFileSystem(),
	}))

	app.All("/api/gm_user/login", UserController{}.Login)
	//routerWebsocket(app)
	app.Use(Timer)
	app.Use(JwtMiddleware)
	app.Use(Rbac)

	return runRouterGroupFn(
		app,
		runGmUser,
		runGmGuid,
		runEsLink,
		runEs,
		runEsMap,
		runEsIndex,
		runDslHistory,
		runEsTask,
		runEsBackUp,
		runEsDoc,
		runEsCrud,
		runDatax,
		runSearch,
	)
}

type routerGroupFn func(app *App)

func runRouterGroupFn(app *App, fns ...routerGroupFn) *App {
	for _, fn := range fns {
		fn(app)
	}
	return app
}
