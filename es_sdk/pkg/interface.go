package pkg

import (
	"context"
	proto2 "github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/escache"
	"net/http"
)

type EsI interface {
	Version() int
	CatNodes(ctx context.Context, h []string) (res *proto2.Response, err error)
	ClusterStats(ctx context.Context, human bool) (res *proto2.Response, err error)
	CatSegments(ctx context.Context, human bool) (res *proto2.Response, err error)
	PerformRequest(ctx context.Context, req *http.Request) (res *proto2.Response, err error)
	Ping(ctx context.Context) (res *proto2.Response, err error)
	Refresh(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	Open(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	Flush(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	IndicesClearCache(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	IndicesClose(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	IndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto2.Response, err error)
	DeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto2.Response, err error)
	SnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson escache.Json) (res *proto2.Response, err error)
	SnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto2.Response, err error)
	RestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson escache.Json) (res *proto2.Response, err error)
	SnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto2.Response, err error)

	SnapshotGetRepository(ctx context.Context, repository []string) (res *proto2.Response, err error)
	SnapshotCreateRepository(ctx context.Context, repository string, reqJson escache.Json) (res *proto2.Response, err error)
	SnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto2.Response, err error)

	GetIndices(ctx context.Context, catIndicesRequest proto2.CatIndicesRequest) (res *proto2.Response, err error)
	CatHealth(ctx context.Context, catRequest proto2.CatHealthRequest) (res *proto2.Response, err error)
	CatShards(ctx context.Context, catRequest proto2.CatShardsRequest) (res *proto2.Response, err error)
	CatCount(ctx context.Context, catRequest proto2.CatCountRequest) (res *proto2.Response, err error)
	CatAllocationRequest(ctx context.Context, catRequest proto2.CatAllocationRequest) (res *proto2.Response, err error)
	CatAliases(ctx context.Context, catRequest proto2.CatAliasesRequest) (res *proto2.Response, err error)

	Delete(ctx context.Context, deleteRequest proto2.DeleteRequest) (res *proto2.Response, err error)
	Update(ctx context.Context, updateRequest proto2.UpdateRequest, body interface{}) (res *proto2.Response, err error)
	Create(ctx context.Context, createRequest proto2.CreateRequest, body interface{}) (res *proto2.Response, err error)
	Search(ctx context.Context, searchRequest proto2.SearchRequest, query interface{}) (res *proto2.Response, err error)

	IndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto2.IndicesPutSettingsRequest, body interface{}) (res *proto2.Response, err error)
	CreateIndex(ctx context.Context, indexCreateRequest proto2.IndicesCreateRequest, body interface{}) (res *proto2.Response, err error)
	DeleteIndex(ctx context.Context, indicesDeleteRequest proto2.IndicesDeleteRequest) (res *proto2.Response, err error)
	Reindex(ctx context.Context, reindexRequest proto2.ReindexRequest, body interface{}) (res *proto2.Response, err error)
	IndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto2.IndicesGetSettingsRequest) (res *proto2.Response, err error)

	PutMapping(ctx context.Context, indicesPutMappingRequest proto2.IndicesPutMappingRequest, body interface{}) (res *proto2.Response, err error)
	GetMapping(ctx context.Context, indexNames []string) (res *proto2.Response, err error)

	GetAliases(ctx context.Context, indexNames []string) (res *proto2.Response, err error)
	AddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto2.Response, err error)
	RemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto2.Response, err error)
	MoveToAnotherIndexAliases(ctx context.Context, body proto2.AliasAction) (res *proto2.Response, err error)

	TaskList(ctx context.Context) (res *proto2.Response, err error)
	TasksCancel(ctx context.Context, taskId string) (res *proto2.Response, err error)
}
