package model

import "github.com/1340691923/ElasticView/engine/db"

// GmRoleModel
type GmRoleModel struct {
	ID          int    `json:"id" db:"id"`
	RoleName    string `json:"name" db:"role_name"`
	Description string `json:"description" db:"description"`
	RoleList    string `json:"routes" db:"role_list"`
}

// GetById
func (this *GmRoleModel) GetById(roleId int) (model GmRoleModel, err error) {
	err = db.Sqlx.Get(&model, "select id,role_name,description,role_list from gm_role where id = ?;", roleId)
	return
}

// Update
func (this *GmRoleModel) Update() (err error) {
	_, err = db.Sqlx.Exec("update gm_role set role_name = ?,description=?,role_list=? where id = ?;", this.RoleName, this.Description, this.RoleList, this.ID)
	return
}

// Delete
func (this *GmRoleModel) Delete() (err error) {
	_, err = db.Sqlx.Exec("delete from gm_role where id = ? ;", this.ID)
	return
}

// Insert
func (this *GmRoleModel) Insert() (id int64, err error) {
	rlt, err := db.Sqlx.Exec("insert into gm_role (role_name,description,role_list) values (?,?,?);", this.RoleName, this.Description, this.RoleList)
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
func (this *GmRoleModel) Select() (model []GmRoleModel, err error) {
	err = db.Sqlx.Select(&model, "select role_name,description,role_list,id from gm_role;")
	return
}
