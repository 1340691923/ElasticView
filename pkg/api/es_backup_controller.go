package api

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/services/es_backup"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
)

// 备份控制器
type EsBackUpController struct {
	*BaseController
	log             *logger.AppLogger
	esClientService *es.EsClientService
	esBackUpService *es_backup.EsBackUpService
}

func NewEsBackUpController(baseController *BaseController, log *logger.AppLogger, esClientService *es.EsClientService, esBackUpService *es_backup.EsBackUpService) *EsBackUpController {
	return &EsBackUpController{BaseController: baseController, log: log, esClientService: esClientService, esBackUpService: esBackUpService}
}

// 快照仓库列表
func (this *EsBackUpController) SnapshotRepositoryListAction(ctx *gin.Context) {
	esSnapshotInfo := new(dto.EsSnapshotInfo)
	err := ctx.Bind(esSnapshotInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(esSnapshotInfo.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	list, res, pathRepo, err := this.esBackUpService.SnapshotRepositoryList(ctx, esI, esSnapshotInfo)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, util.Map{
		"list":     list,
		"res":      res,
		"pathRepo": pathRepo,
	})
}

// 新建快照仓库
func (this *EsBackUpController) SnapshotCreateRepositoryAction(ctx *gin.Context) {
	snapshotCreateRepository := new(dto.SnapshotCreateRepository)
	err := ctx.Bind(snapshotCreateRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotCreateRepository.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.SnapshotCreateRepository(ctx, esI, snapshotCreateRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 清理快照仓库
func (this *EsBackUpController) CleanupeRepositoryAction(ctx *gin.Context) {
	cleanupeRepository := new(dto.CleanupeRepository)
	err := ctx.Bind(cleanupeRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(cleanupeRepository.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.CleanUp(ctx, esI, cleanupeRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 删除快照仓库
func (this *EsBackUpController) SnapshotDeleteRepositoryAction(ctx *gin.Context) {
	snapshotDeleteRepository := new(dto.SnapshotDeleteRepository)
	err := ctx.Bind(snapshotDeleteRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotDeleteRepository.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.SnapshotDeleteRepository(ctx, esI, snapshotDeleteRepository)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 创建快照
func (this *EsBackUpController) CreateSnapshotAction(ctx *gin.Context) {
	createSnapshot := new(dto.CreateSnapshot)
	err := ctx.Bind(createSnapshot)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(createSnapshot.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.CreateSnapshot(ctx, esI, createSnapshot)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 快照列表
func (this *EsBackUpController) SnapshotListAction(ctx *gin.Context) {
	snapshotList := new(dto.SnapshotList)
	err := ctx.Bind(snapshotList)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotList.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.esBackUpService.SnapshotList(ctx, esI, snapshotList)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}

// 删除快照
func (this *EsBackUpController) SnapshotDeleteAction(ctx *gin.Context) {
	snapshotDelete := new(dto.SnapshotDelete)
	err := ctx.Bind(snapshotDelete)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotDelete.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.SnapshotDelete(ctx, esI, snapshotDelete)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 快照详情
func (this *EsBackUpController) SnapshotDetailAction(ctx *gin.Context) {
	snapshotDetail := new(dto.SnapshotDetail)
	err := ctx.Bind(snapshotDetail)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esConnect, err := this.esClientService.GetEsClientByID(snapshotDetail.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.esBackUpService.SnapshotDetail(ctx, esI, snapshotDetail)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}

// 将索引恢复至快照时状态
func (this *EsBackUpController) SnapshotRestoreAction(ctx *gin.Context) {
	snapshotRestore := new(dto.SnapshotRestore)
	err := ctx.Bind(snapshotRestore)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotRestore.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	err = this.esBackUpService.SnapshotRestore(ctx, esI, snapshotRestore)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 得到快照状态
func (this *EsBackUpController) SnapshotStatusAction(ctx *gin.Context) {
	snapshotStatus := new(dto.SnapshotStatus)
	err := ctx.Bind(snapshotStatus)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	esConnect, err := this.esClientService.GetEsClientByID(snapshotStatus.EsConnect)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res, err := this.esBackUpService.SnapshotStatus(ctx, esI, snapshotStatus)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, res)
}
