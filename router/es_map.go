package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES mapping 路由
func runEsMap(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/es_map"
	esMap := app.Group("/api/es_map")
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看mapping列表",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "ListAction",
		}, esMap.(*Group), true, EsMappingController{}.ListAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改mapping",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "UpdateMappingAction",
		}, esMap.(*Group), true, EsMappingController{}.UpdateMappingAction)

	}
}
