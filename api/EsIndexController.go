package api

import (
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/service/alias_service"
	"github.com/1340691923/ElasticView/service/index_service"
	. "github.com/gofiber/fiber/v2"
)

// Es 索引控制器
type EsIndexController struct {
	BaseController
}

// 创建索引
func (this EsIndexController) CreateAction(ctx *Ctx) error {
	esIndexInfo := new(dto.EsIndexInfo)
	err := ctx.BodyParser(&esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := escache.GetEsClientByID(esIndexInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = index_service.NewIndexService(esI).EsIndexCreate(ctx.Context(), esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = index_service.NewIndexService(esI).EsIndexDelete(ctx.Context(), esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)

}

// 获取索引配置信息
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := index_service.NewIndexService(esI).EsIndexGetSettings(ctx.Context(), esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)

}

// 获取所有的索引配置信息
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := index_service.NewIndexService(esI).EsIndexGetSettingsInfo(ctx.Context(), esIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := alias_service.NewAliasService(esI).EsIndexGetAlias(ctx.Context(), esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)

}

func (this EsIndexController) MoveAliasToIndex(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = alias_service.NewAliasService(esI).MoveAliasToIndex(ctx.Context(), esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsIndexController) AddAliasToIndex(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = alias_service.NewAliasService(esI).AddAliasToIndex(ctx.Context(), esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsIndexController) BatchAddAliasToIndex(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = alias_service.NewAliasService(esI).BatchAddAliasToIndex(ctx.Context(), esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsIndexController) RemoveAlias(ctx *Ctx) error {
	esAliasInfo := new(escache.EsAliasInfo)
	err := ctx.BodyParser(&esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esAliasInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = alias_service.NewAliasService(esI).RemoveAlias(ctx.Context(), esAliasInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
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
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := index_service.NewIndexService(esI).EsIndexReindex(ctx.Context(), esReIndexInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
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
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := index_service.NewIndexService(esI).EsIndexNames(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

// 得到所有的索引数量
func (this EsIndexController) IndexsCountAction(ctx *Ctx) error {
	esConnectID := new(escache.EsConnectID)
	err := ctx.BodyParser(&esConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(esConnectID.EsConnectID)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := index_service.NewIndexService(esI).EsIndexCount(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := index_service.NewIndexService(esI).EsIndexStats(ctx.Context(), esIndexInfo.IndexName)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}
