package router

import (
	. "github.com/1340691923/ElasticView/api"
	"github.com/1340691923/ElasticView/pkg/api_config"
	. "github.com/gofiber/fiber/v2"
)

// ES 任务 路由
func runEsTask(app *App) {
	apiRouterConfig := api_config.NewApiRouterConfig()
	const AbsolutePath = "/api/es_task"
	task := app.Group(AbsolutePath)
	{
		taskController := &TaskController{}
		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "查看任务列表",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "ListAction",
		}, task.(*Group), true, taskController.ListAction)

		apiRouterConfig.MountApi(api_config.MountApiBasePramas{
			Remark:       "取消任务",
			Method:       api_config.MethodPost,
			AbsolutePath: AbsolutePath,
			RelativePath: "CancelAction",
		}, task.(*Group), true, taskController.CancelAction)
	}
}
