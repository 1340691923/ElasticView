package es_backup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/pkg/vo"
	"github.com/1340691923/ElasticView/service/cluser_settings_service"
	"net/http"
	"strings"
)

type EsBackUpService struct {
	esClient               pkg.EsI
	clusterSettingsService *cluser_settings_service.ClusterSettingsService
}

func NewEsBackUpService(esClient pkg.EsI, clusterSvr *cluser_settings_service.ClusterSettingsService) *EsBackUpService {
	return &EsBackUpService{esClient: esClient, clusterSettingsService: clusterSvr}
}

func (this *EsBackUpService) SnapshotRepositoryList(ctx context.Context, esSnapshotInfo *escache.EsSnapshotInfo) (list []vo.Snashot, snapshotGetRepository map[string]vo.SnapshotRepository, pathRepo []interface{}, err error) {
	clusterSettings, err := this.clusterSettingsService.GetSettings(ctx)
	if err != nil {
		return
	}
	pathRepo = this.clusterSettingsService.GetPathRepo(clusterSettings)

	if len(pathRepo) == 0 {
		err = my_error.NewError(`path.repo没有设置`, 199999)
		return
	}
	res, err := this.esClient.SnapshotGetRepository(ctx, esSnapshotInfo.SnapshotInfoList)

	if err != nil {
		return
	}

	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	snapshotGetRepository = map[string]vo.SnapshotRepository{}

	err = json.Unmarshal(res.ResByte(), &snapshotGetRepository)
	if err != nil {
		return
	}

	for name, settings := range snapshotGetRepository {
		var t vo.Snashot
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

	return

}

func (this *EsBackUpService) SnapshotCreateRepository(ctx context.Context, snapshotCreateRepository *escache.SnapshotCreateRepository) (err error) {

	clusterSettings, err := this.clusterSettingsService.GetSettings(ctx)
	if err != nil {
		return
	}
	pathRepo := this.clusterSettingsService.GetPathRepo(clusterSettings)
	getAllowedUrls := this.clusterSettingsService.GetAllowedUrls(clusterSettings)

	settings := util.Map{}

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
			err = errors.New("请先设置 path.repo")
			return
		}
		settings["location"] = snapshotCreateRepository.Location
	case "url":
		if len(getAllowedUrls) == 0 {
			err = errors.New("请先设置 allowed_urls")
			return
		}
		settings["url"] = snapshotCreateRepository.Location
	default:
		err = errors.New("无效的type")
		return
	}

	body := map[string]interface{}{
		"type":     snapshotCreateRepository.Type,
		"settings": settings,
	}

	res, err := this.esClient.SnapshotCreateRepository(ctx, snapshotCreateRepository.Repository, body)
	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	return
}

func (this *EsBackUpService) CleanUp(ctx context.Context, cleanupeRepository *escache.CleanupeRepository) (err error) {

	//todo...  Invalid snapshot name [_cleanup], must not start with '_'
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("/_snapshot/%s/_cleanup", cleanupeRepository.Repository), nil)

	if err != nil {
		return
	}

	res, err := this.esClient.PerformRequest(ctx, req)
	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	return
}

func (this *EsBackUpService) SnapshotDeleteRepository(ctx context.Context, repository *escache.SnapshotDeleteRepository) (err error) {

	res, err := this.esClient.SnapshotDeleteRepository(ctx, []string{repository.Repository})
	if err != nil {
		return
	}

	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	return
}

func (this *EsBackUpService) CreateSnapshot(ctx context.Context, createSnapshot *escache.CreateSnapshot) (err error) {

	settings := escache.Json{}

	if len(createSnapshot.IndexList) > 0 {
		settings["indices"] = strings.Join(createSnapshot.IndexList, ",")
	}

	if createSnapshot.Partial != nil {
		settings["partial"] = *createSnapshot.Partial
	}
	if createSnapshot.IncludeGlobalState != nil {
		settings["include_global_state"] = *createSnapshot.IncludeGlobalState
	}

	res, err := this.esClient.SnapshotCreate(
		ctx,
		createSnapshot.RepositoryName,
		createSnapshot.SnapshotName,
		createSnapshot.Wait,
		settings,
	)

	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	return
}

func (this *EsBackUpService) SnapshotList(ctx context.Context, snapshotList *escache.SnapshotList) (res []vo.Snapshot, err error) {

	if snapshotList.Repository == "" {
		err = errors.New("请先选择快照存储库")
		return
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/_cat/snapshots/%s?format=json", snapshotList.Repository), nil)

	if err != nil {
		return
	}

	performRequestRes, err := this.esClient.PerformRequest(ctx, req)

	if err != nil {
		return
	}

	if performRequestRes.StatusErr() != nil {
		err = performRequestRes.StatusErr()
		return
	}

	err = json.Unmarshal(performRequestRes.ResByte(), &res)

	if err != nil {
		return
	}

	return
}

func (this *EsBackUpService) SnapshotDetail(ctx context.Context, snapshotDetail *escache.SnapshotDetail) (res *vo.SnapshotDetail, err error) {

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("/_snapshot/%s/%s", snapshotDetail.Repository, snapshotDetail.Snapshot),
		nil)
	if err != nil {
		return
	}

	esRes, err := this.esClient.PerformRequest(ctx, req)
	if err != nil {
		return
	}

	if esRes.StatusErr() != nil {
		err = esRes.StatusErr()
		return
	}

	res = new(vo.SnapshotDetail)

	err = json.Unmarshal(esRes.ResByte(), res)

	if err != nil {
		return
	}

	return
}

func (this *EsBackUpService) SnapshotStatus(ctx context.Context, snapshotStatus *escache.SnapshotStatus) (res *vo.SnapshotStatus, err error) {
	esRes, err := this.esClient.SnapshotStatus(
		ctx,
		snapshotStatus.RepositoryName,
		[]string{snapshotStatus.SnapshotName},
		nil,
	)
	if err != nil {
		return
	}

	if esRes.StatusErr() != nil {
		err = esRes.StatusErr()
		return
	}

	res = new(vo.SnapshotStatus)
	err = json.Unmarshal(esRes.ResByte(), &res)
	if err != nil {
		return
	}

	return
}

func (this *EsBackUpService) SnapshotDelete(ctx context.Context, snapshotDelete *escache.SnapshotDelete) (err error) {
	esRes, err := this.esClient.
		SnapshotDelete(ctx, snapshotDelete.Repository, snapshotDelete.Snapshot)
	if err != nil {
		return
	}

	if esRes.StatusErr() != nil {
		err = esRes.StatusErr()
		return
	}

	return
}

func (this *EsBackUpService) SnapshotRestore(ctx context.Context, snapshotRestore *escache.SnapshotRestore) (err error) {

	body := map[string]interface{}{}

	if snapshotRestore.IncludeGlobalState != nil {
		body["include_global_state"] = *snapshotRestore.IncludeGlobalState
	}
	if snapshotRestore.Partial != nil {
		body["partial"] = *snapshotRestore.Partial
	}

	if len(snapshotRestore.IndexList) > 0 {
		body["indices"] = strings.Join(snapshotRestore.IndexList, ",")
	}
	if snapshotRestore.RenamePattern != "" {
		body["rename_pattern"] = snapshotRestore.RenamePattern
	}
	if snapshotRestore.RenameReplacement != "" {
		body["rename_replacement"] = snapshotRestore.RenameReplacement
	}

	esRes, err := this.esClient.RestoreSnapshot(
		ctx,
		snapshotRestore.RepositoryName,
		snapshotRestore.SnapshotName,
		snapshotRestore.Wait,
		body,
	)
	if err != nil {
		return
	}

	if esRes.StatusErr() != nil {
		err = esRes.StatusErr()
		return
	}

	return
}
