package controller

import (
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"

	"github.com/1340691923/ElasticView/engine/es"
	. "github.com/gofiber/fiber/v2"
)

// ES 文档控制器
type EsDocController struct {
	BaseController
}

// 删除文档数据
func (this EsDocController) DeleteRowByIDAction(ctx *Ctx) error {
	esDocDeleteRowByID := new(es.EsDocDeleteRowByID)
	err := ctx.BodyParser(esDocDeleteRowByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(esDocDeleteRowByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsDocDeleteRowByID(ctx, esDocDeleteRowByID)
}

// 修改文档
func (this EsDocController) UpdateByIDAction(ctx *Ctx) error {
	esDocUpdateByID := new(es.EsDocUpdateByID)
	err := ctx.BodyParser(esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(esDocUpdateByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsDocUpdateByID(ctx, esDocUpdateByID)
}

// 新增文档
func (this EsDocController) InsertAction(ctx *Ctx) error {
	esDocUpdateByID := new(es.EsDocUpdateByID)
	err := ctx.BodyParser(esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(esDocUpdateByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsDocInsert(ctx, esDocUpdateByID)

}
