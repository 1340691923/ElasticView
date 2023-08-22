package v7

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/escache"
	elasticV7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"net/http"
)

type EsClient7 struct {
	client *elasticV7.Client
}

func NewEsClient7(cfg proto2.Config) (pkg.EsI, error) {
	obj := &EsClient7{}
	esCfg, err := cfg.ConvertV7(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	obj.init(esCfg)
	return obj, nil
}

func (this *EsClient7) init(config elasticV7.Config) (err error) {

	this.client, err = elasticV7.NewClient(config)
	if err != nil {
		fmt.Printf("Error creating the client: %s\n", err)
		return
	}
	return
}

func (this *EsClient7) Version() int {
	return 7
}

func (this *EsClient7) Ping(
	ctx context.Context,
) (
	res *proto2.Response,
	err error,
) {
	httpRes, err := this.client.Ping(
		this.client.Ping.WithContext(ctx),
	)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) GetMapping(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesGetMappingRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) SnapshotCreate(
	ctx context.Context,
	repository string,
	snapshot string,
	waitForCompletion *bool,
	reqJson escache.Json,
) (
	res *proto2.Response, err error) {
	snapshotCreateService := esapi.SnapshotCreateRequest{
		Repository:        repository,
		Snapshot:          snapshot,
		Body:              esutil.NewJSONReader(reqJson),
		WaitForCompletion: waitForCompletion,
	}

	var httpRes *esapi.Response
	httpRes, err = snapshotCreateService.Do(ctx, this.client)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) PerformRequest(
	ctx context.Context,
	req *http.Request,
) (
	res *proto2.Response, err error) {
	if req != nil {
		req = req.WithContext(ctx)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpRes, err := this.client.Perform(req)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)
	return
}

func (this *EsClient7) SnapshotDelete(
	ctx context.Context,
	repository string,
	snapshot string,
) (
	res *proto2.Response,
	err error) {
	req := esapi.SnapshotDeleteRequest{
		Repository: repository,
		Snapshot:   snapshot,
	}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) RestoreSnapshot(
	ctx context.Context,
	repository string,
	snapshot string,
	waitForCompletion *bool,
	reqJson escache.Json,
) (
	res *proto2.Response,
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
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Refresh(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesRefreshRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Open(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesOpenRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Flush(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesFlushRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) IndicesClearCache(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesClearCacheRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) IndicesClose(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesCloseRequest{Index: indexNames}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) IndicesForcemerge(
	ctx context.Context,
	indexNames []string,
	maxNumSegments *int,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.IndicesForcemergeRequest{Index: indexNames, MaxNumSegments: maxNumSegments}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) DeleteByQuery(
	ctx context.Context,
	indexNames []string,
	documents []string,
	body interface{},
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.DeleteByQueryRequest{Index: indexNames, DocumentType: []string{}, Body: esutil.NewJSONReader(body)}
	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) SnapshotStatus(
	ctx context.Context,
	repository string,
	snapshot []string,
	ignoreUnavailable *bool,
) (
	res *proto2.Response,
	err error,
) {
	request := esapi.SnapshotStatusRequest{
		Repository:        repository,
		Snapshot:          snapshot,
		IgnoreUnavailable: ignoreUnavailable,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) SnapshotGetRepository(
	ctx context.Context,
	repository []string,
) (
	res *proto2.Response,
	err error,
) {
	request := esapi.SnapshotGetRepositoryRequest{
		Repository: repository,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) SnapshotCreateRepository(
	ctx context.Context,
	repository string,
	reqJson escache.Json,
) (
	res *proto2.Response,
	err error,
) {
	request := esapi.SnapshotCreateRepositoryRequest{
		Repository: repository,
		Body:       esutil.NewJSONReader(reqJson),
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) SnapshotDeleteRepository(
	ctx context.Context,
	repository []string,
) (
	res *proto2.Response,
	err error,
) {
	request := esapi.SnapshotDeleteRepositoryRequest{
		Repository: repository,
	}

	httpRes, err := request.Do(ctx, this.client)

	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) GetIndices(
	ctx context.Context,
	catRequest proto2.CatIndicesRequest,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.CatIndicesRequest{
		Index:                   catRequest.Index,
		Bytes:                   catRequest.Bytes,
		ExpandWildcards:         catRequest.ExpandWildcards,
		Format:                  catRequest.Format,
		H:                       catRequest.H,
		Health:                  catRequest.Health,
		Help:                    catRequest.Help,
		IncludeUnloadedSegments: catRequest.IncludeUnloadedSegments,
		Local:                   catRequest.Local,
		MasterTimeout:           catRequest.MasterTimeout,
		Pri:                     catRequest.Pri,
		S:                       catRequest.S,
		Time:                    catRequest.Time,
		V:                       catRequest.V,
		Pretty:                  catRequest.Pretty,
		Human:                   catRequest.Human,
		ErrorTrace:              catRequest.ErrorTrace,
		FilterPath:              catRequest.FilterPath,
		Header:                  catRequest.Header,
	}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatHealth(
	ctx context.Context,
	catRequest proto2.CatHealthRequest,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.CatHealthRequest{
		Format:     catRequest.Format,
		H:          catRequest.H,
		Help:       catRequest.Help,
		S:          catRequest.S,
		Time:       catRequest.Time,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatShards(
	ctx context.Context,
	catRequest proto2.CatShardsRequest,
) (
	res *proto2.Response,
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
	req.Time = catRequest.Time
	req.V = catRequest.V
	req.Pretty = catRequest.Pretty
	req.Human = catRequest.Human
	req.ErrorTrace = catRequest.ErrorTrace
	req.FilterPath = catRequest.FilterPath
	req.Header = catRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatCount(
	ctx context.Context,
	catRequest proto2.CatCountRequest,
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatAllocationRequest(
	ctx context.Context,
	catRequest proto2.CatAllocationRequest,
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatAliases(
	ctx context.Context,
	catRequest proto2.CatAliasesRequest,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.CatAliasesRequest{}

	req.Name = catRequest.Name
	req.ExpandWildcards = catRequest.ExpandWildcards
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatSegments(
	ctx context.Context,
	human bool,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.CatSegmentsRequest{Human: human}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) ClusterStats(
	ctx context.Context,
	human bool,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.ClusterStatsRequest{Human: human}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CatNodes(
	ctx context.Context,
	h []string,
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.CatNodesRequest{H: h, Format: "json"}

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Delete(
	ctx context.Context,
	deleteRequest proto2.DeleteRequest,
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Update(
	ctx context.Context,
	updateRequest proto2.UpdateRequest,
	body interface{},
) (
	res *proto2.Response,
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
	req.RequireAlias = updateRequest.RequireAlias
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Create(
	ctx context.Context,
	createRequest proto2.CreateRequest,
	body interface{},
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Search(
	ctx context.Context,
	searchRequest proto2.SearchRequest,
	query interface{},
) (
	res *proto2.Response,
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
	req.MinCompatibleShardNode = searchRequest.MinCompatibleShardNode
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) IndicesPutSettingsRequest(
	ctx context.Context,
	indexSettingsRequest proto2.IndicesPutSettingsRequest,
	body interface{},
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) CreateIndex(
	ctx context.Context,
	indexCreateRequest proto2.IndicesCreateRequest,
	body interface{},
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) DeleteIndex(
	ctx context.Context,
	indicesDeleteRequest proto2.IndicesDeleteRequest,
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) Reindex(
	ctx context.Context,
	reindexRequest proto2.ReindexRequest,
	body interface{},
) (
	res *proto2.Response,
	err error,
) {
	req := esapi.ReindexRequest{}

	req.Body = esutil.NewJSONReader(body)
	req.MaxDocs = reindexRequest.MaxDocs
	req.Refresh = reindexRequest.Refresh
	req.RequestsPerSecond = reindexRequest.RequestsPerSecond
	req.Scroll = reindexRequest.Scroll
	req.Slices = reindexRequest.Slices
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) IndicesGetSettingsRequest(
	ctx context.Context,
	indicesGetSettingsRequest proto2.IndicesGetSettingsRequest,
) (
	res *proto2.Response,
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
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) PutMapping(
	ctx context.Context,
	indicesPutMappingRequest proto2.IndicesPutMappingRequest,
	body interface{},
) (
	res *proto2.Response,
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
	req.WriteIndexOnly = indicesPutMappingRequest.WriteIndexOnly
	req.Pretty = indicesPutMappingRequest.Pretty
	req.Human = indicesPutMappingRequest.Human
	req.ErrorTrace = indicesPutMappingRequest.ErrorTrace
	req.FilterPath = indicesPutMappingRequest.FilterPath
	req.Header = indicesPutMappingRequest.Header

	httpRes, err := req.Do(ctx, this.client)
	if err != nil {
		return nil, err
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) GetAliases(
	ctx context.Context,
	indexNames []string,
) (
	res *proto2.Response,
	err error,
) {

	httpRes, err := this.client.Indices.GetAlias(
		this.client.Indices.GetAlias.WithIndex(indexNames...),
		this.client.Indices.GetAlias.WithContext(ctx),
	)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) AddAliases(
	ctx context.Context,
	indexName []string,
	aliasName string,
) (
	res *proto2.Response,
	err error,
) {

	httpRes, err := this.client.Indices.PutAlias(
		indexName,
		aliasName,
		this.client.Indices.PutAlias.WithContext(ctx),
	)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) RemoveAliases(
	ctx context.Context,
	indexName []string,
	aliasName []string,
) (
	res *proto2.Response,
	err error,
) {

	httpRes, err := this.client.Indices.DeleteAlias(
		indexName,
		aliasName,
		this.client.Indices.DeleteAlias.WithContext(ctx),
	)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) MoveToAnotherIndexAliases(
	ctx context.Context,
	body proto2.AliasAction,
) (
	res *proto2.Response,
	err error,
) {

	httpRes, err := this.client.Indices.UpdateAliases(
		esutil.NewJSONReader(body),
		this.client.Indices.UpdateAliases.WithContext(ctx),
	)
	if err != nil {
		return
	}

	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) TaskList(
	ctx context.Context,
) (
	res *proto2.Response,
	err error,
) {
	httpRes, err := this.client.Tasks.List(
		this.client.Tasks.List.WithDetailed(true),
		this.client.Tasks.List.WithContext(ctx))
	if err != nil {
		return
	}
	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}

func (this *EsClient7) TasksCancel(
	ctx context.Context,
	taskId string,
) (
	res *proto2.Response,
	err error,
) {
	httpRes, err := this.client.Tasks.Cancel(
		this.client.Tasks.Cancel.WithTaskID(taskId),
		this.client.Tasks.Cancel.WithContext(ctx))
	if err != nil {
		return
	}
	res, err = proto2.NewResponse(httpRes.StatusCode, httpRes.Header, httpRes.Body)

	return
}
