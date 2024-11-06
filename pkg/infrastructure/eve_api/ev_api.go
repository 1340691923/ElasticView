package eve_api

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/api"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/vo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/imroc/req/v2"
	"sync"
)

type EvEApi struct {
	log         *logger.AppLogger
	client      *req.Client
	accessToken string
	lock        *sync.RWMutex
	isDebug     bool
}

func NewEvApi(cfg *config.Config, log *logger.AppLogger) *EvEApi {

	var client *req.Client

	if cfg.DeBug {
		client = req.C().DevMode()
	} else {
		client = req.C()
	}

	return &EvEApi{log: log, client: client, lock: new(sync.RWMutex), isDebug: cfg.DeBug}
}

func (this *EvEApi) Request(ctx context.Context, api api.API, requestData interface{}, result *vo.ApiCommonRes) error {

	header := map[string]string{}
	if this.GetAccessToken() != "" {
		header["Ev-Token"] = this.GetAccessToken()
	}

	_, err := this.client.R().
		SetContext(ctx).
		SetBody(requestData).
		SetHeaders(header).
		SetResult(&result).
		Post(fmt.Sprintf("%s%s", this.GetEvBackDomain(), api))
	if err != nil {
		return err
	}
	return result.Error()
}

func (this *EvEApi) SetAccessToken(token string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.accessToken = token
}

func (this *EvEApi) GetAccessToken() string {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.accessToken
}

func (this *EvEApi) GetEvBackDomain() string {
	if this.isDebug {
		return "http://127.0.0.1:8199/v1/api/"
	}
	return "http://dev.elastic-view.cn/v1/api/"
}
