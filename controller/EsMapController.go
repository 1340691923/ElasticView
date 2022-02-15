package controller

import (
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"

	"github.com/1340691923/ElasticView/engine/es"
	. "github.com/gofiber/fiber/v2"
)

// Es 映射控制器
type EsMappingController struct {
	BaseController
}

// Es 映射列表
func (this EsMappingController) ListAction(ctx *Ctx) error {
	esMapGetProperties := new(es.EsMapGetProperties)
	err := ctx.BodyParser(&esMapGetProperties)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(esMapGetProperties.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsMappingList(ctx, esMapGetProperties)
}

// 修改映射
func (this EsMappingController) UpdateMappingAction(ctx *Ctx) error {
	updateMapping := new(es.UpdateMapping)
	err := ctx.BodyParser(&updateMapping)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(updateMapping.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.UpdateMapping(ctx, updateMapping)
}
