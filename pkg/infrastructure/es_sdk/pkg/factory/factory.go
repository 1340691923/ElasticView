package factory

import (
	"errors"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	v6 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v6"
	v7 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v7"
	v8 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v8"
)

var VerError = errors.New("ES版本暂只支持6,7,8")

var EsServiceMap = map[int]func(cfg proto.Config) (pkg.EsI, error){
	6: v6.NewEsClient6,
	7: v7.NewEsClient7,
	8: v8.NewEsClient8,
}

func NewEsService(cfg proto.Config) (pkg.EsI, error) {
	var found bool
	var fn func(cfg proto.Config) (pkg.EsI, error)
	if fn, found = EsServiceMap[cfg.Version]; !found {
		return nil, VerError
	}
	fn = EsServiceMap[cfg.Version]
	return fn(cfg)
}
