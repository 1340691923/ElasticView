package es

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/es_sdk/pkg/cache"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/services/cache_service"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"log"
	"sort"
)

type EsClientService struct {
	cfg         *config.Config
	esCache     *cache_service.EsCache
	esLinkV2Dao *dao.EsLinkV2Dao
	gmUserDao   *dao.GmUserDao
	orm         *orm.Gorm
}

func NewEsClientService(cfg *config.Config, esCache *cache_service.EsCache, esLinkV2Dao *dao.EsLinkV2Dao, gmUserDao *dao.GmUserDao, orm *orm.Gorm) *EsClientService {
	return &EsClientService{cfg: cfg, esCache: esCache, esLinkV2Dao: esLinkV2Dao, gmUserDao: gmUserDao, orm: orm}
}

func (this *EsClientService) GetEsLinkOptions(ctx context.Context, roles []int) ([]vo.EsLinkOpt, error) {

	sort.Ints(roles)

	roleids, _ := json.Marshal(roles)

	isGetCache, optList := this.esCache.EsLinkGet(string(roleids))
	if isGetCache {
		return optList, nil
	}

	type Opt struct {
		ID        int64  `db:"id"`
		Remark    string `db:"remark"`
		CfgRemark string `db:"cfg_remark"`
		Version   string `json:"version"`
	}
	var opt []Opt
	err := this.orm.Raw(`select relation_id as id,cfg_remark,remark,version from 
 (
select * from
(select es_link_id,role_cfg_id,id as relation_id from eslink_role_cfg_reletion )
 j1
left join
(
select * from
(select id, user,pwd,rootpem ,certpem,keypem,remark as cfg_remark  from eslink_cfg_v2   ) ecv
left join
(select role_id,es_link_cfg_id from gm_role_eslink_cfg_v2) grec
on ecv.id = grec.es_link_cfg_id
) j2
on j1.role_cfg_id = j2.id
)
j1 
left join
(
select id,ip,remark,version from es_link_v2 elv
) j3
on  j1.es_link_id = j3.id
where role_id in (?) and remark is not null
`, roles).Scan(&opt).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, v := range opt {
		optList = append(optList, vo.EsLinkOpt{
			ID:      v.ID,
			Remark:  fmt.Sprintf("%s(%s)", v.Remark, v.CfgRemark),
			Version: v.Version,
		})
	}

	this.esCache.EsLinkSet(string(roleids), optList)

	cache.CleanDataSourceCache(true)

	return optList, nil
}

// 获取es配置信息
func (this *EsClientService) GetEsClientByID(ctx context.Context, id int, userId int) (*model.EsConnect, error) {
	var err error

	if id == 0 {
		return nil, errors.New("请先选择数据源")
	}

	roleIds, err := this.gmUserDao.GetRolesFromUser(userId)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	linkOptions, err := this.GetEsLinkOptions(ctx, roleIds)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	isBanLink := true

	for _, link := range linkOptions {
		if link.ID == int64(id) {
			isBanLink = false
		}
	}

	if isBanLink {
		return nil, errors.New("您已经被移除访问该连接的权限")
	}

	esConnect := this.esCache.Get(id)
	if esConnect != nil {
		return esConnect, nil
	}

	esConnectTmp, err := this.esLinkV2Dao.GetEsConnectDataByRoleId(ctx, roleIds, id)

	if util.FilterMysqlNilErr(err) {
		return nil, err
	}

	if esConnectTmp.Ip == "" {
		return nil, errors.New("请先选择ES连接")
	}
	if esConnectTmp.Pwd != "" {
		esConnectTmp.Pwd, err = this.EsPwdESBDecrypt(ctx, esConnectTmp.Pwd)
		if err != nil {
			return nil, err
		}
	}
	header := []vo.HeaderKv{}

	json.Unmarshal([]byte(esConnectTmp.Header), &header)

	header = append(header, vo.HeaderKv{
		Key:   "ev_user_id",
		Value: cast.ToString(userId),
	})
	b, _ := json.Marshal(header)
	esConnectTmp.Header = string(b)
	this.esCache.Set(id, esConnectTmp)
	return esConnectTmp, nil
}

func (this *EsClientService) EsPwdESBDecrypt(ctx context.Context, cryptedStr string) (string, error) {

	pwdByte, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		log.Println(cryptedStr, this.cfg.EsPwdSecret, err)
		return "", err
	}
	b, err := util.ECBDecrypt(pwdByte, this.cfg.EsPwdSecret)
	if err != nil {
		log.Println(cryptedStr, this.cfg.EsPwdSecret, pwdByte, err)
		return "", err
	}
	return string(b), nil
}

func (this *EsClientService) EsPwdESBEncrypt(ctx context.Context, pwd string) (string, error) {
	b, err := util.ECBEncrypt(pwd, this.cfg.EsPwdSecret)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}
