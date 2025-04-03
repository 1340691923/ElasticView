package base

import (
	"context"
	"fmt"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

var NotAllowConnType = errors.New("请选择正确的数据源")

type BaseDatasource struct {
}

func (b *BaseDatasource) RedisExecCommand(ctx context.Context, dbName int, args ...interface{}) (data interface{}, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (columns []string, list []map[string]interface{}, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (data map[string]interface{}, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) Ping(ctx context.Context) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsVersion() (version int, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatNodes(ctx context.Context, h []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsClusterStats(ctx context.Context, human bool) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesSegmentsRequest(ctx context.Context, human bool) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsPerformRequest(ctx context.Context, req *http.Request) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsRefresh(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsOpen(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsFlush(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesClearCache(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesClose(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesForcemerge(ctx context.Context, indexNames []string, maxNumSegments *int) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsDeleteByQuery(ctx context.Context, indexNames []string, documents []string, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotCreate(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotDelete(ctx context.Context, repository string, snapshot string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsRestoreSnapshot(ctx context.Context, repository string, snapshot string, waitForCompletion *bool, reqJson proto.Json) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotStatus(ctx context.Context, repository string, snapshot []string, ignoreUnavailable *bool) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotGetRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotCreateRepository(ctx context.Context, repository string, reqJson proto.Json) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSnapshotDeleteRepository(ctx context.Context, repository []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsGetIndices(ctx context.Context, catIndicesRequest proto.CatIndicesRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatHealth(ctx context.Context, catRequest proto.CatHealthRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatShards(ctx context.Context, catRequest proto.CatShardsRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatCount(ctx context.Context, catRequest proto.CatCountRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatAllocationRequest(ctx context.Context, catRequest proto.CatAllocationRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCatAliases(ctx context.Context, catRequest proto.CatAliasesRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsDelete(ctx context.Context, deleteRequest proto.DeleteRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsUpdate(ctx context.Context, updateRequest proto.UpdateRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCreate(ctx context.Context, createRequest proto.CreateRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsSearch(ctx context.Context, searchRequest proto.SearchRequest, query interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesPutSettingsRequest(ctx context.Context, indexSettingsRequest proto.IndicesPutSettingsRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsCreateIndex(ctx context.Context, indexCreateRequest proto.IndicesCreateRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsDeleteIndex(ctx context.Context, indicesDeleteRequest proto.IndicesDeleteRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsReindex(ctx context.Context, reindexRequest proto.ReindexRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsIndicesGetSettingsRequest(ctx context.Context, indicesGetSettingsRequest proto.IndicesGetSettingsRequest) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsPutMapping(ctx context.Context, indicesPutMappingRequest proto.IndicesPutMappingRequest, body interface{}) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsGetMapping(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsGetAliases(ctx context.Context, indexNames []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsAddAliases(ctx context.Context, indexName []string, aliasName string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsRemoveAliases(ctx context.Context, indexName []string, aliasName []string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsMoveToAnotherIndexAliases(ctx context.Context, body proto.AliasAction) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsTaskList(ctx context.Context) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) EsTasksCancel(ctx context.Context, taskId string) (res *proto.Response, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) ExecMongoCommand(ctx context.Context, dbName string, command bson.D, timeout time.Duration) (res bson.M, err error) {
	err = NotAllowConnType
	return
}

func (b *BaseDatasource) ExtractIPPort(address string) (string, string, error) {
	// 移除协议前缀
	address = strings.TrimPrefix(address, "http://")
	address = strings.TrimPrefix(address, "tcp://")

	// 使用 net.SplitHostPort 获取 IP 和端口
	ip, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", "", fmt.Errorf("提取ip,端口失败: %v", err)
	}

	return ip, port, nil
}

func (b *BaseDatasource) ShowMongoDbs(ctx context.Context) ([]string, error) {
	return nil, NotAllowConnType
}
