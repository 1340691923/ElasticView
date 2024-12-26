package mongo

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/base"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/bson"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/pkg/errors"
	bson2 "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoClient struct {
	base.BaseDatasource
	db *mongo.Client
}

func NewMongoClient(cfg *proto2.Config) (pkg.ClientInterface, error) {
	ds, ok := cache.GetDataSourceCache(cfg.ConnectId)

	if !ok {
		obj := &MongoClient{}

		if len(cfg.Addresses) == 0 {
			return nil, errors.New("ip和端口不能为空")
		}

		ip, port, err := obj.ExtractIPPort(cfg.Addresses[0])

		if err != nil {
			return nil, errors.WithStack(err)
		}

		auth := ""

		if cfg.Username != "" && cfg.Password != "" {
			auth = fmt.Sprintf("%s:%s@", cfg.Username, cfg.Password)
		}

		clientOpts := options.Client().ApplyURI(
			fmt.Sprintf("mongodb://%s%s:%s",
				auth, ip, port),
		)
		client, err := mongo.Connect(context.Background(), clientOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
		}
		obj.db = client

		cache.SaveDataSourceCache(cfg.ConnectId, obj)
		return obj, nil
	}

	return ds, nil
}

func (this *MongoClient) Ping(
	ctx context.Context,
) (
	res *proto.Response,
	err error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = this.db.Ping(context.Background(), nil)
	if err != nil {
		return
	}

	res = proto.NewResponseNotErr()

	return
}

func (this *MongoClient) ExecMongoCommand(ctx context.Context, dbName string, command bson.D, timeout time.Duration) (res bson.M, err error) {

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var command2 bson2.D

	for _, v := range command {
		command2 = append(command2, bson2.E{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	// 执行命令
	var result bson2.M
	err = this.db.Database(dbName).RunCommand(ctx, command2).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute command: %v", err)
	}

	res = bson.M{}

	for k, v := range result {
		res[k] = v
	}

	return res, nil
}

func (this *MongoClient) ShowMongoDbs(ctx context.Context) ([]string, error) {
	return this.db.ListDatabaseNames(ctx, bson2.D{})
}
