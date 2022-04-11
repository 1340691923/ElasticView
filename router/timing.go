package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"

	. "github.com/gofiber/fiber/v2"
)

// 计划任务
func runTiming(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/TimingController"
	gmUser := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "计划任务列表",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "ListAction",
		}, gmUser.(*Group), true, TimingController{}.ListAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "取消计划任务",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "CancelAction",
		}, gmUser.(*Group), true, TimingController{}.CancelAction)

	}
}
