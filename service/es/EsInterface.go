package es

import (
	"errors"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/gofiber/fiber/v2"
)

type EsInterface interface {
	SnapshotRepositoryList(ctx *fiber.Ctx, esSnapshotInfo *escache.EsSnapshotInfo) (err error)
	SnapshotCreateRepository(ctx *fiber.Ctx, snapshotCreateRepository *escache.SnapshotCreateRepository) (err error)
	CleanupeRepository(ctx *fiber.Ctx, repository *escache.CleanupeRepository) (err error)
	SnapshotDeleteRepository(ctx *fiber.Ctx, repository *escache.SnapshotDeleteRepository) (err error)
	CreateSnapshot(ctx *fiber.Ctx, snapshot *escache.CreateSnapshot) (err error)
	SnapshotList(ctx *fiber.Ctx, list *escache.SnapshotList) (err error)
	SnapshotDelete(ctx *fiber.Ctx, snapshotDelete *escache.SnapshotDelete) (err error)
	SnapshotDetail(ctx *fiber.Ctx, detail *escache.SnapshotDetail) (err error)
	SnapshotRestore(ctx *fiber.Ctx, restore *escache.SnapshotRestore) (err error)
	SnapshotStatus(ctx *fiber.Ctx, status *escache.SnapshotStatus) (err error)
	Cat(ctx *fiber.Ctx, rest *escache.EsCat) (err error)
	RunDsl(ctx *fiber.Ctx, optimize *escache.EsRest) (err error)
	Optimize(ctx *fiber.Ctx, optimize *escache.EsOptimize) (err error)
	RecoverCanWrite(ctx *fiber.Ctx) (err error)
	EsDocDeleteRowByID(ctx *fiber.Ctx, id *escache.EsDocDeleteRowByID) (err error)
	EsDocUpdateByID(ctx *fiber.Ctx, id *escache.EsDocUpdateByID) (err error)
	EsDocInsert(ctx *fiber.Ctx, id *escache.EsDocUpdateByID) (err error)
	EsIndexCreate(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsIndexDelete(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsIndexGetSettings(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsIndexGetSettingsInfo(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsIndexGetAlias(ctx *fiber.Ctx, info *escache.EsAliasInfo) (err error)
	EsIndexOperateAlias(ctx *fiber.Ctx, info *escache.EsAliasInfo) (err error)
	EsIndexReindex(ctx *fiber.Ctx, info *escache.EsReIndexInfo) (err error)
	EsIndexIndexNames(ctx *fiber.Ctx) (err error)
	EsIndexStats(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsIndexCatStatus(ctx *fiber.Ctx, info *escache.EsIndexInfo) (err error)
	EsMappingList(ctx *fiber.Ctx, properties *escache.EsMapGetProperties) (err error)
	UpdateMapping(ctx *fiber.Ctx, mapping *escache.UpdateMapping) (err error)
	TaskList(ctx *fiber.Ctx) (err error)
	Cancel(ctx *fiber.Ctx, task *escache.CancelTask) (err error)
	CrudGetList(ctx *fiber.Ctx, filter *escache.CrudFilter) (err error)
	CrudGetDSL(ctx *fiber.Ctx, filter *escache.CrudFilter) (err error)
	CrudDownload(ctx *fiber.Ctx, filter *escache.CrudFilter) (err error)
}

var VerError = errors.New("ES版本暂只支持6,7")

var EsServiceMap = map[int]func(conn *escache.EsConnect) (EsInterface, error){
	6: NewEsServiceV6,
	7: NewEsServiceV7,
	8: NewEsServiceV8,
}

func NewEsService(conn *escache.EsConnect) (EsInterface, error) {
	var found bool
	var fn func(conn *escache.EsConnect) (EsInterface, error)
	if fn, found = EsServiceMap[conn.Version]; !found {
		return nil, VerError
	}
	fn = EsServiceMap[conn.Version]
	return fn(conn)
}
