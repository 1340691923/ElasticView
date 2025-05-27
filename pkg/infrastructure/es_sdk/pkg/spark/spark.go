package spark

import (
	"context"
	"database/sql"
	sql2 "database/sql"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/base"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/utils"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SparkClient struct {
	base.BaseDatasource
	db        *gorm.DB
	transport thrift.TTransport
}

func NewSparkClient(cfg *proto2.Config) (pkg.ClientInterface, error) {
	ds, ok := cache.GetDataSourceCache(cfg.ConnectId)
	
	if !ok {
		obj := &SparkClient{}
		
		if len(cfg.Addresses) == 0 {
			return nil, errors.New("ip和端口不能为空")
		}
		
		ip, port, err := obj.ExtractIPPort(cfg.Addresses[0])
		if err != nil {
			return nil, errors.WithStack(err)
		}
		
		transportFactory := thrift.NewTBufferedTransportFactory(8192)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		
		socket, err := thrift.NewTSocket(fmt.Sprintf("%s:%s", ip, port))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		
		transport := transportFactory.GetTransport(socket)
		if err := transport.Open(); err != nil {
			return nil, errors.WithStack(err)
		}
		
		obj.transport = transport
		
		sqlDB, err := sql.Open("spark", fmt.Sprintf("spark://%s:%s@%s:%s", 
			cfg.Username, 
			cfg.Password, 
			ip, 
			port))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		
		orm, err := gorm.Open(gorm.New(gorm.Config{
			ConnPool: sqlDB,
		}))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		obj.db = orm
		
		cache.SaveDataSourceCache(cfg.ConnectId, obj)
		return obj, nil
	}
	
	return ds, nil
}

func (this *SparkClient) Ping(ctx context.Context) (res *proto.Response, err error) {
	if this.transport != nil && this.transport.IsOpen() {
		res = proto.NewResponseNotErr()
		return
	}
	
	db, err := this.db.WithContext(ctx).DB()
	if err != nil {
		return
	}
	err = db.Ping()
	if err != nil {
		return
	}
	res = proto.NewResponseNotErr()
	return
}

func (this *SparkClient) MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error) {
	tx := this.db.Begin()
	if dbName != "" {
		tx = tx.WithContext(ctx).Exec(fmt.Sprintf("USE %s", dbName))
	}
	result := tx.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		err = result.Error
		return
	}
	rowsAffected = result.RowsAffected
	err = tx.Commit().Error
	return
}

func (this *SparkClient) MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (columns []string, list []map[string]interface{}, err error) {
	tx := this.db.Begin()
	if dbName != "" {
		tx = tx.WithContext(ctx).Exec(fmt.Sprintf("USE %s", dbName))
	}
	columns, list, err = utils.QuerySQL(tx, sql, args...)
	if err != nil {
		return nil, nil, err
	}
	return columns, list, nil
}

func (this *SparkClient) MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (data map[string]interface{}, err error) {
	_, storeRes, err := this.MysqlSelectSql(ctx, dbName, sql, args)
	if err != nil {
		return
	}
	data = map[string]interface{}{}
	if len(storeRes) > 0 {
		data = storeRes[0]
	} else {
		err = sql2.ErrNoRows
		return
	}
	return
}
