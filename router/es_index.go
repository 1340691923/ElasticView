package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES索引 路由
func runEsIndex(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/es_index"
	esIndex := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "删除ES索引",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "DeleteAction",
		}, esIndex.(*Group), true, EsIndexController{}.DeleteAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "创建ES索引",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "CreateAction",
		}, esIndex.(*Group), true, EsIndexController{}.CreateAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取ES索引配置",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetSettingsAction",
		}, esIndex.(*Group), true, EsIndexController{}.GetSettingsAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "得到所有的ES索引名",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "IndexNamesAction",
		}, esIndex.(*Group), true, EsIndexController{}.IndexNamesAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "重建索引",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "ReindexAction",
		}, esIndex.(*Group), true, EsIndexController{}.ReindexAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取别名",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetAliasAction",
		}, esIndex.(*Group), true, EsIndexController{}.GetAliasAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "操作别名",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "OperateAliasAction",
		}, esIndex.(*Group), true, EsIndexController{}.OperateAliasAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取所有的索引配置信息",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetSettingsInfoAction",
		}, esIndex.(*Group), true, EsIndexController{}.GetSettingsInfoAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取索引的Stats",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "StatsAction",
		}, esIndex.(*Group), true, EsIndexController{}.StatsAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看索引的Stats",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "CatStatusAction",
		}, esIndex.(*Group), true, EsIndexController{}.CatStatusAction)

	}
}
