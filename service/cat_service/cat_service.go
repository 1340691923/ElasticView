package cat_service

import (
	"context"
	"github.com/1340691923/ElasticView/es_sdk/pkg"
	proto2 "github.com/1340691923/ElasticView/es_sdk/pkg/proto"
	"strings"
)

type CatService struct {
	esClient pkg.EsI
}

func NewCatService(esClient pkg.EsI) *CatService {
	return &CatService{esClient: esClient}
}

func (this *CatService) CatHealth(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatHealth(ctx, proto2.CatHealthRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *CatService) CatShards(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatShards(ctx, proto2.CatShardsRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *CatService) CatCount(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatCount(ctx, proto2.CatCountRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *CatService) CatAllocation(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatAllocationRequest(ctx, proto2.CatAllocationRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *CatService) CatAliases(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatAliases(ctx, proto2.CatAliasesRequest{
		Format: "json",
		Human:  true,
	})
	return
}

func (this *CatService) CatNodes(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatNodes(ctx, strings.Split("ip,name,heap.percent,heap.current,heap.max,ram.percent,ram.current,ram.max,node.role,master,cpu,load_1m,load_5m,load_15m,disk.used_percent,disk.used,disk.total", ","))
	return
}

func (this *CatService) CatIndices(ctx context.Context, sort []string, indexBytesFormat string) (res *proto2.Response, err error) {

	req := proto2.CatIndicesRequest{}
	req.S = sort
	req.Human = true
	req.Format = "json"
	if indexBytesFormat != "" {
		req.Bytes = indexBytesFormat
	}
	res, err = this.esClient.GetIndices(ctx, req)
	return
}

func (this *CatService) CatSegments(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.CatSegments(ctx, true)
	return
}

func (this *CatService) ClusterStats(ctx context.Context) (res *proto2.Response, err error) {
	res, err = this.esClient.ClusterStats(ctx, true)
	return
}
