package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/es_settings"
	. "github.com/gofiber/fiber/v2"

	jsoniter "github.com/json-iterator/go"
	"github.com/olivere/elastic"
)

// 备份控制器
type EsBackUpController struct {
	BaseController
}

//快照仓库列表
func (this EsBackUpController) SnapshotRepositoryListAction(ctx *Ctx) error {
	esSnapshotInfo := es.EsSnapshotInfo{}
	err := ctx.BodyParser(&esSnapshotInfo)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(esSnapshotInfo.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	clusterSettings, err := es_settings.NewSettings(esClinet.(*es.EsClientV6).Client)
	if err != nil {
		return this.Error(ctx, err)
	}
	pathRepo := clusterSettings.GetPathRepo()

	if len(pathRepo) == 0 {
		return this.Error(ctx, my_error.NewError(`path.repo没有设置`, 199999))
	}

	res, err := esClinet.(*es.EsClientV6).Client.SnapshotGetRepository(esSnapshotInfo.SnapshotInfoList...).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	type tmp struct {
		Name                   string `json:"name"`
		Type                   string `json:"type"`
		Location               string `json:"location"`
		Compress               string `json:"compress"`
		MaxRestoreBytesPerSec  string `json:"max_restore_bytes_per_sec"`
		MaxSnapshotBytesPerSec string `json:"max_snapshot_bytes_per_sec"`
		ChunkSize              string `json:"chunk_size"`
		Readonly               string `json:"readonly"`
	}
	list := []tmp{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for name, settings := range res {
		var t tmp
		t.Type = settings.Type
		t.Name = name
		b, err := json.Marshal(settings.Settings)
		if err != nil {
			logs.Logger.Sugar().Errorf("err", err)
			continue
		}
		err = json.Unmarshal(b, &t)
		if err != nil {
			logs.Logger.Sugar().Errorf("err", err)
			continue
		}
		list = append(list, t)
	}

	return this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"list":     list,
		"res":      res,
		"pathRepo": pathRepo,
	})
}

//新建快照仓库
func (this EsBackUpController) SnapshotCreateRepositoryAction(ctx *Ctx) error {
	snapshotCreateRepository := es.SnapshotCreateRepository{}
	err := ctx.BodyParser(&snapshotCreateRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotCreateRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	clusterSettings, err := es_settings.NewSettings(esClinet.(*es.EsClientV6).Client)
	if err != nil {
		return this.Error(ctx, err)
	}
	pathRepo := clusterSettings.GetPathRepo()
	getAllowedUrls := clusterSettings.GetAllowedUrls()

	settings := map[string]interface{}{}

	if snapshotCreateRepository.Compress != "" {
		compress := snapshotCreateRepository.Compress
		settings["compress"] = compress
	}

	if snapshotCreateRepository.MaxRestoreBytesPerSec != "" {
		settings["max_restore_bytes_per_sec"] = snapshotCreateRepository.MaxRestoreBytesPerSec
	}

	if snapshotCreateRepository.MaxSnapshotBytesPerSec != "" {
		settings["max_snapshot_bytes_per_sec"] = snapshotCreateRepository.MaxSnapshotBytesPerSec
	}

	if snapshotCreateRepository.Readonly != "" {
		settings["readonly"] = snapshotCreateRepository.Readonly
	}

	if snapshotCreateRepository.ChunkSize != "" {
		settings["chunk_size"] = snapshotCreateRepository.ChunkSize
	}

	switch snapshotCreateRepository.Type {
	case "fs":
		if len(pathRepo) == 0 {
			return this.Error(ctx, errors.New("请先设置 path.repo"))

		}
		settings["location"] = snapshotCreateRepository.Location
	case "url":
		if len(getAllowedUrls) == 0 {
			return this.Error(ctx, errors.New("请先设置 allowed_urls"))

		}
		settings["url"] = snapshotCreateRepository.Location
	default:
		return this.Error(ctx, errors.New("无效的type"))

	}

	_, err = esClinet.(*es.EsClientV6).Client.SnapshotCreateRepository(snapshotCreateRepository.Repository).Type(snapshotCreateRepository.Type).Settings(
		settings,
	).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//清理快照仓库
func (this EsBackUpController) CleanupeRepositoryAction(ctx *Ctx) error {
	cleanupeRepository := es.CleanupeRepository{}
	err := ctx.BodyParser(&cleanupeRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(cleanupeRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}
	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/_snapshot/%s/_cleanup", cleanupeRepository.Repository),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

//删除快照仓库
func (this EsBackUpController) SnapshotDeleteRepositoryAction(ctx *Ctx) error {
	snapshotDeleteRepository := es.SnapshotDeleteRepository{}
	err := ctx.BodyParser(&snapshotDeleteRepository)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotDeleteRepository.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	_, err = esClinet.(*es.EsClientV6).Client.SnapshotDeleteRepository(snapshotDeleteRepository.Repository).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//创建快照
func (this EsBackUpController) CreateSnapshotAction(ctx *Ctx) error {
	createSnapshot := es.CreateSnapshot{}
	err := ctx.BodyParser(&createSnapshot)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(createSnapshot.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	snapshotCreateService := esClinet.(*es.EsClientV6).Client.
		SnapshotCreate(createSnapshot.RepositoryName, createSnapshot.SnapshotName)

	if createSnapshot.Wait != nil {
		snapshotCreateService.WaitForCompletion(*createSnapshot.Wait)
	}

	settings := es.Json{}

	if len(createSnapshot.IndexList) > 0 {
		settings["indices"] = strings.Join(createSnapshot.IndexList, ",")
	}

	if createSnapshot.IgnoreUnavailable != nil {
		settings["indices"] = *createSnapshot.IgnoreUnavailable
	}

	if createSnapshot.Partial != nil {
		settings["partial"] = *createSnapshot.Partial
	}
	if createSnapshot.IncludeGlobalState != nil {
		settings["include_global_state"] = *createSnapshot.IncludeGlobalState
	}

	res, err := snapshotCreateService.BodyJson(settings).Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

//快照列表
func (this EsBackUpController) SnapshotListAction(ctx *Ctx) error {
	snapshotList := es.SnapshotList{}
	err := ctx.BodyParser(&snapshotList)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotList.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	if snapshotList.Repository == "" {
		return this.Error(ctx, errors.New("请先选择快照存储库"))

	}

	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/snapshots/%s", snapshotList.Repository),
	})

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}

//删除快照
func (this EsBackUpController) SnapshotDeleteAction(ctx *Ctx) error {
	snapshotDelete := es.SnapshotDelete{}
	err := ctx.BodyParser(&snapshotDelete)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotDelete.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	_, err = esClinet.(*es.EsClientV6).Client.
		SnapshotDelete(snapshotDelete.Repository, snapshotDelete.Snapshot).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

//快照详情
func (this EsBackUpController) SnapshotDetailAction(ctx *Ctx) error {
	snapshotDetail := es.SnapshotDetail{}
	err := ctx.BodyParser(&snapshotDetail)
	if err != nil {
		return this.Error(ctx, err)
	}

	esClinet, err := es.GetEsClientV6ByID(snapshotDetail.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	res, err := esClinet.(*es.EsClientV6).Client.PerformRequest(ctx.Context(), elastic.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_snapshot/%s/%s", snapshotDetail.Repository, snapshotDetail.Snapshot),
	})
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res.Body)
}

// 将索引恢复至快照时状态
func (this EsBackUpController) SnapshotRestoreAction(ctx *Ctx) error {
	snapshotRestore := es.SnapshotRestore{}
	err := ctx.BodyParser(&snapshotRestore)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotRestore.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	snapshotRestoreService := esClinet.(*es.EsClientV6).Client.SnapshotRestore(snapshotRestore.RepositoryName, snapshotRestore.SnapshotName)

	if snapshotRestore.Wait != nil {
		snapshotRestoreService.WaitForCompletion(*snapshotRestore.Wait)
	}

	if snapshotRestore.IgnoreUnavailable != nil {
		snapshotRestoreService.IgnoreUnavailable(*snapshotRestore.IgnoreUnavailable)
	}
	if len(snapshotRestore.IndexList) > 0 {
		snapshotRestoreService.Indices(snapshotRestore.IndexList...)
	}
	if snapshotRestore.Partial != nil {
		snapshotRestoreService.Partial(*snapshotRestore.Partial)
	}
	if snapshotRestore.IncludeGlobalState != nil {
		snapshotRestoreService.IncludeGlobalState(*snapshotRestore.IncludeGlobalState)
	}
	if snapshotRestore.RenamePattern != "" {
		snapshotRestoreService.RenamePattern(snapshotRestore.RenamePattern)
	}
	if snapshotRestore.RenameReplacement != "" {
		snapshotRestoreService.RenameReplacement(snapshotRestore.RenameReplacement)
	}

	res, err := snapshotRestoreService.Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

//得到快照状态
func (this EsBackUpController) SnapshotStatusAction(ctx *Ctx) error {
	snapshotStatus := es.SnapshotStatus{}
	err := ctx.BodyParser(&snapshotStatus)
	if err != nil {
		return this.Error(ctx, err)
	}
	esClinet, err := es.GetEsClientV6ByID(snapshotStatus.EsConnect)
	if err != nil {
		return this.Error(ctx, err)
	}

	snapshotRestoreStatus := esClinet.(*es.EsClientV6).Client.SnapshotStatus().Repository(snapshotStatus.RepositoryName).Snapshot(snapshotStatus.SnapshotName)

	res, err := snapshotRestoreStatus.Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}
