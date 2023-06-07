package escache

import (
	"encoding/base64"
	"errors"
	"github.com/1340691923/ElasticView/pkg/engine/config"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/util"
)

// 获取es配置信息
func GetEsClientByID(id int) (*EsConnect, error) {
	var err error
	esCache := NewEsCache()
	esConnect := esCache.Get(id)
	if esConnect != nil {
		return esConnect, nil
	}
	esConnectTmp := EsConnect{}
	sql, args, err := db.SqlBuilder.
		Select("ip", "user", "pwd", "version", "rootpem", "certpem", "keypem").
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
	if esConnectTmp.Pwd != "" {
		esConnectTmp.Pwd, err = EsPwdESBDecrypt(esConnectTmp.Pwd)
		if err != nil {
			return nil, err
		}
	}
	esCache.Set(id, &esConnectTmp)
	return &esConnectTmp, nil
}

func EsPwdESBDecrypt(cryptedStr string) (string, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		return "", err
	}
	b, err := util.ECBDecrypt(pwdByte, config.GlobConfig.EsPwdSecret)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func EsPwdESBEncrypt(pwd string) (string, error) {
	b, err := util.ECBEncrypt(pwd, config.GlobConfig.EsPwdSecret)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
