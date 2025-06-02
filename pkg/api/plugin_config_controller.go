package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/services/plugin_config_service"
	"github.com/gin-gonic/gin"
)

type PluginConfigController struct {
	*BaseController
	pluginConfigService *plugin_config_service.PluginConfigServie
}

func NewPluginConfigController(baseController *BaseController, pluginConfigService *plugin_config_service.PluginConfigServie) *PluginConfigController {
	return &PluginConfigController{BaseController: baseController, pluginConfigService: pluginConfigService}
}

// GetPluginConfig 获取插件配置
func (this *PluginConfigController) GetPluginConfig(ctx *gin.Context) {
	var req dto.GetPluginConfigReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	config, err := this.pluginConfigService.GetPluginConfig(ctx, req.PluginID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res := dto.PluginConfigRes{
		PluginID:   config.PluginID,
		AutoUpdate: config.AutoUpdate,
	}

	this.Success(ctx, response.SearchSuccess, res)
}

// UpdatePluginConfig 更新插件配置
func (this *PluginConfigController) UpdatePluginConfig(ctx *gin.Context) {
	var req dto.UpdatePluginConfigReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.pluginConfigService.UpdatePluginConfig(ctx, req.PluginID, req.AutoUpdate)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
}
