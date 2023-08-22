package api

import (
	"fmt"
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/service/task_service"
	. "github.com/gofiber/fiber/v2"
	"log"
	"runtime"
)

// Es 任务控制器
type TaskController struct {
	BaseController
}

// 任务列表
func (this TaskController) ListAction(ctx *Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			//打印调用栈信息
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])
			logs.Logger.Sugar().Errorf("panic stack info %s", stackInfo)
			logs.Logger.Sugar().Errorf("--->HaveLoginUserSign Error:", r)
			log.Println(stackInfo)
		}
	}()
	taskListReq := new(escache.TaskList)
	err := ctx.BodyParser(&taskListReq)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(taskListReq.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := task_service.NewTaskService(esI).TaskList(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
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

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}

	err = task_service.NewTaskService(esI).Cancel(ctx.Context(), cancelTask)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)

}
