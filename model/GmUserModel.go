package model

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"

	"go.uber.org/zap"
)

// GmUserModel GM用户
type GmUserModel struct {
	ID       int32  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	RoleId   int32  `json:"role_id" db:"role_id"`
	Realname string `json:"realname" db:"realname"`
}

//密码进行md5混淆
func (this GmUserModel) GetPassword() string {
	return util.MD5HexHash([]byte(this.Password))
}

//是否存在该用户
func (this GmUserModel) Exsit() (b bool) {
	var count int
	err := db.Sqlx.Get(&count, "select count(*) from gm_user where username = ? and role_id = ? limit 1;", this.Username, this.RoleId)
	if err != nil || count == 0 {
		logs.Logger.Error("err", zap.String("err.Error()", err.Error()))
		return false
	}
	return true
}

//登录
func (this GmUserModel) GetUserByUP() (gmUser GmUserModel, err error) {
	err = db.Sqlx.Get(&gmUser, "select id,username,password,role_id,realname from gm_user where username = ? and password = ? limit 1;", this.Username, this.GetPassword())
	return
}

//通过id查询用户信息
func (this GmUserModel) GetUserById() (gmUser GmUserModel, err error) {
	err = db.Sqlx.Get(&gmUser, "select id,username,password,role_id,realname from gm_user where id = ?;", this.ID)
	return
}

//新增用户
func (this GmUserModel) Insert() (id int64, err error) {
	rlt, err := db.Sqlx.Exec("insert into gm_user(username,password,role_id,realname)values(?,?,?,?)", this.Username, this.GetPassword(), this.RoleId, this.Realname)
	if err != nil {
		return
	}
	id, _ = rlt.LastInsertId()
	return
}

// Update
func (this GmUserModel) Update() (err error) {
	_, err = db.Sqlx.Exec("update gm_user set username = ?,password=?,role_id=?,realname=? where id = ? ;", this.Username, this.GetPassword(), this.RoleId, this.Realname, this.ID)
	return
}

// Select
func (this GmUserModel) Select() (gmUser []GmUserModel, err error) {
	err = db.Sqlx.Select(&gmUser, "select id,username,password,role_id,realname from gm_user ;")
	return
}

// Delete
func (this GmUserModel) Delete() (err error) {
	_, err = db.Sqlx.Exec("delete from gm_user where id = ? ;", this.ID)
	return
}
