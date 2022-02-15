package controller

import (
	"fmt"

	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"

	"github.com/olivere/elastic"
)

// ES 文档控制器
type EsDocController struct {
	BaseController
}

// 删除文档数据
func (this EsDocController) DeleteRowByIDAction(ctx *Ctx) error {
	esDocDeleteRowByID := es.EsDocDeleteRowByID{}
	err := ctx.BodyParser(&esDocDeleteRowByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esDocDeleteRowByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/%s/%s/%s", esDocDeleteRowByID.IndexName, esDocDeleteRowByID.Type, esDocDeleteRowByID.ID),
		Body:   nil,
	})

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

// 修改文档
func (this EsDocController) UpdateByIDAction(ctx *Ctx) error {
	esDocUpdateByID := es.EsDocUpdateByID{}
	err := ctx.BodyParser(&esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esDocUpdateByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := esClinet.(*es.EsClientV6).Client.Update().Index(esDocUpdateByID.Index).Type(esDocUpdateByID.Type).Id(esDocUpdateByID.ID).
		Doc(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

// 新增文档
func (this EsDocController) InsertAction(ctx *Ctx) error {
	esDocUpdateByID := es.EsDocUpdateByID{}
	err := ctx.BodyParser(&esDocUpdateByID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esDocUpdateByID.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := esClinet.(*es.EsClientV6).Client.Index().
		Index(esDocUpdateByID.Index).
		Type(esDocUpdateByID.Type).BodyJson(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}
