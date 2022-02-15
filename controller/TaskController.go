package controller

import (
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	. "github.com/gofiber/fiber/v2"

	"github.com/olivere/elastic"
)

// Es 任务控制器
type TaskController struct {
	BaseController
}

// 任务列表
func (this TaskController) ListAction(ctx *Ctx) error {
	taskListReq := es.TaskList{}
	err := ctx.BodyParser(&taskListReq)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(taskListReq.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	tasksListService := esClinet.(*es.EsClientV6).Client.TasksList().Detailed(true)

	tasksListResponse, err := tasksListService.Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	taskListRes := map[string]*elastic.TaskInfo{}

	for _, node := range tasksListResponse.Nodes {
		for taskId, taskInfo := range node.Tasks {
			taskListRes[taskId] = taskInfo
		}
	}

	return this.Success(ctx, response.SearchSuccess, taskListRes)
}

// 取消任务
func (this TaskController) CancelAction(ctx *Ctx) error {
	cancelTask := es.CancelTask{}
	err := ctx.BodyParser(&cancelTask)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(cancelTask.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.TasksCancel().TaskId(cancelTask.TaskID).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}
