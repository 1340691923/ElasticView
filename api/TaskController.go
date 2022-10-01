package api

import (
	"github.com/1340691923/ElasticView/pkg/escache"
	es2 "github.com/1340691923/ElasticView/service/es"
	. "github.com/gofiber/fiber/v2"
)

// Es 任务控制器
type TaskController struct {
	BaseController
}

// 任务列表
func (this TaskController) ListAction(ctx *Ctx) error {
	taskListReq := new(escache.TaskList)
	err := ctx.BodyParser(&taskListReq)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(taskListReq.EsConnect)
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
	cancelTask := new(escache.CancelTask)
	err := ctx.BodyParser(&cancelTask)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(cancelTask.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)

	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.Cancel(ctx, cancelTask)

}
