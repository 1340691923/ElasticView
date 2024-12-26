package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/gin-gonic/gin"
)

type GmOperaterController struct {
	*BaseController
	log                  *logger.AppLogger
	cfg                  *config.Config
	gmOperaterLogService *gm_operater_log.GmOperaterLogService
}

func NewGmOperaterController(baseController *BaseController, log *logger.AppLogger, cfg *config.Config, gmOperaterLogService *gm_operater_log.GmOperaterLogService) *GmOperaterController {
	return &GmOperaterController{BaseController: baseController, log: log, cfg: cfg, gmOperaterLogService: gmOperaterLogService}
}

// @Summary 查看后台操作日志
// @Tags ev后台操作日志
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.GmOperaterLogList false "查询参数"
// @Success 0 {object} vo.GmOperaterLog
// @Router /api/operater_log/ListAction [post]
func (this *GmOperaterController) ListAction(ctx *gin.Context) {

	var reqData dto.GmOperaterLogList

	if err := ctx.Bind(&reqData); err != nil {
		this.Error(ctx, err)
		return
	}

	gmOperaterLogService := this.gmOperaterLogService

	list, count, err := gmOperaterLogService.List(ctx, reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": list, "count": count})
}
