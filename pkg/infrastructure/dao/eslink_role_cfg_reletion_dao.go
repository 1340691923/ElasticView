package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type EslinkRoleCfgReletionDao struct {
	orm *orm.Gorm
}

func NewEslinkRoleCfgReletion(orm *orm.Gorm) *EslinkRoleCfgReletionDao {
	return &EslinkRoleCfgReletionDao{orm: orm}
}

func (this *EslinkRoleCfgReletionDao) DeleteEsCfgRelation(ctx context.Context, tx *gorm.DB, relationId int) (err error) {

	err = tx.Debug().
		Where("id = ?", relationId).
		Delete(model.EslinkRoleCfgReletion{}).
		WithContext(ctx).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (this *EslinkRoleCfgReletionDao) DeleteByEsLinkId(ctx context.Context, tx *gorm.DB, id int) (err error) {
	err = tx.
		Where("es_link_id = ?", id).
		Delete(model.EslinkRoleCfgReletion{}).WithContext(ctx).
		Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (this *EslinkRoleCfgReletionDao) UpdateByEslLinkId(ctx context.Context, roleCfgId, esLinkID int) error {
	err := this.orm.
		Model(model.EslinkRoleCfgReletion{}).
		Where("es_link_id = ?", esLinkID).
		Update("role_cfg_id", roleCfgId).WithContext(ctx).
		Error
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return err
}

func (this *EslinkRoleCfgReletionDao) Save(ctx context.Context, tx *gorm.DB, roleCfgId, esLinkID int) error {
	err := tx.Create(&model.EslinkRoleCfgReletion{
		EsLinkId:  esLinkID,
		RoleCfgId: roleCfgId,
		Created:   time.Now(),
		Updated:   time.Now(),
	}).WithContext(ctx).Error
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return err
}
