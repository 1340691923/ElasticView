package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/services/ws_service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsController struct {
	log       *logger.AppLogger
	cfg       *config.Config
	orm       *orm.Gorm
	jwtSvr    *jwt_svr.Jwt
	wsService *ws_service.WsService
}

func NewWsController(log *logger.AppLogger, cfg *config.Config, orm *orm.Gorm, jwtSvr *jwt_svr.Jwt, wsService *ws_service.WsService) *WsController {
	return &WsController{log: log, cfg: cfg, orm: orm, jwtSvr: jwtSvr, wsService: wsService}
}

func (this *WsController) WsAction(ctx *gin.Context) {
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		this.log.Sugar().Errorf("err", err)
		return
	}
	defer c.Close()
	/*cliams, err := this.jwtSvr.ParseToken(util.GetToken(ctx))
	if err != nil {
		this.log.Sugar().Errorf("err", err)

		return
	}*/
	this.wsService.InitConnect(c, ctx, 1, 1)
}
