package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// DslHistory 路由
func runDslHistory(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/dslHistory"
	dslHistory := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "清空DSL查询历史记录",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "CleanAction",
		}, dslHistory.(*Group), true, DslHistoryController{}.CleanAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看DSL查询历史记录",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "ListAction",
		}, dslHistory.(*Group), true, DslHistoryController{}.ListAction)

	}
}
