package api

import (
	"github.com/1340691923/ElasticView/pkg/escache"
	es2 "github.com/1340691923/ElasticView/service/es"

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

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsMappingList(ctx, esMapGetProperties)
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
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.UpdateMapping(ctx, updateMapping)
}
