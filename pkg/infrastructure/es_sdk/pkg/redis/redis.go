package redis

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/base"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisClient struct {
	base.BaseDatasource
	client *redis.Client
}

func NewRedisClient(cfg *proto2.Config) (pkg.ClientInterface, error) {
	ds, ok := cache.GetDataSourceCache(cfg.ConnectId)

	if !ok {
		obj := &RedisClient{}

		if len(cfg.Addresses) == 0 {
			return nil, errors.New("ip和端口不能为空")
		}

		options := &redis.Options{
			Addr:     cfg.Addresses[0], // Redis 服务器地址
			Password: "",
		}
		if cfg.Username != "" {
			options.Username = cfg.Username
		}
		if cfg.Password != "" {
			options.Password = cfg.Password
		}

		// 打开数据库连接
		rdb := redis.NewClient(options)
		obj.client = rdb

		cache.SaveDataSourceCache(cfg.ConnectId, obj)
		return obj, nil
	}
	return ds, nil
}

func (this *RedisClient) Ping(
	ctx context.Context,
) (
	res *proto.Response,
	err error,
) {
	err = this.client.Ping(ctx).Err()
	if err != nil {
		return
	}

	res = proto.NewResponseNotErr()

	return
}

func (this *RedisClient) RedisExecCommand(ctx context.Context, dbName int, args ...interface{}) (data interface{}, err error) {

	err = this.client.Do(ctx, "SELECT", dbName).Err()
	if err != nil {
		return
	}

	data, err = this.client.Do(ctx, args...).Result()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return
}
