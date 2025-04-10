package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type EslinkCfgV2Dao struct {
	orm *orm.Gorm
}

func NewEslinkCfgV2Dao(orm *orm.Gorm) *EslinkCfgV2Dao {
	return &EslinkCfgV2Dao{orm: orm}
}

func (this *EslinkCfgV2Dao) SelectByCreateBy(ctx context.Context, userid int, page, pageSize int) (esCfgs []model.EslinkCfgV2, count int, err error) {

	db := this.orm.WithContext(ctx).Table("eslink_cfg_v2")

	db = db.Where("create_by = ?", userid)

	err = db.Select("count(*)").Scan(&count).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = db.Select("*").Limit(pageSize).Offset(orm.CreatePage(page, pageSize)).Scan(&esCfgs).Error

	err = errors.WithStack(err)

	return
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
