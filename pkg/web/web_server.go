package web

import (
	"context"
	"fmt"
	_ "github.com/1340691923/ElasticView/docs"
	"github.com/1340691923/ElasticView/pkg/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/middleware"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/resources/views"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"html/template"
	"io/fs"
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
	esController          *api.EsController
	wsController          *api.WsController
	indexController       *api.IndexController
	pluginController      *api.PluginController
}

func NewWebServer(engine *web_engine.WebEngine, log *logger.AppLogger, cfg *config.Config, rbac *access_control.Rbac, middleWareService *middleware.MiddleWareService, gmOperaterController *api.GmOperaterController, managerRoleController *api.ManagerRoleController, esLinkController *api.EsLinkController, managerUserController *api.ManagerUserController, esController *api.EsController, wsController *api.WsController, indexController *api.IndexController, pluginController *api.PluginController) *WebServer {
	return &WebServer{engine: engine, log: log, cfg: cfg, rbac: rbac, middleWareService: middleWareService, gmOperaterController: gmOperaterController, managerRoleController: managerRoleController, esLinkController: esLinkController, managerUserController: managerUserController, esController: esController, wsController: wsController, indexController: indexController, pluginController: pluginController}
}

type Config struct {
	Name string
}

func (this *WebServer) runRouter() {

	htmlTpl := template.Must(template.New("").ParseFS(views.IndexFileTemplate, "dist/*.html"))

	this.engine.GetGinEngine().SetHTMLTemplate(htmlTpl)

	js, _ := fs.Sub(views.JsFs, "dist/js")

	staticServer := this.engine.GetGinEngine()

	staticServer.StaticFS("/js/", http.FS(js))

	css, _ := fs.Sub(views.CssFs, "dist/css")
	staticServer.StaticFS("/css/", http.FS(css))

	img, _ := fs.Sub(views.ImgFs, "dist/img")
	staticServer.StaticFS("/img/", http.FS(img))

	staticServer.StaticFS("/favicon.ico", http.FS(views.FaviconFs))

	staticServer.GET("/", this.indexController.IndexHtml)

	this.engine.GetGinEngine().GET("/api/callback", this.indexController.CallBack)

	this.engine.GetGinEngine().Use(
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
			AllowHeaders: []string{"Accept", "Accept-Encoding", "Accept-Language", "Access-Control-Request-Headers",
				"Access-Control-Request-Method", "Connection", "Referer", "Sec-Fetch-Dest", "User-Agent",
				"Origin", "Authorization", "Content-Type", "X-Token", "x-token", "X-Version", "Current-Language"},
			ExposeHeaders:    []string{"Content-Type"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			MaxAge: 24 * time.Hour,
		}),
	)
	this.engine.GetGinEngine().GET("/api/Health", this.indexController.Health)
	this.engine.GetGinEngine().POST("/api/GetI18nCfg", this.indexController.GetI18nCfg)
	this.engine.GetGinEngine().POST("/api/GetOAuthList", this.managerUserController.GetOAuthList)

	this.engine.GetGinEngine().Any("/api/call_plugin_views/:plugin_id/*action",
		this.pluginController.CallPluginViews)

	this.engine.GetGinEngine().Use(
		this.middleWareService.CheckVersion,
	)

	this.engine.GetGinEngine().POST("/api/gm_user/login", this.managerUserController.Login)

	this.runNoVerificationRouter()

	this.engine.GetGinEngine().Any("/api/call_plugin/:plugin_id/*action",
		this.pluginController.CallPlugin)

	this.engine.GetGinEngine().GET("/swagger/*any", func(c *gin.Context) {
		ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "SWAGGER")(c)
	})
	this.engine.GetGinEngine().Use(
		this.middleWareService.JwtMiddleware,
		this.middleWareService.Rbac,
	)

	this.runOperaterLog()
	this.runManagerUser()
	this.runEsLink()
	this.runEs()
	this.runPlugins()

	this.engine.GetGinEngine().NoRoute(func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"code": 500,
			"msg":  "路由不存在",
		})
	})

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

		if this.cfg.DeBug {
			this.log.Debug("WebServer RUN ", zap.String("端口号", fmt.Sprintf(":%d", this.cfg.Port)))
		} else {
			this.log.Debug("WebServer RUN ", zap.String("端口号", fmt.Sprintf(":%d", this.cfg.Port)))
			this.log.Info("WebServer RUN ", zap.String("端口号", fmt.Sprintf(":%d", this.cfg.Port)))
		}

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
		uri := fmt.Sprintf("%s%s", "http://localhost", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v", uri))
	}
}
