package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/alias_service"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/index_service"
	"github.com/gin-gonic/gin"
)

// Es 索引控制器
type EsIndexController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	indexService    *index_service.IndexService
}

func NewEsIndexController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, indexService *index_service.IndexService) *EsIndexController {
	return &EsIndexController{BaseController: baseController, log: log, esClientService: esClientService, indexService: indexService}
}

// 创建索引
func (this *EsIndexController) CreateAction(ctx *gin.Context) {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.indexService.EsIndexCreate(ctx, esI, esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 删除索引
func (this *EsIndexController) DeleteAction(ctx *gin.Context) {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.indexService.EsIndexDelete(ctx, esI, esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return

}

// 获取索引配置信息
func (this *EsIndexController) GetSettingsAction(ctx *gin.Context) {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.indexService.EsIndexGetSettings(ctx, esI, esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)

}

// 获取所有的索引配置信息
func (this *EsIndexController) GetSettingsInfoAction(ctx *gin.Context) {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.indexService.EsIndexGetSettingsInfo(ctx, esI, esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}

// 获取别名
func (this *EsIndexController) GetAliasAction(ctx *gin.Context) {
	esAliasInfo := new(dto.EsAliasInfo)
	err := ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := alias_service.NewAliasService().EsIndexGetAlias(ctx, esI, esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)

}

func (this *EsIndexController) MoveAliasToIndex(ctx *gin.Context) {
	esAliasInfo := new(dto.EsAliasInfo)
	err := ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = alias_service.NewAliasService().MoveAliasToIndex(ctx, esI, esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

func (this *EsIndexController) AddAliasToIndex(ctx *gin.Context) {
	esAliasInfo := new(dto.EsAliasInfo)
	err := ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = alias_service.NewAliasService().AddAliasToIndex(ctx, esI, esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

func (this *EsIndexController) BatchAddAliasToIndex(ctx *gin.Context) {
	esAliasInfo := new(dto.EsAliasInfo)
	err := ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = alias_service.NewAliasService().BatchAddAliasToIndex(ctx, esI, esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

func (this *EsIndexController) RemoveAlias(ctx *gin.Context) {
	esAliasInfo := new(dto.EsAliasInfo)
	err := ctx.Bind(&esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = alias_service.NewAliasService().RemoveAlias(ctx, esI, esAliasInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 重建索引
func (this *EsIndexController) ReindexAction(ctx *gin.Context) {
	esReIndexInfo := new(dto.EsReIndexInfo)
	err := ctx.Bind(&esReIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esReIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.indexService.EsIndexReindex(ctx, esI, esReIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}

// 得到所有的索引名
func (this *EsIndexController) IndexNamesAction(ctx *gin.Context) {
	esConnectID := new(dto.EsConnectID)
	err := ctx.Bind(&esConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.indexService.EsIndexNames(ctx, esI)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}

// 得到所有的索引数量
func (this *EsIndexController) IndexsCountAction(ctx *gin.Context) {
	esConnectID := new(dto.EsConnectID)
	err := ctx.Bind(&esConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.indexService.EsIndexCount(ctx, esI)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}

// 获取索引的Stats
func (this *EsIndexController) StatsAction(ctx *gin.Context) {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.Bind(&esIndexInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.indexService.EsIndexStats(ctx, esI, esIndexInfo.IndexName)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, res)
}
