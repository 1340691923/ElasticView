package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/navicat_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
)

// ES CRUD操作
type EsCrudController struct {
	*BaseController
	esClientService *es.EsClientService
	log             *logger.AppLogger
	navicatService  *navicat_service.NavicatService
}

func NewEsCrudController(baseController *BaseController, esClientService *es.EsClientService, log *logger.AppLogger, navicatService *navicat_service.NavicatService) *EsCrudController {
	return &EsCrudController{BaseController: baseController, esClientService: esClientService, log: log, navicatService: navicatService}
}

// 可视化筛选获取数据
func (this *EsCrudController) GetList(ctx *gin.Context) {
	crudFilter := new(dto.CrudFilter)
	err := ctx.Bind(crudFilter)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	navicatSvr := this.navicatService

	res, count, err := navicatSvr.CrudGetList(ctx, esI, crudFilter)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "count": count})
}

// 可视化GetDSL
func (this *EsCrudController) GetDSL(ctx *gin.Context) {
	crudFilter := new(dto.CrudFilter)
	err := ctx.Bind(crudFilter)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	navicatSvr := this.navicatService

	res, err := navicatSvr.CrudGetDSL(ctx, esI, crudFilter)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, util.Map{"list": res})
}

// 下载
func (this *EsCrudController) Download(ctx *gin.Context) {

	crudFilter := new(dto.CrudFilter)
	err := ctx.Bind(crudFilter)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	navicatSvr := this.navicatService

	downloadFileName, titleList, searchData, err := navicatSvr.CrudDownload(ctx, esI, crudFilter)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.DownloadExcel(downloadFileName, titleList, searchData, ctx, this.log)
	return
}
