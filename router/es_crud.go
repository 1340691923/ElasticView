package router

import (
	. "github.com/1340691923/ElasticView/controller"
	"github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES基础操作 路由
func runEsCrud(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/es_crud"
	es := app.Group(AbsolutePath)
	{
		esCrudController := EsCrudController{}
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "Navicat（数据筛选）",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetList",
		}, es.(*Group), true, esCrudController.GetList)
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "Navicat（获取查询语句）",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetDSL",
		}, es.(*Group), true, esCrudController.GetDSL)

	}
}
