package v7

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/base"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/es_log"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	logger2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	elasticV7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"net/http"
)

type EsClient7 struct {
	client *elasticV7.Client
	base.BaseDatasource
}

func NewEsClient7(cfg *proto2.Config) (pkg.ClientInterface, error) {
	ds, ok := cache.GetDataSourceCache(cfg.ConnectId)
	if !ok {
		obj := &EsClient7{}
		esCfg, err := cfg.ConvertV7(
			es_log.NewLogger(cfg.Cfg.EnableLogEs, true, cfg.Cfg.EnableLogEsRes, logger2.EsReqLog.Named("elasticsearch7.x")))
		if err != nil {
			err = errors.WithStack(err)
			return nil, err
		}
		obj.init(esCfg)
		cache.SaveDataSourceCache(cfg.ConnectId, obj)
		return obj, nil
	}

	return ds, nil
}

func (this *EsClient7) init(config elasticV7.Config) (err error) {

	this.client, err = elasticV7.NewClient(config)
	if err != nil {
		fmt.Printf("Error creating the client: %s\n", err)
		return
	}
	return
}

func (this *EsClient7) EsVersion() (version int, err error) {
	return 7, nil
}

func (this *EsClient7) Ping(
	ctx context.Context,
) (
	res *proto.Response,
	err error,
) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return
	}

	res, err = this.EsPerformRequest(ctx, req)
	if err != nil {
		return
	}
	if res.StatusCode() == 401 {
		err = errors.New("ES地址OK，但是密码验证失败")
		return
	}

	return
}

func (this *EsClient7) EsGetMapping(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesGetMappingRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSnapshotCreate(
	ctx context.Context,
	repository string,
	snapshot string,
	waitForCompletion *bool,
	reqJson proto.Json,
) (
	res *proto.Response, err error) {
	snapshotCreateService := esapi.SnapshotCreateRequest{
		Repository:        repository,
		Snapshot:          snapshot,
		Body:              esutil.NewJSONReader(reqJson),
		WaitForCompletion: waitForCompletion,
	}

	var httpRes *esapi.Response
	httpRes, err = snapshotCreateService.Do(ctx, this.client)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsPerformRequest(
	ctx context.Context,
	req *http.Request,
) (
	res *proto.Response, err error) {
	if req != nil {
		req = req.WithContext(ctx)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpRes, err := this.client.Perform(req)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)
	return
}

func (this *EsClient7) EsSnapshotDelete(
	ctx context.Context,
	repository string,
	snapshot string,
) (
	res *proto.Response,
	err error) {
	req := esapi.SnapshotDeleteRequest{
		Repository: repository,
		Snapshot:   snapshot,
	}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsRestoreSnapshot(
	ctx context.Context,
	repository string,
	snapshot string,
	waitForCompletion *bool,
	reqJson proto.Json,
) (
	res *proto.Response,
	err error,
) {
	request := esapi.SnapshotRestoreRequest{
		Snapshot:          snapshot,
		Repository:        repository,
		Body:              esutil.NewJSONReader(reqJson),
		WaitForCompletion: waitForCompletion,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsRefresh(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesRefreshRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsOpen(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesOpenRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsFlush(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesFlushRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesClearCache(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesClearCacheRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesClose(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesCloseRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesForcemerge(
	ctx context.Context,
	indexNames []string,
	maxNumSegments *int,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesForcemergeRequest{Index: indexNames, MaxNumSegments: maxNumSegments}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsDeleteByQuery(
	ctx context.Context,
	indexNames []string,
	documents []string,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.DeleteByQueryRequest{Index: indexNames, DocumentType: []string{}, Body: esutil.NewJSONReader(body)}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSnapshotStatus(
	ctx context.Context,
	repository string,
	snapshot []string,
	ignoreUnavailable *bool,
) (
	res *proto.Response,
	err error,
) {
	request := esapi.SnapshotStatusRequest{
		Repository:        repository,
		Snapshot:          snapshot,
		IgnoreUnavailable: ignoreUnavailable,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSnapshotGetRepository(
	ctx context.Context,
	repository []string,
) (
	res *proto.Response,
	err error,
) {
	request := esapi.SnapshotGetRepositoryRequest{
		Repository: repository,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSnapshotCreateRepository(
	ctx context.Context,
	repository string,
	reqJson proto.Json,
) (
	res *proto.Response,
	err error,
) {
	request := esapi.SnapshotCreateRepositoryRequest{
		Repository: repository,
		Body:       esutil.NewJSONReader(reqJson),
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSnapshotDeleteRepository(
	ctx context.Context,
	repository []string,
) (
	res *proto.Response,
	err error,
) {
	request := esapi.SnapshotDeleteRepositoryRequest{
		Repository: repository,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsGetIndices(
	ctx context.Context,
	catRequest proto.CatIndicesRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatIndicesRequest{
		Index:                   catRequest.Index,
		Bytes:                   catRequest.Bytes,
		Format:                  catRequest.Format,
		H:                       catRequest.H,
		Health:                  catRequest.Health,
		Help:                    catRequest.Help,
		IncludeUnloadedSegments: catRequest.IncludeUnloadedSegments,
		Local:                   catRequest.Local,
		MasterTimeout:           catRequest.MasterTimeout,
		Pri:                     catRequest.Pri,
		S:                       catRequest.S,
		V:                       catRequest.V,
		Pretty:                  catRequest.Pretty,
		Human:                   catRequest.Human,
		ErrorTrace:              catRequest.ErrorTrace,
		FilterPath:              catRequest.FilterPath,
		Header:                  catRequest.Header,
	}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatHealth(
	ctx context.Context,
	catRequest proto.CatHealthRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatHealthRequest{
		Format:     catRequest.Format,
		H:          catRequest.H,
		Help:       catRequest.Help,
		S:          catRequest.S,
		Ts:         catRequest.Ts,
		V:          catRequest.V,
		Pretty:     catRequest.Pretty,
		Human:      catRequest.Human,
		ErrorTrace: catRequest.ErrorTrace,
		FilterPath: catRequest.FilterPath,
		Header:     catRequest.Header,
	}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatShards(
	ctx context.Context,
	catRequest proto.CatShardsRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatShardsRequest{}

	req.Index = catRequest.Index
	req.Bytes = catRequest.Bytes
	req.Format = catRequest.Format
	req.H = catRequest.H
	req.Help = catRequest.Help
	req.Local = catRequest.Local
	req.MasterTimeout = catRequest.MasterTimeout
	req.S = catRequest.S
	req.V = catRequest.V
	req.Pretty = catRequest.Pretty
	req.Human = catRequest.Human
	req.ErrorTrace = catRequest.ErrorTrace
	req.FilterPath = catRequest.FilterPath
	req.Header = catRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatCount(
	ctx context.Context,
	catRequest proto.CatCountRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatCountRequest{}

	req.Index = catRequest.Index
	req.Format = catRequest.Format
	req.H = catRequest.H
	req.Help = catRequest.Help
	req.S = catRequest.S
	req.V = catRequest.V
	req.Pretty = catRequest.Pretty
	req.Human = catRequest.Human
	req.ErrorTrace = catRequest.ErrorTrace
	req.FilterPath = catRequest.FilterPath
	req.Header = catRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatAllocationRequest(
	ctx context.Context,
	catRequest proto.CatAllocationRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatAllocationRequest{}

	req.NodeID = catRequest.NodeID
	req.Bytes = catRequest.Bytes
	req.Format = catRequest.Format
	req.H = catRequest.H
	req.Help = catRequest.Help
	req.Local = catRequest.Local
	req.MasterTimeout = catRequest.MasterTimeout
	req.S = catRequest.S
	req.V = catRequest.V
	req.Pretty = catRequest.Pretty
	req.Human = catRequest.Human
	req.ErrorTrace = catRequest.ErrorTrace
	req.FilterPath = catRequest.FilterPath
	req.Header = catRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatAliases(
	ctx context.Context,
	catRequest proto.CatAliasesRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatAliasesRequest{}

	req.Name = catRequest.Name
	req.Format = catRequest.Format
	req.H = catRequest.H
	req.Help = catRequest.Help
	req.Local = catRequest.Local
	req.S = catRequest.S
	req.V = catRequest.V
	req.Pretty = catRequest.Pretty
	req.Human = catRequest.Human
	req.ErrorTrace = catRequest.ErrorTrace
	req.FilterPath = catRequest.FilterPath
	req.Header = catRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesSegmentsRequest(
	ctx context.Context,
	human bool,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatSegmentsRequest{Human: human, Format: "json"}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsClusterStats(
	ctx context.Context,
	human bool,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.ClusterStatsRequest{Human: human}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCatNodes(
	ctx context.Context,
	h []string,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.CatNodesRequest{H: h, Format: "json"}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsDelete(
	ctx context.Context,
	deleteRequest proto.DeleteRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.DeleteRequest{}

	req.Index = deleteRequest.Index
	req.DocumentType = "_doc" //deleteRequest.DocumentType
	req.DocumentID = deleteRequest.DocumentID
	req.IfPrimaryTerm = deleteRequest.IfPrimaryTerm
	req.IfSeqNo = deleteRequest.IfSeqNo
	req.Refresh = deleteRequest.Refresh
	req.Routing = deleteRequest.Routing
	req.Timeout = deleteRequest.Timeout
	req.Version = deleteRequest.Version
	req.VersionType = deleteRequest.VersionType
	req.WaitForActiveShards = deleteRequest.WaitForActiveShards
	req.Pretty = deleteRequest.Pretty
	req.Human = deleteRequest.Human
	req.ErrorTrace = deleteRequest.ErrorTrace
	req.FilterPath = deleteRequest.FilterPath
	req.Header = deleteRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsUpdate(
	ctx context.Context,
	updateRequest proto.UpdateRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.UpdateRequest{}

	req.Index = updateRequest.Index
	req.DocumentType = "_doc" //updateRequest.DocumentType
	req.DocumentID = updateRequest.DocumentID
	req.Body = esutil.NewJSONReader(map[string]interface{}{
		"doc": body,
	})
	req.IfPrimaryTerm = updateRequest.IfPrimaryTerm
	req.IfSeqNo = updateRequest.IfSeqNo
	req.Lang = updateRequest.Lang
	req.Refresh = updateRequest.Refresh
	req.RetryOnConflict = updateRequest.RetryOnConflict
	req.Routing = updateRequest.Routing
	req.Source = updateRequest.Source
	req.SourceExcludes = updateRequest.SourceExcludes
	req.SourceIncludes = updateRequest.SourceIncludes
	req.Timeout = updateRequest.Timeout
	req.WaitForActiveShards = updateRequest.WaitForActiveShards
	req.Pretty = updateRequest.Pretty
	req.Human = updateRequest.Human
	req.ErrorTrace = updateRequest.ErrorTrace
	req.FilterPath = updateRequest.FilterPath
	req.Header = updateRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCreate(
	ctx context.Context,
	createRequest proto.CreateRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndexRequest{}
	req.Index = createRequest.Index
	req.DocumentType = "_doc" //createRequest.DocumentType
	req.DocumentID = createRequest.DocumentID
	req.Body = esutil.NewJSONReader(body)
	req.Pipeline = createRequest.Pipeline
	req.Refresh = createRequest.Refresh
	req.Routing = createRequest.Routing
	req.Timeout = createRequest.Timeout
	req.Version = createRequest.Version
	req.VersionType = createRequest.VersionType
	req.WaitForActiveShards = createRequest.WaitForActiveShards
	req.Pretty = createRequest.Pretty
	req.Human = createRequest.Human
	req.ErrorTrace = createRequest.ErrorTrace
	req.FilterPath = createRequest.FilterPath
	req.Header = createRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsSearch(
	ctx context.Context,
	searchRequest proto.SearchRequest,
	query interface{},
) (
	res *proto.Response,
	err error,
) {

	req := esapi.SearchRequest{}

	req.Index = searchRequest.Index
	req.DocumentType = []string{} //searchRequest.DocumentType
	req.Body = esutil.NewJSONReader(query)
	req.AllowNoIndices = searchRequest.AllowNoIndices
	req.AllowPartialSearchResults = searchRequest.AllowPartialSearchResults
	req.Analyzer = searchRequest.Analyzer
	req.AnalyzeWildcard = searchRequest.AnalyzeWildcard
	req.BatchedReduceSize = searchRequest.BatchedReduceSize
	req.CcsMinimizeRoundtrips = searchRequest.CcsMinimizeRoundtrips
	req.DefaultOperator = searchRequest.DefaultOperator
	req.Df = searchRequest.Df
	req.DocvalueFields = searchRequest.DocvalueFields
	req.ExpandWildcards = searchRequest.ExpandWildcards
	req.Explain = searchRequest.Explain
	req.From = searchRequest.From
	req.IgnoreThrottled = searchRequest.IgnoreThrottled
	req.IgnoreUnavailable = searchRequest.IgnoreUnavailable
	req.Lenient = searchRequest.Lenient
	req.MaxConcurrentShardRequests = searchRequest.MaxConcurrentShardRequests
	req.Preference = searchRequest.Preference
	req.PreFilterShardSize = searchRequest.PreFilterShardSize
	req.Query = searchRequest.Query
	req.RequestCache = searchRequest.RequestCache
	req.RestTotalHitsAsInt = searchRequest.RestTotalHitsAsInt
	req.Routing = searchRequest.Routing
	req.Scroll = searchRequest.Scroll
	req.SearchType = searchRequest.SearchType
	req.SeqNoPrimaryTerm = searchRequest.SeqNoPrimaryTerm
	req.Size = searchRequest.Size
	req.Sort = searchRequest.Sort
	req.Source = searchRequest.Source
	req.SourceExcludes = searchRequest.SourceExcludes
	req.SourceIncludes = searchRequest.SourceIncludes
	req.Stats = searchRequest.Stats
	req.StoredFields = searchRequest.StoredFields
	req.SuggestField = searchRequest.SuggestField
	req.SuggestMode = searchRequest.SuggestMode
	req.SuggestSize = searchRequest.SuggestSize
	req.SuggestText = searchRequest.SuggestText
	req.TerminateAfter = searchRequest.TerminateAfter
	req.Timeout = searchRequest.Timeout
	req.TrackScores = searchRequest.TrackScores
	req.TrackTotalHits = searchRequest.TrackTotalHits
	req.TypedKeys = searchRequest.TypedKeys
	req.Version = searchRequest.Version
	req.Pretty = searchRequest.Pretty
	req.Human = searchRequest.Human
	req.ErrorTrace = searchRequest.ErrorTrace
	req.FilterPath = searchRequest.FilterPath
	req.Header = searchRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesPutSettingsRequest(
	ctx context.Context,
	indexSettingsRequest proto.IndicesPutSettingsRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesPutSettingsRequest{}

	req.Index = indexSettingsRequest.Index
	req.Body = esutil.NewJSONReader(body)
	req.AllowNoIndices = indexSettingsRequest.AllowNoIndices
	req.ExpandWildcards = indexSettingsRequest.ExpandWildcards
	req.FlatSettings = indexSettingsRequest.FlatSettings
	req.IgnoreUnavailable = indexSettingsRequest.IgnoreUnavailable
	req.MasterTimeout = indexSettingsRequest.MasterTimeout
	req.PreserveExisting = indexSettingsRequest.PreserveExisting
	req.Timeout = indexSettingsRequest.Timeout
	req.Pretty = indexSettingsRequest.Pretty
	req.Human = indexSettingsRequest.Human
	req.ErrorTrace = indexSettingsRequest.ErrorTrace
	req.FilterPath = indexSettingsRequest.FilterPath
	req.Header = indexSettingsRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsCreateIndex(
	ctx context.Context,
	indexCreateRequest proto.IndicesCreateRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesCreateRequest{}

	req.Index = indexCreateRequest.Index
	req.Body = esutil.NewJSONReader(body)
	req.IncludeTypeName = indexCreateRequest.IncludeTypeName
	req.MasterTimeout = indexCreateRequest.MasterTimeout
	req.Timeout = indexCreateRequest.Timeout
	req.WaitForActiveShards = indexCreateRequest.WaitForActiveShards
	req.Pretty = indexCreateRequest.Pretty
	req.Human = indexCreateRequest.Human
	req.ErrorTrace = indexCreateRequest.ErrorTrace
	req.FilterPath = indexCreateRequest.FilterPath
	req.Header = indexCreateRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsDeleteIndex(
	ctx context.Context,
	indicesDeleteRequest proto.IndicesDeleteRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesDeleteRequest{}

	req.Index = indicesDeleteRequest.Index
	req.AllowNoIndices = indicesDeleteRequest.AllowNoIndices
	req.ExpandWildcards = indicesDeleteRequest.ExpandWildcards
	req.IgnoreUnavailable = indicesDeleteRequest.IgnoreUnavailable
	req.MasterTimeout = indicesDeleteRequest.MasterTimeout
	req.Timeout = indicesDeleteRequest.Timeout
	req.Pretty = indicesDeleteRequest.Pretty
	req.Human = indicesDeleteRequest.Human
	req.ErrorTrace = indicesDeleteRequest.ErrorTrace
	req.FilterPath = indicesDeleteRequest.FilterPath
	req.Header = indicesDeleteRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsReindex(
	ctx context.Context,
	reindexRequest proto.ReindexRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.ReindexRequest{}

	req.Body = esutil.NewJSONReader(body)
	req.MaxDocs = reindexRequest.MaxDocs
	req.Refresh = reindexRequest.Refresh
	req.RequestsPerSecond = reindexRequest.RequestsPerSecond
	req.Scroll = reindexRequest.Scroll
	S := cast.ToInt(reindexRequest.Slices)
	req.Slices = &S
	req.Timeout = reindexRequest.Timeout
	req.WaitForActiveShards = reindexRequest.WaitForActiveShards
	req.WaitForCompletion = reindexRequest.WaitForCompletion
	req.Pretty = reindexRequest.Pretty
	req.Human = reindexRequest.Human
	req.ErrorTrace = reindexRequest.ErrorTrace
	req.FilterPath = reindexRequest.FilterPath
	req.Header = reindexRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsIndicesGetSettingsRequest(
	ctx context.Context,
	indicesGetSettingsRequest proto.IndicesGetSettingsRequest,
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesGetSettingsRequest{}

	req.Index = indicesGetSettingsRequest.Index
	req.Name = indicesGetSettingsRequest.Name
	req.AllowNoIndices = indicesGetSettingsRequest.AllowNoIndices
	req.ExpandWildcards = indicesGetSettingsRequest.ExpandWildcards
	req.FlatSettings = indicesGetSettingsRequest.FlatSettings
	req.IgnoreUnavailable = indicesGetSettingsRequest.IgnoreUnavailable
	req.IncludeDefaults = indicesGetSettingsRequest.IncludeDefaults
	req.Local = indicesGetSettingsRequest.Local
	req.MasterTimeout = indicesGetSettingsRequest.MasterTimeout
	req.Pretty = indicesGetSettingsRequest.Pretty
	req.Human = indicesGetSettingsRequest.Human
	req.ErrorTrace = indicesGetSettingsRequest.ErrorTrace
	req.FilterPath = indicesGetSettingsRequest.FilterPath
	req.Header = indicesGetSettingsRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsPutMapping(
	ctx context.Context,
	indicesPutMappingRequest proto.IndicesPutMappingRequest,
	body interface{},
) (
	res *proto.Response,
	err error,
) {
	req := esapi.IndicesPutMappingRequest{}

	req.Index = indicesPutMappingRequest.Index
	req.DocumentType = "_doc" //indicesPutMappingRequest.DocumentType
	req.Body = esutil.NewJSONReader(body)
	req.AllowNoIndices = indicesPutMappingRequest.AllowNoIndices
	req.ExpandWildcards = indicesPutMappingRequest.ExpandWildcards
	req.IgnoreUnavailable = indicesPutMappingRequest.IgnoreUnavailable
	req.IncludeTypeName = indicesPutMappingRequest.IncludeTypeName
	req.MasterTimeout = indicesPutMappingRequest.MasterTimeout
	req.Timeout = indicesPutMappingRequest.Timeout
	req.Pretty = indicesPutMappingRequest.Pretty
	req.Human = indicesPutMappingRequest.Human
	req.ErrorTrace = indicesPutMappingRequest.ErrorTrace
	req.FilterPath = indicesPutMappingRequest.FilterPath
	req.Header = indicesPutMappingRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsGetAliases(
	ctx context.Context,
	indexNames []string,
) (
	res *proto.Response,
	err error,
) {

	httpRes, err := this.client.Indices.GetAlias(
		this.client.Indices.GetAlias.WithIndex(indexNames...),
		this.client.Indices.GetAlias.WithContext(ctx),
	)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsAddAliases(
	ctx context.Context,
	indexName []string,
	aliasName string,
) (
	res *proto.Response,
	err error,
) {

	httpRes, err := this.client.Indices.PutAlias(
		indexName,
		aliasName,
		this.client.Indices.PutAlias.WithContext(ctx),
	)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsRemoveAliases(
	ctx context.Context,
	indexName []string,
	aliasName []string,
) (
	res *proto.Response,
	err error,
) {

	httpRes, err := this.client.Indices.DeleteAlias(
		indexName,
		aliasName,
		this.client.Indices.DeleteAlias.WithContext(ctx),
	)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsMoveToAnotherIndexAliases(
	ctx context.Context,
	body proto.AliasAction,
) (
	res *proto.Response,
	err error,
) {

	httpRes, err := this.client.Indices.UpdateAliases(
		esutil.NewJSONReader(body),
		this.client.Indices.UpdateAliases.WithContext(ctx),
	)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsTaskList(
	ctx context.Context,
) (
	res *proto.Response,
	err error,
) {
	httpRes, err := this.client.Tasks.List(
		this.client.Tasks.List.WithDetailed(true),
		this.client.Tasks.List.WithContext(ctx))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) EsTasksCancel(
	ctx context.Context,
	taskId string,
) (
	res *proto.Response,
	err error,
) {
	httpRes, err := this.client.Tasks.Cancel(
		this.client.Tasks.Cancel.WithTaskID(taskId),
		this.client.Tasks.Cancel.WithContext(ctx))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	res, err = proto.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}
