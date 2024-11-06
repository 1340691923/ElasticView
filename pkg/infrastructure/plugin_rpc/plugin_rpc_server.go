package plugin_rpc

import (
	"context"
	"fmt"
	_ "github.com/1340691923/ElasticView/docs"
	"github.com/1340691923/ElasticView/pkg/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/middleware"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type PluginRpcServer struct {
	engine               *web_engine.WebEngine
	log                  *logger.AppLogger
	cfg                  *config.Config
	middleWareService    *middleware.MiddleWareService
	pluginUtilController *api.PluginUtilController
}

func NewPluginRpcServer(log *logger.AppLogger, cfg *config.Config, middleWareService *middleware.MiddleWareService, pluginUtilController *api.PluginUtilController) *PluginRpcServer {
	return &PluginRpcServer{engine: web_engine.NewWebEngine(), log: log, cfg: cfg, middleWareService: middleWareService, pluginUtilController: pluginUtilController}
}

type Config struct {
	Name string
}

func (this *PluginRpcServer) runRouter() {

	this.runPluginUtil()

	this.engine.GetGinEngine().NoRoute(func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"code": 500,
			"msg":  "路由不存在",
		})
	})

}

func (this *PluginRpcServer) Run(ctx context.Context) (err error) {

	this.runRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", this.cfg.PluginRpcPort),
		Handler: this.engine.GetGinEngine(),
	}

	go func() {
		this.log.Info("PluginRpcServer RUN ", zap.String("端口号", fmt.Sprintf(":%d", this.cfg.PluginRpcPort)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			this.log.Error("PluginRpcServer RUN err", zap.Error(err))
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
