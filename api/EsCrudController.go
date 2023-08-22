package api

import (
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/service/navicat_service"
	"github.com/gofiber/fiber/v2"
)

// ES CRUD操作
type EsCrudController struct {
	BaseController
}

// 可视化筛选获取数据
func (this EsCrudController) GetList(ctx *fiber.Ctx) error {
	crudFilter := new(escache.CrudFilter)
	err := ctx.BodyParser(crudFilter)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	navicatSvr := navicat_service.NewNavicatService(esI)

	res, count, err := navicatSvr.CrudGetList(ctx.Context(), crudFilter)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "count": count})
}

// 可视化GetDSL
func (this EsCrudController) GetDSL(ctx *fiber.Ctx) error {
	crudFilter := new(escache.CrudFilter)
	err := ctx.BodyParser(crudFilter)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	navicatSvr := navicat_service.NewNavicatService(esI)

	res, err := navicatSvr.CrudGetDSL(ctx.Context(), crudFilter)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res})
}

// 下载
func (this EsCrudController) Download(ctx *fiber.Ctx) error {

	crudFilter := new(escache.CrudFilter)
	err := ctx.BodyParser(crudFilter)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	navicatSvr := navicat_service.NewNavicatService(esI)

	downloadFileName, titleList, searchData, err := navicatSvr.CrudDownload(ctx.Context(), crudFilter)

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.DownloadExcel(downloadFileName, titleList, searchData, ctx)
}
