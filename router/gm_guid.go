package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES 新手引导 路由
func runGmGuid(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/gm_guid"
	guid := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "完成新手引导",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "Finish",
		}, guid.(*Group), false, GuidController{}.Finish)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "判断是否完成新手引导",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "IsFinish",
		}, guid.(*Group), false, GuidController{}.Finish)
	}
}
