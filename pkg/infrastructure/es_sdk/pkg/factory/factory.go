package factory

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/clickhouse"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/mongo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/mysql"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/postgresql"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/redis"
	v6 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v6"
	v7 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v7"
	v8 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v8"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"strings"
)

var EsServiceMap = map[string]func(cfg *proto.Config) (pkg.ClientInterface, error){
	"elasticsearch6.x": v6.NewEsClient6,
	"elasticsearch7.x": v7.NewEsClient7,
	"elasticsearch8.x": v8.NewEsClient8,
	"mysql":            mysql.NewMysqlClient,
	"redis":            redis.NewRedisClient,
	"clickhouse":       clickhouse.NewClickhouseClient,
	"postgres":         postgresql.NewPostgresqlClient,
	"mongo":            mongo.NewMongoClient,
}

func NewEsService(cfg *proto.Config) (pkg.ClientInterface, error) {
	var found bool
	var fn func(cfg *proto.Config) (pkg.ClientInterface, error)
	if fn, found = EsServiceMap[cfg.Version]; !found {
		return nil, VersionErr()
	}
	fn = EsServiceMap[cfg.Version]
	return fn(cfg)
}

func VersionErr() error {
	datasources := []string{}
	for key := range EsServiceMap {
		datasources = append(datasources, key)
	}

	return fmt.Errorf("暂只支持（%s）", strings.Join(datasources, ","))
}
