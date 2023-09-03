package model

import "github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"

// GmRoleModel
type GmRoleModel struct {
	ID          int    `json:"id" db:"id"`
	RoleName    string `json:"name" db:"role_name"`
	Description string `json:"description" db:"description"`
	RoleList    string `json:"routes" db:"role_list"`
}

// GetById
func (this *GmRoleModel) GetById(sqlx *sqlstore.SqlStore, roleId int) (model GmRoleModel, err error) {
	err = sqlx.Get(&model, "select id,role_name,description,role_list from gm_role where id = ?;", roleId)
	return
}

// Update
func (this *GmRoleModel) Update(sqlx *sqlstore.SqlStore) (err error) {
	_, err = sqlx.Exec(
		"update gm_role set role_name = ?,description=?,role_list=? where id = ?;",
		this.RoleName, this.Description, this.RoleList, this.ID)
	return
}

// Delete
func (this *GmRoleModel) Delete(sqlx *sqlstore.SqlStore) (err error) {
	_, err = sqlx.Exec("delete from gm_role where id = ? ;", this.ID)
	return
}

// Insert
func (this *GmRoleModel) Insert(sqlx *sqlstore.SqlStore) (id int64, err error) {
	rlt, err := sqlx.Exec(
		"insert into gm_role (role_name,description,role_list) values (?,?,?);",
		this.RoleName, this.Description, this.RoleList)
	if err != nil {
		return
	}
	id, err = rlt.LastInsertId()
	if err != nil {
		return
	}
	return
}

// Select
func (this *GmRoleModel) Select(sqlx *sqlstore.SqlStore) (model []GmRoleModel, err error) {
	err = sqlx.Select(&model, "select role_name,description,role_list,id from gm_role;")
	return
}
