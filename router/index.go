//路由层
package router

import (
	. "github.com/1340691923/ElasticView/controller"
	. "github.com/1340691923/ElasticView/middleware"
	"github.com/1340691923/ElasticView/views"
	. "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Init() *App {
	app := New()
	app.Use(cors.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use("/", filesystem.New(filesystem.Config{
		Root: views.GetFileSystem(),
	}))

	app.All("/api/gm_user/login", UserController{}.Login)

	app.Use(Timer)
	app.Use(JwtMiddleware)
	app.Use(Rbac)
	/*app.Use(OperaterLog)*/
	runGmUser(app)
	runGmGuid(app)
	runEsLink(app)
	runEs(app)
	runEsMap(app)
	runEsIndex(app)
	runDslHistory(app)
	runEsTask(app)
	runEsBackUp(app)
	runEsDoc(app)
	return app
}
