package router

import (
	. "github.com/1340691923/ElasticView/controller"
	api_config "github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES连接 路由
func runEsLink(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/es_link"
	esLink := app.Group(AbsolutePath)
	{
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "新增ES连接树信息",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "InsertAction",
		}, esLink.(*Group), true, EsLinkController{}.InsertAction)
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "删除ES连接树信息",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "DeleteAction",
		}, esLink.(*Group), true, EsLinkController{}.DeleteAction)
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改ES连接树信息",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "UpdateAction",
		}, esLink.(*Group), true, EsLinkController{}.UpdateAction)
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看ES连接树列表",
			Method:       api_config.MethodGet,
			AbsolutePath: AbsolutePath,
			RelativePath: "ListAction",
		}, esLink.(*Group), true, EsLinkController{}.ListAction)
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看ES连接配置下拉选",
			Method:       api_config.MethodGet,
			AbsolutePath: AbsolutePath,
			RelativePath: "OptAction",
		}, esLink.(*Group), false, EsLinkController{}.OptAction)
	}
}
