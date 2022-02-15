package controller

import (
	"github.com/1340691923/ElasticView/engine/es"
	es2 "github.com/1340691923/ElasticView/platform-basic-libs/service/es"
	. "github.com/gofiber/fiber/v2"
)

// 备份控制器
type EsBackUpController struct {
	BaseController
}

//快照仓库列表
func (this EsBackUpController) SnapshotRepositoryListAction(ctx *Ctx) error {
	esSnapshotInfo := new(es.EsSnapshotInfo)
	err := ctx.BodyParser(esSnapshotInfo)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := es.GetEsClientByID(esSnapshotInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotRepositoryList(ctx, esSnapshotInfo)
}

//新建快照仓库
func (this EsBackUpController) SnapshotCreateRepositoryAction(ctx *Ctx) error {
	snapshotCreateRepository := new(es.SnapshotCreateRepository)
	err := ctx.BodyParser(snapshotCreateRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotCreateRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotCreateRepository(ctx, snapshotCreateRepository)
}

//清理快照仓库
func (this EsBackUpController) CleanupeRepositoryAction(ctx *Ctx) error {
	cleanupeRepository := new(es.CleanupeRepository)
	err := ctx.BodyParser(cleanupeRepository)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := es.GetEsClientByID(cleanupeRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.CleanupeRepository(ctx, cleanupeRepository)
}

//删除快照仓库
func (this EsBackUpController) SnapshotDeleteRepositoryAction(ctx *Ctx) error {
	snapshotDeleteRepository := new(es.SnapshotDeleteRepository)
	err := ctx.BodyParser(snapshotDeleteRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotDeleteRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	return esService.SnapshotDeleteRepository(ctx, snapshotDeleteRepository)
}

//创建快照
func (this EsBackUpController) CreateSnapshotAction(ctx *Ctx) error {
	createSnapshot := new(es.CreateSnapshot)
	err := ctx.BodyParser(createSnapshot)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(createSnapshot.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	return esService.CreateSnapshot(ctx, createSnapshot)
}

//快照列表
func (this EsBackUpController) SnapshotListAction(ctx *Ctx) error {
	snapshotList := new(es.SnapshotList)
	err := ctx.BodyParser(snapshotList)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotList.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotList(ctx, snapshotList)
}

//删除快照
func (this EsBackUpController) SnapshotDeleteAction(ctx *Ctx) error {
	snapshotDelete := new(es.SnapshotDelete)
	err := ctx.BodyParser(snapshotDelete)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotDelete.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotDelete(ctx, snapshotDelete)
}

//快照详情
func (this EsBackUpController) SnapshotDetailAction(ctx *Ctx) error {
	snapshotDetail := new(es.SnapshotDetail)
	err := ctx.BodyParser(snapshotDetail)
	if err != nil {
		return this.Error(ctx, err)
	}

	esConnect, err := es.GetEsClientByID(snapshotDetail.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotDetail(ctx, snapshotDetail)
}

// 将索引恢复至快照时状态
func (this EsBackUpController) SnapshotRestoreAction(ctx *Ctx) error {
	snapshotRestore := new(es.SnapshotRestore)
	err := ctx.BodyParser(snapshotRestore)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotRestore.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotRestore(ctx, snapshotRestore)
}

//得到快照状态
func (this EsBackUpController) SnapshotStatusAction(ctx *Ctx) error {
	snapshotStatus := new(es.SnapshotStatus)
	err := ctx.BodyParser(snapshotStatus)
	if err != nil {
		return this.Error(ctx, err)
	}
	esConnect, err := es.GetEsClientByID(snapshotStatus.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	esService, err := es2.NewEsService(esConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	return esService.SnapshotStatus(ctx, snapshotStatus)
}
