package data_conversion

import (
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/request"
)

type Datasource interface {
	Ping() error
	GetTables() ([]string, error)
	GetTableColumns(tableName string) (interface{}, error)
	Transfer(id int, transferReq *request.TransferReq) (err error)
}

const (
	MysqlSource = "mysql"
	//MssqlSource = "sqlserver"
	CkSource = "clickhouse"
	//MongoSource = "mongodb"
)

var DataSourceMap = map[string]func(data request.DataxInfoTestLinkReq) Datasource{
	MysqlSource: NewMysql,
	//MssqlSource: NewMssql,
	CkSource: NewClickhouse,
	//MongoSource: NewMongoDb,
}

func NewDataSource(data request.DataxInfoTestLinkReq) (i Datasource, err error) {
	var fn func(data request.DataxInfoTestLinkReq) Datasource
	var found bool
	if fn, found = DataSourceMap[data.Typ]; !found {
		return nil, errors.New(fmt.Sprintf("没有找到该命令:%v", data.Typ))
	}
	return fn(data), nil
}
