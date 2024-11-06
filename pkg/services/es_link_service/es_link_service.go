package es_link_service

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/services/es"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// EsLinkModel es连接信息表
type EsLinkService struct {
	orm                      *sqlstore.SqlStore
	logger                   *logger.AppLogger
	esClientService          *es.EsClientService
	eslinkCfgV2Dao           *dao.EslinkCfgV2Dao
	gmRoleEslinkCfgV2Dao     *dao.GmRoleEslinkCfgV2Dao
	eslinkRoleCfgReletionDao *dao.EslinkRoleCfgReletionDao
	gmUserDao                *dao.GmUserDao
	esLinkV2Dao              *dao.EsLinkV2Dao
}

func NewEsLinkService(orm *sqlstore.SqlStore, logger *logger.AppLogger, esClientService *es.EsClientService, eslinkCfgV2Dao *dao.EslinkCfgV2Dao, gmRoleEslinkCfgV2Dao *dao.GmRoleEslinkCfgV2Dao, eslinkRoleCfgReletionDao *dao.EslinkRoleCfgReletionDao, gmUserDao *dao.GmUserDao, esLinkV2Dao *dao.EsLinkV2Dao) *EsLinkService {
	return &EsLinkService{orm: orm, logger: logger, esClientService: esClientService, eslinkCfgV2Dao: eslinkCfgV2Dao, gmRoleEslinkCfgV2Dao: gmRoleEslinkCfgV2Dao, eslinkRoleCfgReletionDao: eslinkRoleCfgReletionDao, gmUserDao: gmUserDao, esLinkV2Dao: esLinkV2Dao}
}

// 获取列表信息
func (this *EsLinkService) GetListAction(ctx context.Context, userID int, roleID []int) (esLinkList []*vo.EsLink, err error) {
	isAdmin := false
	if util.InArr(roleID, 1) {
		isAdmin = true
	}
	builder := this.orm.Model(model.EsLinkV2{})

	if !isAdmin {
		builder = builder.Where("create_by = ?", userID)
	}

	var esLinks []model.EsLinkV2

	err = builder.Scan(&esLinks).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	for _, esLink := range esLinks {

		userInfo, err := this.gmUserDao.GetUserById(ctx, esLink.CreateBy)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		esLinkVo := &vo.EsLink{
			ID:               int(esLink.Id),
			Remark:           esLink.Remark,
			Ip:               esLink.Ip,
			Version:          esLink.Version,
			CreateById:       esLink.CreateBy,
			CreateByUserName: userInfo.Username,
			Created:          esLink.Created.Format(util.TimeFormat),
			Updated:          esLink.Updated.Format(util.TimeFormat),
			EsLinkConfigs:    nil,
		}

		var esCfgs []model.EslinkCfgModelTmp
		err = this.orm.Raw(`select relation_id,es_link_id,role_cfg_id,id,user,pwd,rootpem,certpem,keypem,remark,header from 
(select es_link_id,role_cfg_id,id as relation_id from eslink_role_cfg_reletion ) j1
left join 
(
select id, user,pwd,rootpem ,certpem,keypem,remark,header   from eslink_cfg_v2 
) j2
on j1.role_cfg_id = j2.id 
where es_link_id = ? `, esLink.Id).Scan(&esCfgs).Error
		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, esCfg := range esCfgs {
			if esCfg.Pwd != "" {
				esCfg.Pwd, err = this.esClientService.EsPwdESBDecrypt(ctx, esCfg.Pwd)
				if err != nil {
					return nil, errors.WithStack(err)
				}
			}

			header := []vo.HeaderKv{}
			if esCfg.Header == "" {
				esCfg.Header = "[]"
			}
			err = json.Unmarshal([]byte(esCfg.Header), &header)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			esLinkConfig := &vo.EsLinkConfig{
				CfgRelationId: esCfg.RelationId,
				Id:            esCfg.ID,
				Ip:            esLink.Ip,
				Version:       esLink.Version,
				EsLinkId:      esCfg.EsLinkId,
				User:          esCfg.User,
				Pwd:           esCfg.Pwd,
				Remark:        esCfg.Remark,
				Created:       esCfg.Created.Format(util.TimeFormat),
				Updated:       esCfg.Updated.Format(util.TimeFormat),
				RootPEM:       esCfg.RootPEM,
				CertPEM:       esCfg.CertPEM,
				KeyPEM:        esCfg.KeyPEM,
				Header:        header,
			}
			sharedRoleIds := []string{}
			type SharedRole struct {
				RoleId int `db:"role_id"`
			}
			sharedRoles, err := this.gmRoleEslinkCfgV2Dao.GetRoleIdById(ctx, esLinkConfig.Id)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			for _, v := range sharedRoles {
				sharedRoleIds = append(sharedRoleIds, cast.ToString(v))
			}

			esLinkConfig.ShareRoles = sharedRoleIds
			esLinkVo.EsLinkConfigs = append(esLinkVo.EsLinkConfigs, esLinkConfig)
		}
		esLinkList = append(esLinkList, esLinkVo)
	}

	return esLinkList, nil
}

func (this *EsLinkService) GetEsCfgList(ctx context.Context, userId int) (res []*vo.EsLinkConfigV2, err error) {
	var esCfgs []model.EslinkCfgV2
	esCfgs, err = this.eslinkCfgV2Dao.SelectByCreateBy(ctx, userId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, esCfg := range esCfgs {
		if esCfg.Pwd != "" {
			esCfg.Pwd, err = this.esClientService.EsPwdESBDecrypt(ctx, esCfg.Pwd)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}

		headerString := esCfg.Header

		if headerString == "" {
			headerString = "[]"
		}
		header := []vo.HeaderKv{}
		json.Unmarshal([]byte(headerString), &header)

		esLinkConfig := &vo.EsLinkConfigV2{
			Id:      esCfg.Id,
			User:    esCfg.User,
			Pwd:     esCfg.Pwd,
			Remark:  esCfg.Remark,
			Created: esCfg.Created.Format(util.TimeFormat),
			Updated: esCfg.Updated.Format(util.TimeFormat),
			RootPEM: cast.ToString(esCfg.Rootpem),
			CertPEM: cast.ToString(esCfg.Certpem),
			KeyPEM:  cast.ToString(esCfg.Keypem),
			Header:  header,
		}
		sharedRoleIds := []string{}

		roles, err := this.gmRoleEslinkCfgV2Dao.GetRoleIdById(ctx, esCfg.Id)

		if err != nil {
			return nil, errors.WithStack(err)
		}

		for _, roleid := range roles {
			sharedRoleIds = append(sharedRoleIds, cast.ToString(roleid))
		}

		esLinkConfig.ShareRoles = sharedRoleIds

		res = append(res, esLinkConfig)
	}
	return res, nil
}

func (this *EsLinkService) GetEsCfgOpt(ctx context.Context, userId int) (res []*vo.EsLinkConfigOpt, err error) {
	var esCfgs []model.EslinkCfgV2

	esCfgs, err = this.eslinkCfgV2Dao.SelectByCreateBy(ctx, userId)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, esCfg := range esCfgs {

		esLinkConfig := &vo.EsLinkConfigOpt{
			Id:     esCfg.Id,
			Remark: esCfg.Remark,
		}

		res = append(res, esLinkConfig)
	}

	return res, nil
}

func (this *EsLinkService) DeleteEsCfgRelation(ctx context.Context, tx *gorm.DB, relationId int) (err error) {
	err = this.eslinkRoleCfgReletionDao.DeleteEsCfgRelation(ctx, tx, relationId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EsLinkService) DeleteRoleEslinkCfgByRoleId(ctx context.Context, tx *gorm.DB, roleId int) (err error) {

	err = this.gmRoleEslinkCfgV2Dao.DeleteEslinkRoleCfgReletionByRoleID(ctx, tx, roleId)
	if err != nil {
		return errors.WithStack(err)
	}
	err = this.gmRoleEslinkCfgV2Dao.DeleteByRoleID(ctx, tx, roleId)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EsLinkService) DeleteRoleEslinkCfgByEsLinkCfgId(ctx context.Context, tx *gorm.DB, roleId int) (err error) {
	return errors.WithStack(this.gmRoleEslinkCfgV2Dao.DeleteByEsLinkCfgId(ctx, tx, roleId))
}

func (this *EsLinkService) UpdateEsLinkCfgById(ctx context.Context, tx *gorm.DB, updateMap map[string]interface{}, id int) error {
	return errors.WithStack(this.eslinkCfgV2Dao.UpdateById(ctx, tx, updateMap, id))
}

func (this *EsLinkService) DeleteEsCfgRelationByEsLinkId(ctx context.Context, tx *gorm.DB, eslinkId int) (err error) {

	return errors.WithStack(this.eslinkRoleCfgReletionDao.DeleteByEsLinkId(ctx, tx, eslinkId))
}

func (this *EsLinkService) DeleteById(ctx context.Context, tx *gorm.DB, id int) (err error) {
	return errors.WithStack(this.esLinkV2Dao.Delete(ctx, tx, id))
}

func (this *EsLinkService) DeleteEsCfg(ctx context.Context, id int) (err error) {
	tx := this.orm.Begin()

	err = tx.Where("id = ?", id).Delete(model.EslinkCfgV2{}).WithContext(ctx).Error
	if err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	err = tx.Where("es_link_cfg_id = ?", id).Delete(model.GmRoleEslinkCfgV2{}).WithContext(ctx).Error
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Where("role_cfg_id = ?", id).Delete(model.EslinkRoleCfgReletion{}).WithContext(ctx).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()

	return

}

func (this *EsLinkService) SaveEsLink(ctx context.Context, tx *gorm.DB, data *model.EsLinkV2) (id int, err error) {
	err = this.esLinkV2Dao.Save(ctx, tx, data)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return data.Id, nil
}

func (this *EsLinkService) UpdateEsLink(ctx context.Context, tx *gorm.DB, data map[string]interface{}, id int) (err error) {
	err = tx.Model(model.EsLinkV2{}).Where("id = ?", id).WithContext(ctx).Updates(data).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EsLinkService) SaveEsLinkCfgV2(ctx context.Context, data *model.EslinkCfgV2) (id int, err error) {
	err = this.eslinkCfgV2Dao.Save(ctx, data)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return data.Id, nil
}

func (this *EsLinkService) UpdateEslinkRoleCfgByEsLinkId(ctx context.Context, roleCfgId, esLinkID int) (err error) {
	err = this.eslinkRoleCfgReletionDao.UpdateByEslLinkId(ctx, roleCfgId, esLinkID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EsLinkService) SaveEslinkRoleCfgByEsLinkId(ctx context.Context, tx *gorm.DB, roleCfgId, esLinkID int) (err error) {
	err = this.eslinkRoleCfgReletionDao.Save(ctx, tx, roleCfgId, esLinkID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
