package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/services/big_mode_service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type AiController struct {
	*BaseController
	log     *logger.AppLogger
	bigMode *big_mode_service.BigMode
}

func NewAiController(baseController *BaseController, log *logger.AppLogger, bigMode *big_mode_service.BigMode) *AiController {
	return &AiController{BaseController: baseController, log: log, bigMode: bigMode}
}

func (this *AiController) SearchBigMode(ctx *gin.Context) {
	var req dto.SearchBigModeReq
	err := ctx.BindJSON(&req)
	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}

	typ := big_mode_service.CommonSysMd
	switch req.SysType {
	case 1:
		typ = big_mode_service.DefaultSysContent
	case 2:
		typ = big_mode_service.DefaultSysMd
	}

	res, err := this.bigMode.BigModelSearch(typ, req.Content)
	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}
