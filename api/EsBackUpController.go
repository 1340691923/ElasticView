package api

import (
	"github.com/1340691923/ElasticView/es_sdk/pkg/factory"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/service/cluser_settings_service"
	"github.com/1340691923/ElasticView/service/es_backup"
	. "github.com/gofiber/fiber/v2"
)

// 备份控制器
type EsBackUpController struct {
	BaseController
}

// 快照仓库列表
func (this EsBackUpController) SnapshotRepositoryListAction(ctx *Ctx) error {
	esSnapshotInfo := new(escache.EsSnapshotInfo)
	err := ctx.BodyParser(esSnapshotInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := escache.GetEsClientByID(esSnapshotInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	list, res, pathRepo, err := es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotRepositoryList(ctx.Context(), esSnapshotInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, util.Map{
		"list":     list,
		"res":      res,
		"pathRepo": pathRepo,
	})
}

// 新建快照仓库
func (this EsBackUpController) SnapshotCreateRepositoryAction(ctx *Ctx) error {
	snapshotCreateRepository := new(escache.SnapshotCreateRepository)
	err := ctx.BodyParser(snapshotCreateRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotCreateRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotCreateRepository(ctx.Context(), snapshotCreateRepository)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 清理快照仓库
func (this EsBackUpController) CleanupeRepositoryAction(ctx *Ctx) error {
	cleanupeRepository := new(escache.CleanupeRepository)
	err := ctx.BodyParser(cleanupeRepository)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := escache.GetEsClientByID(cleanupeRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).CleanUp(ctx.Context(), cleanupeRepository)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 删除快照仓库
func (this EsBackUpController) SnapshotDeleteRepositoryAction(ctx *Ctx) error {
	snapshotDeleteRepository := new(escache.SnapshotDeleteRepository)
	err := ctx.BodyParser(snapshotDeleteRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotDeleteRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotDeleteRepository(ctx.Context(), snapshotDeleteRepository)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 创建快照
func (this EsBackUpController) CreateSnapshotAction(ctx *Ctx) error {
	createSnapshot := new(escache.CreateSnapshot)
	err := ctx.BodyParser(createSnapshot)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(createSnapshot.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).CreateSnapshot(ctx.Context(), createSnapshot)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 快照列表
func (this EsBackUpController) SnapshotListAction(ctx *Ctx) error {
	snapshotList := new(escache.SnapshotList)
	err := ctx.BodyParser(snapshotList)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotList.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotList(ctx.Context(), snapshotList)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

// 删除快照
func (this EsBackUpController) SnapshotDeleteAction(ctx *Ctx) error {
	snapshotDelete := new(escache.SnapshotDelete)
	err := ctx.BodyParser(snapshotDelete)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotDelete.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotDelete(ctx.Context(), snapshotDelete)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 快照详情
func (this EsBackUpController) SnapshotDetailAction(ctx *Ctx) error {
	snapshotDetail := new(escache.SnapshotDetail)
	err := ctx.BodyParser(snapshotDetail)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := escache.GetEsClientByID(snapshotDetail.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotDetail(ctx.Context(), snapshotDetail)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}

// 将索引恢复至快照时状态
func (this EsBackUpController) SnapshotRestoreAction(ctx *Ctx) error {
	snapshotRestore := new(escache.SnapshotRestore)
	err := ctx.BodyParser(snapshotRestore)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotRestore.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	err = es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotRestore(ctx.Context(), snapshotRestore)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 得到快照状态
func (this EsBackUpController) SnapshotStatusAction(ctx *Ctx) error {
	snapshotStatus := new(escache.SnapshotStatus)
	err := ctx.BodyParser(snapshotStatus)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := escache.GetEsClientByID(snapshotStatus.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esI, err := factory.NewEsService(esConnect.ToEsSdkCfg())
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := es_backup.NewEsBackUpService(
		esI, cluser_settings_service.NewClusterSettingsService(esI),
	).SnapshotStatus(ctx.Context(), snapshotStatus)
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res)
}
