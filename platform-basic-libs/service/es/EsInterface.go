package es

import . "github.com/gofiber/fiber/v2"

type EsBackUpInterface interface {
	SnapshotRepositoryList(ctx *Ctx)
	SnapshotCreateRepository(ctx *Ctx)
	CleanupeRepository(ctx *Ctx)
	SnapshotDeleteRepository(ctx *Ctx)
	CreateSnapshot(ctx *Ctx)
	SnapshotList(ctx *Ctx)
	SnapshotDelete(ctx *Ctx)
	SnapshotDetail(ctx *Ctx)
	SnapshotRestore(ctx *Ctx)
	SnapshotStatus(ctx *Ctx)
}

type EsInterface interface {
	Ping(ctx *Ctx)
	Cat(ctx *Ctx)
	RunDsl(ctx *Ctx)
	SqlToDsl(ctx *Ctx)
	Optimize(ctx *Ctx)
	RecoverCanWrite(ctx *Ctx)
}

type EsDocInterface interface {
	DeleteRowByID(ctx *Ctx)
	UpdateByID(ctx *Ctx)
	Insert(ctx *Ctx)
}

type EsIndexInterface interface {
	Create(ctx *Ctx)
	Delete(ctx *Ctx)
	GetSettings(ctx *Ctx)
	GetSettingsInfo(ctx *Ctx)
	GetAlias(ctx *Ctx)
	OperateAlias(ctx *Ctx)
	Reindex(ctx *Ctx)
	IndexNames(ctx *Ctx)
	Stats(ctx *Ctx)
	CatStatus(ctx *Ctx)
}

type EsLinkInterface interface {
	List(ctx *Ctx)
	Opt(ctx *Ctx)
	Insert(ctx *Ctx)
	Update(ctx *Ctx)
	Delete(ctx *Ctx)
}

type EsMappingInterface interface {
	List(ctx *Ctx)
	UpdateMapping(ctx *Ctx)
}

type GuidInterface interface {
	Finish(ctx *Ctx)
	IsFinish(ctx *Ctx)
}

type TaskInterface interface {
	List(ctx *Ctx)
	Cancel(ctx *Ctx)
}
