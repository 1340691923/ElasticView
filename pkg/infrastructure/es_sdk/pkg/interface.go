package pkg

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"net/http"
)

type EsI interface {
	Version() int
	CatNodes(ctx context.Context, h []string) (res *proto.Response, err error)
	ClusterStats(ctx context.Context, human bool) (res *proto.Response, err error)
	IndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error)
	PerformRequest(ctx context.Context, req *http.Request) (res *proto.Response, err error)
	Ping(ctx context.Context) (res *proto.Response, err error)
	Refresh(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	Open(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	Flush(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	IndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	IndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	IndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error)
	DeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error)
	SnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error)
	SnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error)
	RestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error)
	SnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error)

	SnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error)
	SnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error)
	SnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error)

	GetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error)
	CatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error)
	CatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error)
	CatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error)
	CatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error)
	CatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error)

	Delete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error)
	Update(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error)
	Create(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error)
	Search(ctx context.Context, searchRequest proto.SearchRequest, query interface{}) (res *proto.Response, err error)

	IndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error)
	CreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error)
	DeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error)
	Reindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error)
	IndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error)

	PutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error)
	GetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error)

	GetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error)
	AddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error)
	RemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error)
	MoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error)

	TaskList(ctx context.Context) (res *proto.Response, err error)
	TasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error)
}
