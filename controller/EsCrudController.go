package controller

import (
	"github.com/1340691923/ElasticView/engine/es"
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"
	"github.com/gofiber/fiber/v2"
)

// ES CRUD操作
type EsCrudController struct {
	BaseController
}

// 可视化筛选获取数据
func (this EsCrudController) GetList(ctx *fiber.Ctx) error {
	crudFilter := new(es.CrudFilter)
	err := ctx.BodyParser(crudFilter)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	return esService.CrudGetList(ctx, crudFilter)
}

// 可视化GetDSL
func (this EsCrudController) GetDSL(ctx *fiber.Ctx) error {
	crudFilter := new(es.CrudFilter)
	err := ctx.BodyParser(crudFilter)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(crudFilter.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	return esService.CrudGetDSL(ctx, crudFilter)
}