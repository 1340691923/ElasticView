package controller

import (
	"context"

	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"
)

// Es 映射控制器
type EsMappingController struct {
	BaseController
}

// Es 映射列表
func (this EsMappingController) ListAction(ctx *Ctx) error {
	esConnect := es.EsMapGetProperties{}
	err := ctx.BodyParser(&esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esConnect.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	if esConnect.IndexName == "" {
		res, err := esClinet.GetMappings()
		if err != nil {
			return this.Error(ctx, err)
		}
		return this.Success(ctx, response.SearchSuccess, res)
	} else {
		res, err := esClinet.(*es.EsClientV6).Client.GetMapping().Index(esConnect.IndexName).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}
		return this.Success(ctx, response.SearchSuccess, res)
	}

}

/*func (this EsMappingController) GetPropertiesAction(ctx *fiber.Ctx) {
	esConnect := es.EsMapGetProperties{}
	err := ctx.BodyParser(&esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esConnect.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.GetMapping().Index("")
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}*/

// 修改映射
func (this EsMappingController) UpdateMappingAction(ctx *Ctx) error {
	updateMapping := es.UpdateMapping{}
	err := ctx.BodyParser(&updateMapping)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(updateMapping.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.PutMapping().
		Index(updateMapping.IndexName).
		Type(updateMapping.TypeName).
		UpdateAllTypes(true).
		BodyJson(updateMapping.Properties).
		Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}
