package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_doc_service"
	"github.com/gin-gonic/gin"
)

// ES 文档控制器
type EsDocController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	esDocService    *es_doc_service.EsDocService
}

func NewEsDocController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, esDocService *es_doc_service.EsDocService) *EsDocController {
	return &EsDocController{BaseController: baseController, log: log, esClientService: esClientService, esDocService: esDocService}
}

// 删除文档数据
func (this *EsDocController) DeleteRowByIDAction(ctx *gin.Context) {
	esDocDeleteRowByID := new(dto.EsDocDeleteRowByID)
	err := ctx.Bind(esDocDeleteRowByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esDocDeleteRowByID.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.esDocService.DeleteRowByIDAction(ctx, esI, esDocDeleteRowByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 修改文档
func (this *EsDocController) UpdateByIDAction(ctx *gin.Context) {
	esDocUpdateByID := new(dto.EsDocUpdateByID)
	err := ctx.Bind(esDocUpdateByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esDocUpdateByID.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.esDocService.EsDocUpdateByID(ctx, esI, esDocUpdateByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 新增文档
func (this *EsDocController) InsertAction(ctx *gin.Context) {
	esDocUpdateByID := new(dto.EsDocUpdateByID)
	err := ctx.Bind(esDocUpdateByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(esDocUpdateByID.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.esDocService.EsDocInsert(ctx, esI, esDocUpdateByID)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, res)

}
