package controller

import (
	"github.com/1340691923/ElasticView/engine/es"
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"
	"github.com/gofiber/fiber/v2"
	"log"
)

// ES CRUD操作
type EsCrudController struct {
	BaseController
}

// 删除文档数据
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
	log.Println(crudFilter)

	return esService.CrudGetList(ctx, crudFilter)
}
