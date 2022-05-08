package api

import (
	"github.com/1340691923/ElasticView/engine/es"
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"
	. "github.com/gofiber/fiber/v2"
)

// Es 任务控制器
type TaskController struct {
	BaseController
}

// 任务列表
func (this TaskController) ListAction(ctx *Ctx) error {
	taskListReq := new(es.TaskList)
	err := ctx.BodyParser(&taskListReq)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(taskListReq.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)

	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.TaskList(ctx)
}

// 取消任务
func (this TaskController) CancelAction(ctx *Ctx) error {
	cancelTask := new(es.CancelTask)
	err := ctx.BodyParser(&cancelTask)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(cancelTask.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)

	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.Cancel(ctx, cancelTask)

}
