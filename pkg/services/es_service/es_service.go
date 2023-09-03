package es_service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/ElasticView/pkg/vo"
	"net/http"
)

type EsService struct{}

func NewEsService() *EsService {
	return &EsService{}
}

func (this *EsService) Ping(ctx context.Context, esClient pkg.EsI) (res *vo.PingResult, err error) {

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return
	}

	resp, err := esClient.PerformRequest(ctx, req)

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

func (this *EsService) RecoverCanWrite(ctx context.Context, esClient pkg.EsI) (err error) {

	body := util.Map{
		"index": util.Map{
			"blocks": util.Map{
				"read_only_allow_delete": "false",
			},
		},
	}

	js, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, "/_settings", bytes.NewBuffer(js))

	if err != nil {
		return
	}

	res, err := esClient.PerformRequest(ctx, req)

	if err != nil {
		return
	}
	if res.StatusErr() != nil {
		err = res.StatusErr()
		return
	}

	return
}
