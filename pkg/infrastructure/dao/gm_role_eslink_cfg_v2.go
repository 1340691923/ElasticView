package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GmRoleEslinkCfgV2Dao struct {
	orm *orm.Gorm
}

func NewGmRoleEslinkCfgV2Dao(orm *orm.Gorm) *GmRoleEslinkCfgV2Dao {
	return &GmRoleEslinkCfgV2Dao{orm: orm}
}

func (this *GmRoleEslinkCfgV2Dao) GetRoleIdById(ctx context.Context, id int) (roleIds []int, err error) {
	err = this.orm.Raw("select role_id from gm_role_eslink_cfg_v2 where es_link_cfg_id = ?", id).Scan(&roleIds).WithContext(ctx).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return roleIds, nil
}

func (this *GmRoleEslinkCfgV2Dao) QueryByRoleID(ctx context.Context, roleId int) (list []model.GmRoleEslinkCfgV2, err error) {
	err = this.orm.Table("gm_role_eslink_cfg_v2").Where("role_id = ?", roleId).Find(&list).WithContext(ctx).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

func (this *GmRoleEslinkCfgV2Dao) DeleteByRoleID(ctx context.Context, tx *gorm.DB, roleId int) (err error) {
	err = tx.Where("role_id = ?", roleId).Delete(model.GmRoleEslinkCfgV2{}).WithContext(ctx).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *GmRoleEslinkCfgV2Dao) DeleteEslinkRoleCfgReletionByRoleID(ctx context.Context, tx *gorm.DB, roleId int) (err error) {
	SQL := `delete from eslink_role_cfg_reletion where role_cfg_id in
             (
                  select id from gm_role_eslink_cfg_v2 where role_id = ?
              )`
	err = tx.Exec(SQL, roleId).WithContext(ctx).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

func (this *GmRoleEslinkCfgV2Dao) DeleteByEsLinkCfgId(ctx context.Context, tx *gorm.DB, cfgId int) (err error) {
	err = tx.Where("es_link_cfg_id = ?", cfgId).Delete(model.GmRoleEslinkCfgV2{}).WithContext(ctx).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
