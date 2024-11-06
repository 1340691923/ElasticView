package es_service

import (
	"bytes"
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

type EsService struct {
	orm *sqlstore.SqlStore
}

func NewEsService(orm *sqlstore.SqlStore) *EsService {
	return &EsService{orm: orm}
}

func (this *EsService) Ping(ctx context.Context, esClient pkg.ClientInterface) (res *vo.PingResult, err error) {

	resp, err := esClient.Ping(ctx)

	if err != nil {
		return
	}

	if resp.StatusErr() != nil {
		err = resp.StatusErr()
		return
	}

	res = new(vo.PingResult)

	err = json.Unmarshal(resp.ResByte(), res)

	if err != nil {
		return
	}

	return
}

func (this *EsService) RunDsl(
	ctx context.Context,
	esI pkg.ClientInterface,
	userID int,
	method string,
	path string,
	body string,
) (res *proto.Response, err error) {
	method = strings.ToUpper(method)

	u, err := url.Parse(path)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	path = u.Path
	query := u.Query()
	query.Add("format", "json")

	var req *http.Request

	version, err := esI.EsVersion()

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if version > 6 {
		if len(path) > 0 {
			if path[0:1] != "/" {
				path = "/" + path
			}
		}
	}

	if body == "" {
		req, err = http.NewRequest(method, path+"?"+query.Encode(), nil)
	} else {
		req, err = http.NewRequest(method, path+"?"+query.Encode(), bytes.NewReader([]byte(body)))
	}

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	res, err = esI.EsPerformRequest(ctx, req)

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (this *EsService) EsIndexCount(ctx context.Context, esClient pkg.ClientInterface) (indexNameLen int, err error) {
	catIndicesResponse, err := esClient.EsGetIndices(ctx, proto.CatIndicesRequest{
		Format: "json",
	})
	if err != nil {
		return
	}
	if catIndicesResponse.StatusErr() != nil {
		err = catIndicesResponse.StatusErr()
		return
	}
	var list []proto.CatIndex
	err = json.Unmarshal(catIndicesResponse.ResByte(), &list)
	if err != nil {
		return
	}
	indexNameLen = len(list)
	return
}

func (this *EsService) CatHealth(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatHealth(ctx, proto.CatHealthRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *EsService) CatShards(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatShards(ctx, proto.CatShardsRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *EsService) CatCount(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatCount(ctx, proto.CatCountRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *EsService) CatAllocation(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatAllocationRequest(ctx, proto.CatAllocationRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *EsService) CatAliases(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatAliases(ctx, proto.CatAliasesRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *EsService) CatNodes(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsCatNodes(ctx, strings.Split("ip,name,heap.percent,heap.current,heap.max,ram.percent,ram.current,ram.max,node.role,master,cpu,load_1m,load_5m,load_15m,disk.used_percent,disk.used,disk.total", ","))
	return
}

func (this *EsService) CatIndices(ctx context.Context, esClient pkg.ClientInterface, sort []string, indexBytesFormat string) (res *proto.Response, err error) {

	req := proto.CatIndicesRequest{}
	req.S = sort
	req.Human = true
	req.Format = "json"
	if indexBytesFormat != "" {
		req.Bytes = indexBytesFormat
	}
	res, err = esClient.EsGetIndices(ctx, req)
	return
}

func (this *EsService) IndicesSegmentsRequest(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsIndicesSegmentsRequest(ctx, true)
	return
}

func (this *EsService) ClusterStats(ctx context.Context, esClient pkg.ClientInterface) (res *proto.Response, err error) {
	res, err = esClient.EsClusterStats(ctx, true)
	return
}
