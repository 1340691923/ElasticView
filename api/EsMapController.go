package api

import (
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/service/index_service"

	. "github.com/gofiber/fiber/v2"
)

// Es 映射控制器
type EsMappingController struct {
	BaseController
}

// Es 映射列表
func (this EsMappingController) ListAction(ctx *Ctx) error {
	esMapGetProperties := new(escache.EsMapGetProperties)
	err := ctx.BodyParser(&esMapGetProperties)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esMapGetProperties.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := index_service.NewIndexService(esI).EsMappingList(ctx.Context(), esMapGetProperties)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "ver": esI.Version()})
}

// 修改映射
func (this EsMappingController) UpdateMappingAction(ctx *Ctx) error {
	updateMapping := new(escache.UpdateMapping)
	err := ctx.BodyParser(&updateMapping)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(updateMapping.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := index_service.NewIndexService(esI).UpdateMapping(ctx.Context(), updateMapping)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}
