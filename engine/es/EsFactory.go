package es

import (
	"errors"

	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
)

//获取es配置信息
func GetEsClientByID(id int) (*EsConnect, error) {
	var err error
	esCache := NewEsCache()
	esConnect := esCache.Get(id)
	if esConnect != nil {
		return esConnect, nil
	}
	esConnectTmp := EsConnect{}
	sql, args, err := db.SqlBuilder.
		Select("ip", "user", "pwd", "version").
		From("es_link").
		Where(db.Eq{"id": id}).
		Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	err = db.Sqlx.Get(&esConnectTmp, sql, args...)

	if util.FilterMysqlNilErr(err) {
		return nil, err
	}
	if esConnectTmp.Ip == "" {
		return nil, errors.New("请先选择ES连接")
	}

	esCache.Set(id, &esConnectTmp)
	return &esConnectTmp, nil
}
