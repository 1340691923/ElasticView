package router

import (
	. "github.com/1340691923/ElasticView/api"
	"github.com/1340691923/ElasticView/pkg/api_config"

	. "github.com/gofiber/fiber/v2"
)

func runSearch(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/search"
	gmUser := app.Group(AbsolutePath)
	{
		c := &SearchController{}
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "设置索引备注",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "setIndexCfg",
		}, gmUser.(*Group), true, c.SetIndexConfig)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取索引备注",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "getIndexCfg",
		}, gmUser.(*Group), true, c.GetIndexConfigs)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查询索引日志数据",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "SearchLog",
		}, gmUser.(*Group), true, c.SearchLog)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "获取映射字段别名",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "GetMappingAlias",
		}, gmUser.(*Group), true, c.GetMappingAlias)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "修改映射字段别名",
			Method:       api_config.MethodAny,
			AbsolutePath: AbsolutePath,
			RelativePath: "SetMappingAlias",
		}, gmUser.(*Group), true, c.SetMappingAlias)

	}
}
