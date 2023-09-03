package web

import (
	"context"
	"fmt"
	_ "github.com/1340691923/ElasticView/docs"
	"github.com/1340691923/ElasticView/pkg/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/middleware"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/resources/views"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WebServer struct {
	engine                *web_engine.WebEngine
	log                   *logger.AppLogger
	cfg                   *config.Config
	rbac                  *access_control.Rbac
	middleWareService     *middleware.MiddleWareService
	gmOperaterController  *api.GmOperaterController
	managerRoleController *api.ManagerRoleController
	esLinkController      *api.EsLinkController
	managerUserController *api.ManagerUserController
	dslHistoryController  *api.DslHistoryController
	esBackUpController    *api.EsBackUpController
	esController          *api.EsController
	esCrudController      *api.EsCrudController
	esDocController       *api.EsDocController
	esIndexController     *api.EsIndexController
	esMappingController   *api.EsMappingController
	esTaskController      *api.EsTaskController
	guidController        *api.GuidController
}

func NewWebServer(engine *web_engine.WebEngine, log *logger.AppLogger, cfg *config.Config, rbac *access_control.Rbac, middleWareService *middleware.MiddleWareService, gmOperaterController *api.GmOperaterController, managerRoleController *api.ManagerRoleController, managerUserController *api.ManagerUserController, esLinkController *api.EsLinkController, dslHistoryController *api.DslHistoryController, esBackUpController *api.EsBackUpController, esController *api.EsController, esCrudController *api.EsCrudController, esDocController *api.EsDocController, esIndexController *api.EsIndexController, esMappingController *api.EsMappingController, esTaskController *api.EsTaskController, guidController *api.GuidController) *WebServer {
	return &WebServer{engine: engine, log: log, cfg: cfg, rbac: rbac, middleWareService: middleWareService, gmOperaterController: gmOperaterController, managerRoleController: managerRoleController, esLinkController: esLinkController, managerUserController: managerUserController, dslHistoryController: dslHistoryController, esBackUpController: esBackUpController, esController: esController, esCrudController: esCrudController, esDocController: esDocController, esIndexController: esIndexController, esMappingController: esMappingController, esTaskController: esTaskController, guidController: guidController}
}

func (this *WebServer) runRouter() {
	this.engine.GetGinEngine().Use(static.Serve("/", views.EmbedFolder(views.StatisFs, "dist")))

	this.engine.GetGinEngine().Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Token", "x-token"},
			ExposeHeaders:    []string{"Content-Type"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 24 * time.Hour,
		}),
	)

	this.engine.GetGinEngine().POST("/api/gm_user/login", this.managerUserController.Login)

	this.runNoVerificationRouter()

	this.engine.GetGinEngine().GET("/swagger/*any", func(c *gin.Context) {
		ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER")(c)
	})
	this.engine.GetGinEngine().Use(
		this.middleWareService.JwtMiddleware,
		this.middleWareService.Rbac,
	)

	this.runManagerUser()
	this.runOperaterLog()
	this.runDslHistory()

	this.runEs()
	this.runEsBackUp()
	this.runEsCrud()
	this.runEsDoc()
	this.runEsIndex()
	this.runEsLink()
	this.runEsMap()
	this.runEsTask()
	this.runGmGuid()

}

func (this *WebServer) Run(ctx context.Context) (err error) {

	err = this.rbac.LoadPolicy()
	if err != nil {
		return
	}
	this.runRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", this.cfg.Port),
		Handler: this.engine.GetGinEngine(),
	}

	go func() {
		this.InitOpenWinBrowser()
		this.log.Info("WebServer RUN ", zap.String("端口号", fmt.Sprintf(":%d", this.cfg.Port)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			this.log.Error("WebServer RUN err", zap.Error(err))
			panic(err)
		}
	}()

	<-ctx.Done()
	if srv == nil {
		return
	}
	err = srv.Shutdown(context.Background())
	if err != nil {
		return
	}
	return
}

func (this *WebServer) InitOpenWinBrowser() {
	if !this.cfg.DeBug {
		port := ":" + strconv.Itoa(this.cfg.Port)
		uri := fmt.Sprintf("%s%s", "http://127.0.0.1", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v", uri))
	}
}
