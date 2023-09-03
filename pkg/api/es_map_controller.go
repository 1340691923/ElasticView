package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/index_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
)

// Es 映射控制器
type EsMappingController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	indexService    *index_service.IndexService
}

func NewEsMappingController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, indexService *index_service.IndexService) *EsMappingController {
	return &EsMappingController{BaseController: baseController, log: log, esClientService: esClientService, indexService: indexService}
}

// Es 映射列表
func (this *EsMappingController) ListAction(ctx *gin.Context) {
	esMapGetProperties := new(dto.EsMapGetProperties)
	err := ctx.Bind(&esMapGetProperties)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esMapGetProperties.EsConnectID)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.indexService.EsMappingList(ctx, esI, esMapGetProperties)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "ver": esI.Version()})
}

// 修改映射
func (this *EsMappingController) UpdateMappingAction(ctx *gin.Context) {
	updateMapping := new(dto.UpdateMapping)
	err := ctx.Bind(&updateMapping)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(updateMapping.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.indexService.UpdateMapping(ctx, esI, updateMapping)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}
