package api

import (
	"github.com/1340691923/ElasticView/platform-basic-libs/escache"
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"
	. "github.com/gofiber/fiber/v2"
)

// Es 索引控制器
type EsIndexController struct {
	BaseController
}

//创建索引
func (this EsIndexController) CreateAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexCreate(ctx, esIndexInfo)
}

// 删除索引
func (this EsIndexController) DeleteAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexDelete(ctx, esIndexInfo)

}

//获取索引配置信息
func (this EsIndexController) GetSettingsAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexGetSettings(ctx, esIndexInfo)

}

//获取所有的索引配置信息
func (this EsIndexController) GetSettingsInfoAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexGetSettingsInfo(ctx, esIndexInfo)

}

// 获取别名
func (this EsIndexController) GetAliasAction(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexGetAlias(ctx, esAliasInfo)

}

// 操作别名
func (this EsIndexController) OperateAliasAction(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexOperateAlias(ctx, esAliasInfo)
}

// 重建索引
func (this EsIndexController) ReindexAction(ctx *Ctx) error {
	esReIndexInfo := new(escache.EsReIndexInfo)
	err := ctx.BodyParser(&esReIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esReIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexReindex(ctx, esReIndexInfo)
}

// 得到所有的索引名
func (this EsIndexController) IndexNamesAction(ctx *Ctx) error {
	esConnectID := new(escache.EsConnectID)
	err := ctx.BodyParser(&esConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexIndexNames(ctx)
}

// 获取索引的Stats
func (this EsIndexController) StatsAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexCatStatus(ctx, esIndexInfo)
}

func (this EsIndexController) CatStatusAction(ctx *Ctx) error {
	esIndexInfo := new(escache.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.EsIndexCatStatus(ctx, esIndexInfo)
}
