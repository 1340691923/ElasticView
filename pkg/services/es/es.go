package es

import (
	"encoding/base64"
	"errors"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/util"
)

type EsClientService struct {
	esCache  *EsCache
	cfg      *config.Config
	sqlStore *sqlstore.SqlStore
}

func NewEsClientService(esCache *EsCache, cfg *config.Config, sqlStore *sqlstore.SqlStore) *EsClientService {
	return &EsClientService{esCache: esCache, cfg: cfg, sqlStore: sqlStore}
}

// 获取es配置信息
func (this *EsClientService) GetEsClientByID(id int) (*model.EsConnect, error) {
	var err error
	esConnect := this.esCache.Get(id)
	if esConnect != nil {
		return esConnect, nil
	}
	esConnectTmp := model.EsConnect{}
	sql, args, err := sqlstore.SqlBuilder.
		Select("ip", "user", "pwd", "version", "rootpem", "certpem", "keypem").
		From("es_link").
		Where(sqlstore.Eq{"id": id}).
		Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	err = this.sqlStore.Get(&esConnectTmp, sql, args...)

	if util.FilterMysqlNilErr(err) {
		return nil, err
	}
	if esConnectTmp.Ip == "" {
		return nil, errors.New("请先选择ES连接")
	}
	if esConnectTmp.Pwd != "" {
		esConnectTmp.Pwd, err = this.EsPwdESBDecrypt(esConnectTmp.Pwd)
		if err != nil {
			return nil, err
		}
	}
	esCache.Set(id, &esConnectTmp)
	return &esConnectTmp, nil
}

func (this *EsClientService) EsPwdESBDecrypt(cryptedStr string) (string, error) {
	pwdByte, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		return "", err
	}
	b, err := util.ECBDecrypt(pwdByte, this.cfg.EsPwdSecret)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (this *EsClientService) EsPwdESBEncrypt(pwd string) (string, error) {
	b, err := util.ECBEncrypt(pwd, this.cfg.EsPwdSecret)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
