package es

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/escache"
	"github.com/1340691923/ElasticView/pkg/jwt"
	"github.com/1340691923/ElasticView/pkg/my_error"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/service/es/es6_utils"
	"github.com/1340691923/ElasticView/service/es/es7_utils"
	"github.com/1340691923/ElasticView/service/es_optimize"
	"github.com/1340691923/ElasticView/service/es_settings"
	"github.com/gofiber/fiber/v2"
	elasticV7 "github.com/olivere/elastic/v7"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type EsServiceV7 struct {
	response.Response
	request.Request
	esClient *elasticV7.Client
}

func (this EsServiceV7) CrudGetList(ctx *fiber.Ctx, crudFilter *escache.CrudFilter) (err error) {
	q, err := es6_utils.GetWhereSql(crudFilter.Relation)
	if err != nil {
		return this.Error(ctx, err)
	}
	search := this.esClient.Search(crudFilter.IndexName)
	q2 := search.Query(q)
	for _, tmp := range crudFilter.SortList {
		switch tmp.SortRule {
		case "desc":
			q2 = q2.Sort(tmp.Col, false)
		case "asc":
			q2 = q2.Sort(tmp.Col, true)
		}
	}

	res, err := q2.From(int(db.CreatePage(crudFilter.Page, crudFilter.Limit))).Size(crudFilter.Limit).Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "count": res.Hits.TotalHits.Value})
}

func (this EsServiceV7) CreateSnapshot(ctx *fiber.Ctx, createSnapshot *escache.CreateSnapshot) (err error) {
	snapshotCreateService := this.esClient.
		SnapshotCreate(createSnapshot.RepositoryName, createSnapshot.SnapshotName)

	if createSnapshot.Wait != nil {
		snapshotCreateService.WaitForCompletion(*createSnapshot.Wait)
	}

	settings := escache.Json{}

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

func (this EsServiceV7) SnapshotList(ctx *fiber.Ctx, snapshotList *escache.SnapshotList) (err error) {
	if snapshotList.Repository == "" {
		return this.Error(ctx, errors.New("请先选择快照存储库"))
	}

	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/snapshots/%s", snapshotList.Repository),
	})

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV7) SnapshotDelete(ctx *fiber.Ctx, snapshotDelete *escache.SnapshotDelete) (err error) {
	_, err = this.esClient.
		SnapshotDelete(snapshotDelete.Repository, snapshotDelete.Snapshot).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV7) SnapshotDetail(ctx *fiber.Ctx, snapshotDetail *escache.SnapshotDetail) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_snapshot/%s/%s", snapshotDetail.Repository, snapshotDetail.Snapshot),
	})
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV7) SnapshotRestore(ctx *fiber.Ctx, snapshotRestore *escache.SnapshotRestore) (err error) {

	snapshotRestoreService := this.esClient.SnapshotRestore(snapshotRestore.RepositoryName, snapshotRestore.SnapshotName)

	if snapshotRestore.Wait != nil {
		snapshotRestoreService = snapshotRestoreService.WaitForCompletion(*snapshotRestore.Wait)
	}

	if snapshotRestore.IgnoreUnavailable != nil {
		snapshotRestoreService = snapshotRestoreService.IgnoreUnavailable(*snapshotRestore.IgnoreUnavailable)
	}
	if len(snapshotRestore.IndexList) > 0 {
		snapshotRestoreService = snapshotRestoreService.Indices(snapshotRestore.IndexList...)
	}
	if snapshotRestore.Partial != nil {
		snapshotRestoreService = snapshotRestoreService.Partial(*snapshotRestore.Partial)
	}
	if snapshotRestore.IncludeGlobalState != nil {
		snapshotRestoreService = snapshotRestoreService.IncludeGlobalState(*snapshotRestore.IncludeGlobalState)
	}
	if snapshotRestore.RenamePattern != "" {
		snapshotRestoreService = snapshotRestoreService.RenamePattern(snapshotRestore.RenamePattern)
	}
	if snapshotRestore.RenameReplacement != "" {
		snapshotRestoreService = snapshotRestoreService.RenameReplacement(snapshotRestore.RenameReplacement)
	}

	res, err := snapshotRestoreService.Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) SnapshotStatus(ctx *fiber.Ctx, snapshotStatus *escache.SnapshotStatus) (err error) {
	snapshotRestoreStatus := this.esClient.SnapshotStatus().Repository(snapshotStatus.RepositoryName).Snapshot(snapshotStatus.SnapshotName)

	res, err := snapshotRestoreStatus.Do(ctx.Context())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, res)
}

func (this EsServiceV7) Cat(ctx *fiber.Ctx, esCat *escache.EsCat) (err error) {
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
			data, err = this.esClient.CatIndices().Sort("store.size:desc").Human(true).Bytes(esCat.IndexBytesFormat).Do(ctx.Context())
		} else {
			data, err = this.esClient.CatIndices().Sort("store.size:desc").Human(true).Do(ctx.Context())
		}
	case "CatSegments":
		data, err = this.esClient.IndexSegments().Human(true).Do(ctx.Context())
	case "CatStats":
		data, err = this.esClient.ClusterStats().Human(true).Do(ctx.Context())
	case "Node":
		parmas := url.Values{}
		parmas.Set("h", "ip,name,heap.percent,heap.current,heap.max,ram.percent,ram.current,ram.max,node.role,master,cpu,load_1m,load_5m,load_15m,disk.used_percent,disk.used,disk.total")
		var res *elasticV7.Response
		res, err = this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
			Method: "GET",
			Params: parmas,
			Path:   "/_cat/nodes",
		})
		if err != nil {
			return this.Error(ctx, err)
		}
		data = res.Body
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, data)
}

func (this EsServiceV7) EsIndexCount(ctx *fiber.Ctx) (err error) {
	catIndicesResponse, err := this.esClient.CatIndices().Columns("status").Human(true).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	count := len(catIndicesResponse)
	return this.Success(ctx, response.SearchSuccess, count)
}

func (this EsServiceV7) RunDsl(ctx *fiber.Ctx, esRest *escache.EsRest) (err error) {
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

	var performRequestOptions elasticV7.PerformRequestOptions

	if len(esRest.Path) > 0 {

		if esRest.Path[0:1] != "/" {
			esRest.Path = "/" + esRest.Path
		}

		u, err := url.Parse(esRest.Path)

		if err != nil {
			return this.Error(ctx, err)
		}
		path := strings.Split(esRest.Path, "?")[0]
		if len(strings.Split(esRest.Path, "/")) == 2 || strings.Contains(esRest.Path, "/_cat") {

			performRequestOptions = elasticV7.PerformRequestOptions{
				Method: esRest.Method,
				Path:   path,
				Body:   nil,
			}
			performRequestOptions.Params = u.Query()
		} else {
			performRequestOptions = elasticV7.PerformRequestOptions{
				Method: esRest.Method,
				Path:   path,
				Body:   esRest.Body,
			}
			performRequestOptions.Params = u.Query()
		}
	}

	res, err := this.esClient.PerformRequest(context.Background(), performRequestOptions)

	if err != nil {
		return this.Error(ctx, err)
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, util.Map{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV7) Optimize(ctx *fiber.Ctx, esOptimize *escache.EsOptimize) (err error) {
	optimize := es_optimize.OptimizeFactory(esOptimize.Command)

	if optimize == nil {
		return this.Error(ctx, errors.New("不支持该指令"))
	}
	optimize.CleanIndexName()
	if esOptimize.IndexName != "" {
		optimize.SetIndexName(esOptimize.IndexName)
	}
	err = optimize.DoV7(this.esClient)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV7) RecoverCanWrite(ctx *fiber.Ctx) (err error) {

	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
		Method: "PUT",
		Path:   "/_settings",
		Body: util.Map{
			"index": util.Map{
				"blocks": util.Map{
					"read_only_allow_delete": "false",
				},
			},
		},
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return this.Output(ctx, util.Map{
			"code": res.StatusCode,
			"msg":  fmt.Sprintf("请求异常! 错误码 :" + strconv.Itoa(res.StatusCode)),
			"data": res.Body,
		})
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV7) EsDocDeleteRowByID(ctx *fiber.Ctx, esDocDeleteRowByID *escache.EsDocDeleteRowByID) (err error) {

	res, err := this.esClient.Delete().Index(esDocDeleteRowByID.IndexName).Id(esDocDeleteRowByID.ID).Do(context.Background())

	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsDocUpdateByID(ctx *fiber.Ctx, esDocUpdateByID *escache.EsDocUpdateByID) (err error) {
	res, err := this.esClient.Update().Index(esDocUpdateByID.Index).Type(esDocUpdateByID.Type).Id(esDocUpdateByID.ID).
		Doc(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsDocInsert(ctx *fiber.Ctx, esDocUpdateByID *escache.EsDocUpdateByID) (err error) {
	res, err := this.esClient.Index().
		Index(esDocUpdateByID.Index).
		Type(esDocUpdateByID.Type).BodyJson(esDocUpdateByID.JSON).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsIndexCreate(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
	}
	var res interface{}
	if esIndexInfo.Types == "update" {
		res, err = this.esClient.IndexPutSettings().Index(esIndexInfo.IndexName).BodyJson(esIndexInfo.Settings).Do(context.TODO())
		if err != nil {
			return this.Error(ctx, err)
		}

	} else {
		res, err = this.esClient.CreateIndex(esIndexInfo.IndexName).BodyJson(util.Map{
			"settings": esIndexInfo.Settings,
		}).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsIndexDelete(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))

	}
	_, err = this.esClient.DeleteIndex(strings.Split(esIndexInfo.IndexName, ",")...).Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

func (this EsServiceV7) EsIndexGetSettings(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
	}

	res, err := this.esClient.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res[esIndexInfo.IndexName].Settings)
}

func (this EsServiceV7) EsIndexGetSettingsInfo(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
	}

	res, err := this.esClient.IndexGetSettings(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsIndexGetAlias(ctx *fiber.Ctx, esAliasInfo *escache.EsAliasInfo) (err error) {
	if esAliasInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
	}

	aliasRes, err := this.esClient.Aliases().Index(esAliasInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, aliasRes.Indices[esAliasInfo.IndexName].Aliases)
}

func (this EsServiceV7) EsIndexOperateAlias(ctx *fiber.Ctx, esAliasInfo *escache.EsAliasInfo) (err error) {
	const Add = 1
	const Delete = 2
	const MoveToAnotherIndex = 3
	const PatchAdd = 4
	var res interface{}
	switch esAliasInfo.Types {
	case Add:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
		}
		res, err = this.esClient.Alias().Add(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case Delete:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
		}
		res, err = this.esClient.Alias().Remove(esAliasInfo.IndexName, esAliasInfo.AliasName).Do(ctx.Context())
	case MoveToAnotherIndex:
		res, err = this.esClient.Alias().Action(elasticV7.NewAliasAddAction(esAliasInfo.AliasName).Index(esAliasInfo.NewIndexList...)).Do(ctx.Context())
	case PatchAdd:
		if esAliasInfo.IndexName == "" {
			return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))
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
		err = escache.ReqParmasValid
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsIndexReindex(ctx *fiber.Ctx, esReIndexInfo *escache.EsReIndexInfo) (err error) {
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

func (this EsServiceV7) EsIndexIndexNames(ctx *fiber.Ctx) (err error) {
	catIndicesResponse, err := this.esClient.CatIndices().Columns("index").Human(true).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	indexNames := []string{}
	for _, catIndices := range catIndicesResponse {
		indexNames = append(indexNames, catIndices.Index)
	}

	return this.Success(ctx, response.SearchSuccess, indexNames)
}

func (this EsServiceV7) EsIndexStats(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	if esIndexInfo.IndexName == "" {
		return this.Error(ctx, my_error.NewBusiness(escache.ParmasNullError, escache.IndexNameNullError))

	}

	res, err := this.esClient.IndexStats(esIndexInfo.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) EsIndexCatStatus(ctx *fiber.Ctx, esIndexInfo *escache.EsIndexInfo) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/_cat/indices/%s?h=status", esIndexInfo.IndexName),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, res.Body)
}

func (this EsServiceV7) EsMappingList(ctx *fiber.Ctx, esConnect *escache.EsMapGetProperties) (err error) {
	if esConnect.IndexName == "" {
		res, err := this.esClient.GetMapping().Do(context.Background())
		if err != nil {
			return this.Error(ctx, err)
		}

		return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "ver": 7})
	} else {
		res, err := this.esClient.GetMapping().Index(esConnect.IndexName).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}

		return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "ver": 7})
	}
}

func (this EsServiceV7) UpdateMapping(ctx *fiber.Ctx, updateMapping *escache.UpdateMapping) (err error) {
	res, err := this.esClient.PutMapping().
		Index(updateMapping.IndexName).
		BodyJson(updateMapping.Properties).
		Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) TaskList(ctx *fiber.Ctx) (err error) {
	tasksListService := this.esClient.TasksList().Detailed(true)

	tasksListResponse, err := tasksListService.Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	taskListRes := map[string]*elasticV7.TaskInfo{}

	for _, node := range tasksListResponse.Nodes {
		for taskId, taskInfo := range node.Tasks {
			taskListRes[taskId] = taskInfo
		}
	}

	return this.Success(ctx, response.SearchSuccess, taskListRes)
}

func (this EsServiceV7) Cancel(ctx *fiber.Ctx, cancelTask *escache.CancelTask) (err error) {
	res, err := this.esClient.TasksCancel().TaskId(cancelTask.TaskID).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, res)
}

func (this EsServiceV7) SnapshotRepositoryList(ctx *fiber.Ctx, esSnapshotInfo *escache.EsSnapshotInfo) (err error) {

	clusterSettings, err := es_settings.NewSettingsByV7(this.esClient)
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

	return this.Success(ctx, response.SearchSuccess, util.Map{
		"list":     list,
		"res":      res,
		"pathRepo": pathRepo,
	})
}

func (this EsServiceV7) SnapshotCreateRepository(ctx *fiber.Ctx, snapshotCreateRepository *escache.SnapshotCreateRepository) (err error) {

	clusterSettings, err := es_settings.NewSettingsByV7(this.esClient)
	if err != nil {
		return this.Error(ctx, err)
	}
	pathRepo := clusterSettings.GetPathRepo()
	getAllowedUrls := clusterSettings.GetAllowedUrls()

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

func (this EsServiceV7) CleanupeRepository(ctx *fiber.Ctx, cleanupeRepository *escache.CleanupeRepository) (err error) {
	res, err := this.esClient.PerformRequest(ctx.Context(), elasticV7.PerformRequestOptions{
		Method: "POST",
		Path:   fmt.Sprintf("/_snapshot/%s/_cleanup", cleanupeRepository.Repository),
	})
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, res.Body)
}

func (this EsServiceV7) SnapshotDeleteRepository(ctx *fiber.Ctx, repository *escache.SnapshotDeleteRepository) (err error) {
	_, err = this.esClient.SnapshotDeleteRepository(repository.Repository).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, nil)
}

func NewEsServiceV7(connect *escache.EsConnect) (service EsInterface, err error) {
	esClinet, err := escache.NewEsClientV7(connect)

	if err != nil {
		return nil, err
	}

	return &EsServiceV7{esClient: esClinet}, nil
}

func (this EsServiceV7) CrudGetDSL(ctx *fiber.Ctx, crudFilter *escache.CrudFilter) (err error) {
	q, err := es7_utils.GetWhereSql(crudFilter.Relation)
	if err != nil {
		return this.Error(ctx, err)
	}

	search := elasticV7.NewSearchSource()

	q2 := search.Query(q)
	for _, tmp := range crudFilter.SortList {
		switch tmp.SortRule {
		case "desc":
			q2 = q2.Sort(tmp.Col, false)
		case "asc":
			q2 = q2.Sort(tmp.Col, true)
		}
	}

	res, err := q2.From(int(db.CreatePage(crudFilter.Page, crudFilter.Limit))).Size(crudFilter.Limit).Source()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res})
}

func (this EsServiceV7) CrudDownload(ctx *fiber.Ctx, filter *escache.CrudFilter) (err error) {

	fields, err := this.esClient.GetMapping().Index(filter.IndexName).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}
	fieldsArr := []string{"_index", "_type", "_id"}
	data, ok := fields[filter.IndexName].(map[string]interface{})
	if !ok {
		return this.Error(ctx, errors.New("该索引没有映射结构"))
	}
	mappings, ok := data["mappings"].(map[string]interface{})
	if !ok {
		return this.Error(ctx, errors.New("该索引没有映射结构"))
	}
	properties, ok := mappings["properties"].(map[string]interface{})
	if !ok {
		return this.Error(ctx, errors.New("该索引没有映射结构"))
	}
	propertiesArr := []string{}
	for key := range properties {
		propertiesArr = append(propertiesArr, key)
	}
	sort.Strings(propertiesArr)
	fieldsArr = append(fieldsArr, propertiesArr...)
	q, err := es7_utils.GetWhereSql(filter.Relation)
	if err != nil {
		return this.Error(ctx, err)
	}
	search := this.esClient.Search(filter.IndexName)
	res, err := search.Query(q).Sort("_id", false).Size(8000).Do(ctx.Context())
	if err != nil {
		return this.Error(ctx, err)
	}

	lastIdArr := res.Hits.Hits[len(res.Hits.Hits)-1].Sort

	llist := [][]string{}

	flushHitsDataFn := func(hits []*elasticV7.SearchHit) {
		for _, data := range hits {
			list := []string{}
			list = append(list, data.Index, "_doc", data.Id)
			m := map[string]interface{}{}

			json.Unmarshal(data.Source, &m)

			for _, field := range fieldsArr {
				if field == "_index" || field == "_type" || field == "_id" {
					continue
				}
				if value, ok := m[field]; ok {
					list = append(list, util.ToExcelData(value))
				} else {
					list = append(list, "")
				}
			}

			llist = append(llist, list)
		}
	}

	flushHitsDataFn(res.Hits.Hits)
	haveData := true
	for haveData {
		search := this.esClient.Search(filter.IndexName)
		res, err := search.Query(q).Sort("_id", false).Size(8000).SearchAfter(lastIdArr...).Do(ctx.Context())
		if err != nil {
			return this.Error(ctx, err)
		}
		if len(res.Hits.Hits) == 0 {
			break
		}

		lastIdArr = res.Hits.Hits[len(res.Hits.Hits)-1].Sort
		flushHitsDataFn(res.Hits.Hits)
	}

	return this.DownloadExcel(
		"test",
		fieldsArr,
		llist, ctx)

}

func (this EsServiceV7) SearchLog(ctx *fiber.Ctx, req *escache.SearchlogReq) (err error) {
	var search *elasticV7.SearchService

	if len(req.IndexNames) == 0 {
		return this.Error(ctx, errors.New("索引不能为空"))
	}

	for _, indexName := range req.IndexNames {
		if strings.TrimSpace(indexName) == "" {
			return this.Error(ctx, errors.New("索引不能为空"))
		}
	}

	if req.Mode == 0 {
		search = this.esClient.Search(req.IndexNames...)
	} else {
		search = this.esClient.Search("*" + req.IndexNames[0] + "*")
	}

	source := []string{}

	searchCfgs := []model.SearchConfig{}
	builder := db.SqlBuilder.
		Select("*").
		From("search_index_config").
		Where(db.Eq{"es_connect": req.EsConnect}).
		Where(db.Eq{"index_name": req.IndexNames})
	sql, args, err := builder.ToSql()

	if err != nil {
		return
	}
	err = db.Sqlx.Select(&searchCfgs, sql, args...)

	if util.FilterMysqlNilErr(err) {
		return
	}
	outputCol := map[string]struct{}{}
	for _, v := range searchCfgs {
		arr := strings.Split(v.OutputCols, ",")
		for _, col := range arr {
			outputCol[col] = struct{}{}
		}
	}
	for k := range outputCol {
		source = append(source, k)
	}

	search = search.Source(elasticV7.NewFetchSourceContext(true).Include(source...))

	for _, v := range req.SearchLogFilter {
		search = search.Query(elasticV7.NewWildcardQuery(v.SearchCol, "*"+v.SearchText+"*"))
	}
	search = search.From(int(db.CreatePage(req.Page, req.Limit))).Size(req.Limit)

	res, err := search.
		Do(context.Background())
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, util.Map{"list": res, "count": res.Hits.TotalHits})
}
