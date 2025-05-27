package factory

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/clickhouse"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/dameng"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/hive"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/mariadb"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/mongo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/mysql"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/oracle"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/postgresql"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/redis"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/spark"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/sqlserver"
	v6 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v6"
	v7 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v7"
	v8 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/v8"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"strings"
)

var EsServiceMap = map[string]func(cfg *proto.Config) (pkg.ClientInterface, error){
	pkg.ElasticSearch6: v6.NewEsClient6,
	pkg.ElasticSearch7: v7.NewEsClient7,
	pkg.ElasticSearch8: v8.NewEsClient8,
	pkg.Mysql:          mysql.NewMysqlClient,
	pkg.Redis:          redis.NewRedisClient,
	pkg.ClickHouse:     clickhouse.NewClickhouseClient,
	pkg.Postgres:       postgresql.NewPostgresqlClient,
	pkg.Mongo:          mongo.NewMongoClient,
	"dameng":           dameng.NewDamengClient,
	"oracle":           oracle.NewOracleClient,
	"sqlserver":        sqlserver.NewSqlServerClient,
	"mariadb":          mariadb.NewMariaDBClient,
	"hive":             hive.NewHiveClient,
	"spark":            spark.NewSparkClient,
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
