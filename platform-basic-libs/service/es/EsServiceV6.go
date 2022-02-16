package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/engine/es"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/es_optimize"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/es_settings"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	elasticV6 "github.com/olivere/elastic"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
	"sync"
)

type EsServiceV6 struct {
	response.Response
	request.Request
	esClient *elasticV6.Client
}

func (this EsServiceV6) CreateSnapshot(ctx *fiber.Ctx, createSnapshot *es.CreateSnapshot) (err error) {
	snapshotCreateService := this.esClient.
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

func (this EsServiceV6) SnapshotList(ctx *fiber.Ctx, snapshotList *es.SnapshotList) (err error) {
	if snapshotList.Repository == "" {
		return this.Error(ctx, errors.New("请先选择快照存储库"))
	}

	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/snapshots/%s", snapshotList.Repository),
	})

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV6) SnapshotDelete(ctx *fiber.Ctx, snapshotDelete *es.SnapshotDelete) (err error) {
	_, err = this.esClient.
		SnapshotDelete(snapshotDelete.Repository, snapshotDelete.Snapshot).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV6) SnapshotDetail(ctx *fiber.Ctx, snapshotDetail *es.SnapshotDetail) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_snapshot/%s/%s", snapshotDetail.Repository, snapshotDetail.Snapshot),
	})
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV6) SnapshotRestore(ctx *fiber.Ctx, snapshotRestore *es.SnapshotRestore) (err error) {

	snapshotRestoreService := this.esClient.SnapshotRestore(snapshotRestore.RepositoryName, snapshotRestore.SnapshotName)

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

func (this EsServiceV6) SnapshotStatus(ctx *fiber.Ctx, snapshotStatus *es.SnapshotStatus) (err error) {
	snapshotRestoreStatus := this.esClient.SnapshotStatus().Repository(snapshotStatus.RepositoryName).Snapshot(snapshotStatus.SnapshotName)

	res, err := snapshotRestoreStatus.Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

func (this EsServiceV6) Cat(ctx *fiber.Ctx, esCat *es.EsCat) (err error) {
	var data interface{}

	switch esCat.Cat {
	case "CatHealth":
		data, err = this.esClient.CatHealth().Human(true).Do(ctx.Context())
	case "CatShards":
		data, err = this.esClient.CatShards().Human(true).Do(ctx.Context())
	case "CatCount":
		data, err = this.esClient.CatCount().Human(true).Do(ctx.Context())
	case "CatAllocation":
		data, err = this.esClient.CatAllocation().Human(true).Do(ctx.Context())
	case "CatAliases":
		data, err = this.esClient.CatAliases().Human(true).Do(ctx.Context())
	case "CatIndices":
		if esCat.IndexBytesFormat != "" {
			data, err = this.esClient.CatIndices().Human(true).Bytes(esCat.IndexBytesFormat).Do(ctx.Context())
		} else {
			data, err = this.esClient.CatIndices().Human(true).Do(ctx.Context())
		}
	case "CatSegments":
		data, err = this.esClient.IndexSegments().Human(true).Do(ctx.Context())
	case "CatStats":
		data, err = this.esClient.ClusterStats().Human(true).Do(ctx.Context())
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, data)
}

func (this EsServiceV6) RunDsl(ctx *fiber.Ctx, esRest *es.EsRest) (err error) {
	esRest.Method = strings.ToUpper(esRest.Method)
	if esRest.Method == "GET" {
		c, err := jwt.ParseToken(ctx.Get("X-Token"))
		if err != nil {
			return this.Error(ctx, err)
		}

		gmDslHistoryModel := model.GmDslHistoryModel{
			Uid:    int(c.ID),
			Method: esRest.Method,
			Path:   esRest.Path,
			Body:   esRest.Body,
		}

		err = gmDslHistoryModel.Insert()

		if err != nil {
			return this.Error(ctx, err)
		}
	}

	res, err := this.esClient.PerformRequest(context.Background(), elasticV6.PerformRequestOptions{
		Method: esRest.Method,
		Path:   esRest.Path,
		Body:   esRest.Body,
	})

	if err != nil {
		return this.Error(ctx, err)
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, map[string]interface{}{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV6) Optimize(ctx *fiber.Ctx, esOptimize *es.EsOptimize) (err error) {
	optimize := es_optimize.OptimizeFactory(esOptimize.Command)

	if optimize == nil {
		return this.Error(ctx, errors.New("不支持该指令"))

	}
	if esOptimize.IndexName != "" {
		optimize.SetIndexName(esOptimize.IndexName)
	}
	err = optimize.DoV6(this.esClient)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV6) RecoverCanWrite(ctx *fiber.Ctx) (err error) {

	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "PUT",
		Path:   "/_settings",
		Body: map[string]interface{}{
			"index": map[string]interface{}{
				"blocks": map[string]interface{}{
					"read_only_allow_delete": "false",
				},
			},
		},
	})

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, map[string]interface{}{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV6) EsDocDeleteRowByID(ctx *fiber.Ctx, esDocDeleteRowByID *es.EsDocDeleteRowByID) (err error) {

	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/%s/%s/%s", esDocDeleteRowByID.IndexName, esDocDeleteRowByID.Type, esDocDeleteRowByID.ID),
		Body:   nil,
	})

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsDocUpdateByID(ctx *fiber.Ctx, esDocUpdateByID *es.EsDocUpdateByID) (err error) {
	res, err := this.esClient.Update().Index(esDocUpdateByID.Index).Type(esDocUpdateByID.Type).Id(esDocUpdateByID.ID).
		Doc(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsDocInsert(ctx *fiber.Ctx, esDocUpdateByID *es.EsDocUpdateByID) (err error) {
	res, err := this.esClient.Index().
		Index(esDocUpdateByID.Index).
		Type(esDocUpdateByID.Type).BodyJson(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexCreate(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}
	var res interface{}
	if esIndexInfo.Types == "update" {
		res, err = this.esClient.IndexPutSettings().Index(esIndexInfo.IndexName).BodyJson(esIndexInfo.Settings).Do(context.TODO())
		if err != nil {
			return this.Error(ctx, err)
		}

	} else {
		res, err = this.esClient.CreateIndex(esIndexInfo.IndexName).BodyJson(map[string]interface{}{
			"settings": esIndexInfo.Settings,
		}).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexDelete(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))

	}
	_, err = this.esClient.DeleteIndex(strings.Split(esIndexInfo.IndexName, ",")...).Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV6) EsIndexGetSettings(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	res, err := this.esClient.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res[esIndexInfo.IndexName].Settings)
}

func (this EsServiceV6) EsIndexGetSettingsInfo(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	res, err := this.esClient.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexGetAlias(ctx *fiber.Ctx, esAliasInfo *es.EsAliasInfo) (err error) {
	if esAliasInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
	}

	aliasRes, err := this.esClient.Aliases().Index(esAliasInfo.IndexName).Do(ctx.Context())

	return this.Success(ctx, response.OperateSuccess, aliasRes.Indices[esAliasInfo.IndexName].Aliases)
}

func (this EsServiceV6) EsIndexOperateAlias(ctx *fiber.Ctx, esAliasInfo *es.EsAliasInfo) (err error) {
	const Add = 1
	const Delete = 2
	const MoveToAnotherIndex = 3
	const PatchAdd = 4
	var res interface{}
	switch esAliasInfo.Types {
	case Add:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		res, err = this.esClient.Alias().Add(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case Delete:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		res, err = this.esClient.Alias().Remove(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case MoveToAnotherIndex:
		res, err = this.esClient.Alias().Action(elastic.NewAliasAddAction(esAliasInfo.AliasName).Index(esAliasInfo.NewIndexList...)).Do(ctx.Context())
	case PatchAdd:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))
		}
		wg := sync.WaitGroup{}
		NewAliasNameListLen := len(esAliasInfo.NewAliasNameList)
		if len(esAliasInfo.NewAliasNameList) > 10 {
			err = errors.New("别名列表数量不能大于10")
			break
		} else {
			wg.Add(NewAliasNameListLen)
			for _, aliasName := range esAliasInfo.NewAliasNameList {
				go func(aliasName string) {
					defer wg.Done()
					res, err = this.esClient.Alias().
						Add(esAliasInfo.IndexName, aliasName).
						Do(context.TODO())
				}(aliasName)
			}
			wg.Wait()
		}
	default:
		err = es.ReqParmasValid
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexReindex(ctx *fiber.Ctx, esReIndexInfo *es.EsReIndexInfo) (err error) {
	reindex := this.esClient.Reindex()
	urlValues := esReIndexInfo.UrlValues
	if urlValues.WaitForActiveShards != "" {
		reindex = reindex.WaitForActiveShards(urlValues.WaitForActiveShards)
	}
	if urlValues.Slices != 0 {
		reindex = reindex.Slices(urlValues.Slices)
	}
	if urlValues.Refresh != "" {
		reindex = reindex.Refresh(urlValues.Refresh)
	}
	if urlValues.Timeout != "" {
		reindex = reindex.Timeout(urlValues.Refresh)
	}
	if urlValues.RequestsPerSecond != 0 {
		reindex = reindex.RequestsPerSecond(urlValues.RequestsPerSecond)
	}
	if urlValues.WaitForCompletion != nil {
		reindex = reindex.WaitForCompletion(*urlValues.WaitForCompletion)
	}

	res, err := reindex.Body(esReIndexInfo.Body).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexIndexNames(ctx *fiber.Ctx) (err error) {
	catIndicesResponse, err := this.esClient.CatIndices().Human(true).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	indexNames := []string{}

	for _, catIndices := range catIndicesResponse {
		indexNames = append(indexNames, catIndices.Index)
	}

	return this.Success(ctx, response.SearchSuccess, indexNames)
}

func (this EsServiceV6) EsIndexStats(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(es.ParmasNullError, es.IndexNameNullError))

	}

	res, err := this.esClient.IndexStats(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) EsIndexCatStatus(ctx *fiber.Ctx, esIndexInfo *es.EsIndexInfo) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/indices/%s?h=status", esIndexInfo.IndexName),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV6) EsMappingList(ctx *fiber.Ctx, esConnect *es.EsMapGetProperties) (err error) {
	if esConnect.IndexName == "" {
		res, err := this.esClient.GetMapping().Do(context.Background())
		if err != nil {
			return this.Error(ctx, err)
		}

		return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res, "ver": 6})
	} else {
		res, err := this.esClient.GetMapping().Index(esConnect.IndexName).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}

		return this.Success(ctx, response.SearchSuccess, map[string]interface{}{"list": res, "ver": 6})
	}
}

func (this EsServiceV6) UpdateMapping(ctx *fiber.Ctx, updateMapping *es.UpdateMapping) (err error) {
	res, err := this.esClient.PutMapping().
		Index(updateMapping.IndexName).
		Type(updateMapping.TypeName).
		UpdateAllTypes(true).
		BodyJson(updateMapping.Properties).
		Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) TaskList(ctx *fiber.Ctx) (err error) {
	tasksListService := this.esClient.TasksList().Detailed(true)

	tasksListResponse, err := tasksListService.Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	taskListRes := map[string]*elasticV6.TaskInfo{}

	for _, node := range tasksListResponse.Nodes {
		for taskId, taskInfo := range node.Tasks {
			taskListRes[taskId] = taskInfo
		}
	}

	return this.Success(ctx, response.SearchSuccess, taskListRes)
}

func (this EsServiceV6) Cancel(ctx *fiber.Ctx, cancelTask *es.CancelTask) (err error) {
	res, err := this.esClient.TasksCancel().TaskId(cancelTask.TaskID).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV6) SnapshotRepositoryList(ctx *fiber.Ctx, esSnapshotInfo *es.EsSnapshotInfo) (err error) {

	clusterSettings, err := es_settings.NewSettingsByV6(this.esClient)
	if err != nil {
		return this.Error(ctx, err)
	}
	pathRepo := clusterSettings.GetPathRepo()

	if len(pathRepo) == 0 {
		return this.Error(ctx, my_error.NewError(`path.repo没有设置`, 199999))
	}

	res, err := this.esClient.SnapshotGetRepository(esSnapshotInfo.SnapshotInfoList...).Do(ctx.Context())
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

func (this EsServiceV6) SnapshotCreateRepository(ctx *fiber.Ctx, snapshotCreateRepository *es.SnapshotCreateRepository) (err error) {

	clusterSettings, err := es_settings.NewSettingsByV6(this.esClient)
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

	_, err = this.esClient.SnapshotCreateRepository(snapshotCreateRepository.Repository).Type(snapshotCreateRepository.Type).Settings(
		settings,
	).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV6) CleanupeRepository(ctx *fiber.Ctx, cleanupeRepository *es.CleanupeRepository) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV6.PerformRequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/_snapshot/%s/_cleanup", cleanupeRepository.Repository),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV6) SnapshotDeleteRepository(ctx *fiber.Ctx, repository *es.SnapshotDeleteRepository) (err error) {
	_, err = this.esClient.SnapshotDeleteRepository(repository.Repository).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func NewEsServiceV6(connect *es.EsConnect) (service EsInterface, err error) {
	esClinet, err := es.NewEsClientV6(connect)

	if err != nil {
		return nil, err
	}

	return &EsServiceV6{esClient: esClinet}, nil
}
