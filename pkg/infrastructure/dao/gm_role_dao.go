package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GmRoleDao struct {
	orm *sqlstore.SqlStore
}

func NewGmRoleDao(orm *sqlstore.SqlStore) *GmRoleDao {
	return &GmRoleDao{orm: orm}
}

// GetById
func (this *GmRoleDao) GetById(ctx context.Context, roleId int) (role *model.GmRole, err error) {
	role = new(model.GmRole)
	err = this.orm.Raw("select id,role_name,description,role_list from gm_role where id = ?;", roleId).Scan(&role).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

// Update
func (this *GmRoleDao) Update(ctx context.Context, gmRole model.GmRole) (err error) {
	err = this.orm.Exec(
		"update gm_role set role_name = ?,description=?,role_list=? where id = ?;",
		gmRole.RoleName, gmRole.Description, gmRole.RoleList, gmRole.Id).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

// Delete
func (this *GmRoleDao) Delete(ctx context.Context, tx *gorm.DB, id int) (err error) {
	err = tx.Exec("delete from gm_role where id = ? ;", id).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

// Insert
func (this *GmRoleDao) Insert(ctx context.Context, gmRole model.GmRole) (id int64, err error) {
	err = this.orm.Raw(
		"insert into gm_role (role_name,description,role_list) values (?,?,?) RETURNING id;",
		gmRole.RoleName, gmRole.Description, gmRole.RoleList).Scan(&id).WithContext(ctx).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

// Select
func (this *GmRoleDao) Select(ctx context.Context, isAdmin bool) (model []model.GmRole, err error) {
	if isAdmin {
		err = this.orm.Raw("select role_name,description,role_list,id from gm_role;").Scan(&model).WithContext(ctx).Error
	} else {
		err = this.orm.Raw("select role_name,description,role_list,id from gm_role where id != 1;").Scan(&model).WithContext(ctx).Error
	}
	err = errors.WithStack(err)

	return
}

func (this *GmRoleDao) GetApis(ctx context.Context, roleId int) []string {
	apis := []string{}

	this.orm.Raw("select v1 from casbin_rule where v0 = ?;", roleId).Scan(&apis)
	return apis
}
