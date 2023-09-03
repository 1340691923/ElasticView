package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_task_service"
	"github.com/gin-gonic/gin"
)

// Es 任务控制器
type EsTaskController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	taskService     *es_task_service.EsTaskService
}

func NewEsTaskController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, taskService *es_task_service.EsTaskService) *EsTaskController {
	return &EsTaskController{BaseController: baseController, log: log, esClientService: esClientService, taskService: taskService}
}

// 任务列表
func (this *EsTaskController) ListAction(ctx *gin.Context) {
	taskListReq := new(dto.TaskList)
	err := ctx.Bind(&taskListReq)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(taskListReq.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	res, err := this.taskService.TaskList(ctx, esI)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}

// 取消任务
func (this *EsTaskController) CancelAction(ctx *gin.Context) {
	cancelTask := new(dto.CancelTask)
	err := ctx.Bind(&cancelTask)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(cancelTask.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.taskService.Cancel(ctx, esI, cancelTask)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return

}
