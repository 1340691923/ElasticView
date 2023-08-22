package api

import (
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/service/es_doc_service"
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = es_doc_service.NewEsDocService(esI).DeleteRowByIDAction(ctx.Context(), esDocDeleteRowByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
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
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = es_doc_service.NewEsDocService(esI).EsDocUpdateByID(ctx.Context(), esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
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
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := es_doc_service.NewEsDocService(esI).EsDocInsert(ctx.Context(), esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)

}
