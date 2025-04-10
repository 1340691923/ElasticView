package mysql

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/base"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	proto2 "github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/proto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/utils"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/pkg"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlClient struct {
	base.BaseDatasource
	db *gorm.DB
}

func NewMysqlClient(cfg *proto2.Config) (pkg.ClientInterface, error) {
	ds, ok := cache.GetDataSourceCache(cfg.ConnectId)

	if !ok {
		obj := &MysqlClient{}

		if len(cfg.Addresses) == 0 {
			return nil, errors.New("ip和端口不能为空")
		}
		// 打开数据库连接
		orm, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true",
			cfg.Username,
			cfg.Password,
			cfg.Addresses[0],
		)))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		obj.db = orm

		cache.SaveDataSourceCache(cfg.ConnectId, obj)
		return obj, nil
	}

	return ds, nil
}

func (this *MysqlClient) Ping(
	ctx context.Context,
) (
	res *proto.Response,
	err error,
) {
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

// todo... 表名鉴权 任何插件表名必须以 plugin_ 开头
func (this *MysqlClient) MysqlExecSql(ctx context.Context, dbName, sql string, args ...interface{}) (rowsAffected int64, err error) {

	tx := this.db.Begin() // start transaction

	if dbName != "" {
		tx = tx.WithContext(ctx).Exec(fmt.Sprintf("use %s", dbName))
	}

	result := tx.WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		err = result.Error
		return
	}
	rowsAffected = result.RowsAffected

	err = tx.Commit().Error // end transaction

	return
}

// todo... 表名鉴权
func (this *MysqlClient) MysqlSelectSql(ctx context.Context, dbName, sql string, args ...interface{}) (columns []string, list []map[string]interface{}, err error) {

	tx := this.db.Begin()

	if dbName != "" {
		tx = tx.WithContext(ctx).Exec(fmt.Sprintf("use %s", dbName))
	}

	columns, list, err = utils.QuerySQL(tx, sql, args...)

	if err != nil {
		return nil, nil, err
	}

	return columns, list, nil
}

func (this *MysqlClient) MysqlFirstSql(ctx context.Context, dbName, sql string, args ...interface{}) (data map[string]interface{}, err error) {

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
