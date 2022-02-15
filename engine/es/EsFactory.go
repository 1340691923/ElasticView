package es

import (
	"errors"

	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
)

//获取es6.x的链接
func GetEsClientV6ByID(id int) (esClient EsClient, err error) {

	esCache := NewEsCache()
	esClient = esCache.Get(id)
	if esClient != nil {
		return esClient, nil
	}

	esConnect := EsConnect{}
	sql, args, err := db.SqlBuilder.
		Select("ip", "user", "pwd", "version").
		From("es_link").
		Where(db.Eq{"id": id}).
		Limit(1).ToSql()
	if err != nil {
		return
	}
	err = db.Sqlx.Get(&esConnect, sql, args...)

	if util.FilterMysqlNilErr(err) {
		return
	}
	if esConnect.Ip == "" {
		return nil, errors.New("请先选择ES连接")
	}

	esClient, err = NewEsClientV6(esConnect)
	if err != nil {
		return
	}
	esCache.Set(id, esClient)
	return esClient, nil
}

//根据版本号获取示例
func GetEsClient(esConnect EsConnect) (esClient EsClient, err error) {

	if esConnect.Ip == "" {
		return nil, errors.New("请先选择ES连接")
	}
	switch esConnect.Version {
	case 6:
		return NewEsClientV6(esConnect)
	case 7:
		return NewEsClientV7(esConnect)
	default:
		return nil, errors.New("无效的版本号")
	}
}
