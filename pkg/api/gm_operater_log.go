package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/gin-gonic/gin"
)

type GmOperaterController struct {
	*BaseController
	log                  *logger.AppLogger
	cfg                  *config.Config
	sqlx                 *sqlstore.SqlStore
	gmOperaterLogService *gm_operater_log.GmOperaterLogService
}

func NewGmOperaterController(baseController *BaseController, log *logger.AppLogger, gmOperaterLogService *gm_operater_log.GmOperaterLogService, cfg *config.Config, sqlx *sqlstore.SqlStore) *GmOperaterController {
	return &GmOperaterController{BaseController: baseController, log: log, gmOperaterLogService: gmOperaterLogService, cfg: cfg, sqlx: sqlx}
}

// 查看后台操作日志
func (this *GmOperaterController) ListAction(ctx *gin.Context) {

	var reqData dto.GmOperaterLogList

	if err := ctx.Bind(&reqData); err != nil {
		this.Error(ctx, err)
		return
	}

	gmOperaterLogService := this.gmOperaterLogService

	list, count, err := gmOperaterLogService.List(reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": list, "count": count})
}
