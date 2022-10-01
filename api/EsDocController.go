package api

import (
	"github.com/1340691923/ElasticView/pkg/escache"
	es2 "github.com/1340691923/ElasticView/service/es"

	. "github.com/gofiber/fiber/v2"
)

// ES 文档控制器
type EsDocController struct {
	BaseController
}

// 删除文档数据
func (this EsDocController) DeleteRowByIDAction(ctx *Ctx) error {
	esDocDeleteRowByID := new(escache.EsDocDeleteRowByID)
	err := ctx.BodyParser(esDocDeleteRowByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esDocDeleteRowByID.EsConnect)
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
	esDocUpdateByID := new(escache.EsDocUpdateByID)
	err := ctx.BodyParser(esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esDocUpdateByID.EsConnect)
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
	esDocUpdateByID := new(escache.EsDocUpdateByID)
	err := ctx.BodyParser(esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esDocUpdateByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsDocInsert(ctx, esDocUpdateByID)

}
