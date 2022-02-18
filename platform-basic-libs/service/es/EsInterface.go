package es

import (
	"errors"
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/gofiber/fiber/v2"
)

type EsInterface interface {
	SnapshotRepositoryList(ctx *fiber.Ctx, esSnapshotInfo *es.EsSnapshotInfo) (err error)
	SnapshotCreateRepository(ctx *fiber.Ctx, snapshotCreateRepository *es.SnapshotCreateRepository) (err error)
	CleanupeRepository(ctx *fiber.Ctx, repository *es.CleanupeRepository) (err error)
	SnapshotDeleteRepository(ctx *fiber.Ctx, repository *es.SnapshotDeleteRepository) (err error)
	CreateSnapshot(ctx *fiber.Ctx, snapshot *es.CreateSnapshot) (err error)
	SnapshotList(ctx *fiber.Ctx, list *es.SnapshotList) (err error)
	SnapshotDelete(ctx *fiber.Ctx, snapshotDelete *es.SnapshotDelete) (err error)
	SnapshotDetail(ctx *fiber.Ctx, detail *es.SnapshotDetail) (err error)
	SnapshotRestore(ctx *fiber.Ctx, restore *es.SnapshotRestore) (err error)
	SnapshotStatus(ctx *fiber.Ctx, status *es.SnapshotStatus) (err error)
	Cat(ctx *fiber.Ctx, rest *es.EsCat) (err error)
	RunDsl(ctx *fiber.Ctx, optimize *es.EsRest) (err error)
	Optimize(ctx *fiber.Ctx, optimize *es.EsOptimize) (err error)
	RecoverCanWrite(ctx *fiber.Ctx) (err error)
	EsDocDeleteRowByID(ctx *fiber.Ctx, id *es.EsDocDeleteRowByID) (err error)
	EsDocUpdateByID(ctx *fiber.Ctx, id *es.EsDocUpdateByID) (err error)
	EsDocInsert(ctx *fiber.Ctx, id *es.EsDocUpdateByID) (err error)
	EsIndexCreate(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsIndexDelete(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsIndexGetSettings(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsIndexGetSettingsInfo(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsIndexGetAlias(ctx *fiber.Ctx, info *es.EsAliasInfo) (err error)
	EsIndexOperateAlias(ctx *fiber.Ctx, info *es.EsAliasInfo) (err error)
	EsIndexReindex(ctx *fiber.Ctx, info *es.EsReIndexInfo) (err error)
	EsIndexIndexNames(ctx *fiber.Ctx) (err error)
	EsIndexStats(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsIndexCatStatus(ctx *fiber.Ctx, info *es.EsIndexInfo) (err error)
	EsMappingList(ctx *fiber.Ctx, properties *es.EsMapGetProperties) (err error)
	UpdateMapping(ctx *fiber.Ctx, mapping *es.UpdateMapping) (err error)
	TaskList(ctx *fiber.Ctx) (err error)
	Cancel(ctx *fiber.Ctx, task *es.CancelTask) (err error)
	CrudGetList(ctx *fiber.Ctx, task *es.CrudFilter) (err error)
}

var VerError = errors.New("ES版本暂只支持6,7")

var EsServiceMap = map[int]func(conn *es.EsConnect) (EsInterface, error){
	6: NewEsServiceV6,
	7: NewEsServiceV7,
}

func NewEsService(conn *es.EsConnect) (EsInterface, error) {
	var found bool
	var fn func(conn *es.EsConnect) (EsInterface, error)
	if fn, found = EsServiceMap[conn.Version]; !found {
		return nil, VerError
	}
	fn = EsServiceMap[conn.Version]
	return fn(conn)
}
