package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type EslinkCfgV2Dao struct {
	orm *sqlstore.SqlStore
}

func NewEslinkCfgV2Dao(orm *sqlstore.SqlStore) *EslinkCfgV2Dao {
	return &EslinkCfgV2Dao{orm: orm}
}

func (this *EslinkCfgV2Dao) SelectByCreateBy(ctx context.Context, userid int) ([]model.EslinkCfgV2, error) {
	var esCfgs []model.EslinkCfgV2
	err := this.orm.Raw("select * from eslink_cfg_v2 where create_by = ?", userid).Scan(&esCfgs).WithContext(ctx).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return esCfgs, nil
}

func (this *EslinkCfgV2Dao) UpdateById(ctx context.Context, tx *gorm.DB, update map[string]interface{}, id int) error {
	err := tx.Model(model.EslinkCfgV2{}).Where("id=?", id).Updates(update).WithContext(ctx).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EslinkCfgV2Dao) FindById(ctx context.Context, id int) (model.EslinkCfgV2, error) {
	var eslinkCfgV2 model.EslinkCfgV2
	err := this.orm.Model(model.EslinkCfgV2{}).Where("id=?", id).Find(&eslinkCfgV2).WithContext(ctx).Error
	if err != nil {
		return eslinkCfgV2, errors.WithStack(err)
	}
	return eslinkCfgV2, nil
}

func (this *EslinkCfgV2Dao) Save(ctx context.Context, data *model.EslinkCfgV2) error {
	err := this.orm.Save(data).WithContext(ctx).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
