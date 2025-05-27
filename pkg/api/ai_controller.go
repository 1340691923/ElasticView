package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
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

func (this *AiController) GetAIConfig(ctx *gin.Context) {
	cfg := config.GetConfig()
	
	res := dto.AIConfigRes{
		QwenEnabled:     cfg.Ai.BigModeKey != "",
		BigModeKey:      cfg.Ai.BigModeKey,
		OpenAIEnabled:   cfg.Ai.OpenAIKey != "",
		OpenAIKey:       cfg.Ai.OpenAIKey,
		DeepSeekEnabled: cfg.Ai.DeepSeekKey != "",
		DeepSeekKey:     cfg.Ai.DeepSeekKey,
	}
	
	this.Success(ctx, response.SearchSuccess, res)
}

func (this *AiController) SaveAIConfig(ctx *gin.Context) {
	var req dto.AIConfigReq
	err := ctx.BindJSON(&req)
	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}
	
	cfg := config.GetConfig()
	
	if req.QwenEnabled {
		cfg.Ai.BigModeKey = req.BigModeKey
	} else {
		cfg.Ai.BigModeKey = ""
	}
	
	if req.OpenAIEnabled {
		cfg.Ai.OpenAIKey = req.OpenAIKey
	} else {
		cfg.Ai.OpenAIKey = ""
	}
	
	if req.DeepSeekEnabled {
		cfg.Ai.DeepSeekKey = req.DeepSeekKey
	} else {
		cfg.Ai.DeepSeekKey = ""
	}
	
	err = config.SaveConfig(cfg)
	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}
	
	this.Success(ctx, response.InsertSuccess, "AI配置保存成功")
}

func (this *AiController) TestAIConnection(ctx *gin.Context) {
	var req dto.AIConfigReq
	err := ctx.BindJSON(&req)
	if err != nil {
		this.Error(ctx, errors.WithStack(err))
		return
	}
	
	if req.QwenEnabled && req.BigModeKey != "" {
		_, err = this.bigMode.BigModelSearch(big_mode_service.CommonSysMd, "测试连接")
		if err != nil {
			this.Error(ctx, errors.WithStack(err))
			return
		}
	}
	
	this.Success(ctx, response.SearchSuccess, "AI服务连接测试成功")
}
