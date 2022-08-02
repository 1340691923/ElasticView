//路由层
package router

import (
	. "github.com/1340691923/ElasticView/api"
	. "github.com/1340691923/ElasticView/middleware"
	"github.com/1340691923/ElasticView/views"
	. "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	jsoniter "github.com/json-iterator/go"
)

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
	)
}

type routerGroupFn func(app *App)

func runRouterGroupFn(app *App, fns ...routerGroupFn) *App {
	for _, fn := range fns {
		fn(app)
	}
	return app
}
