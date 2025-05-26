package dao

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// GmUserModel
type GmUserDao struct {
	orm *orm.Gorm
}

func NewGmUserDao(orm *orm.Gorm) *GmUserDao {
	return &GmUserDao{orm: orm}
}

// 是否存在该用户
func (this *GmUserDao) Exsit(ctx context.Context, user model.GmUserModel) (exsit bool, err error) {
	var count int

	err = this.orm.Raw("select count(*) from gm_user where username = ? and is_ban = 0  limit 1;", user.Username).Scan(&count).WithContext(ctx).Error

	err = errors.WithStack(err)

	exsit = count > 0
	return
}

// 登录
func (this *GmUserDao) GetUserByUP(ctx context.Context, user model.GmUserModel) (gmUser model.GmUserModel, err error) {
	err = this.orm.Raw("select id,username,password,realname,is_ban from gm_user where username = ? and password = ? limit 1;", user.Username, user.GetPassword()).WithContext(ctx).Scan(&gmUser).Error

	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (this *GmUserDao) UpdateLastLoginTime(ctx context.Context, userid int) (err error) {
	err = this.orm.Model(model.GmUserModel{}).Where("id = ?", userid).
		Updates(map[string]interface{}{"last_login_time": time.Now().Format(util.TimeFormat)}).
		WithContext(ctx).
		Error
	err = errors.WithStack(err)
	return
}

// 获取用户信息
func (this *GmUserDao) GetUserByUserName(ctx context.Context, userName string) (gmUser model.GmUserModel, err error) {
	err = this.orm.Raw("select id,username,password,realname,is_ban from gm_user where username = ?  limit 1;", userName).WithContext(ctx).Scan(&gmUser).Error
	err = errors.WithStack(err)
	return
}

// 通过id查询用户信息
func (this *GmUserDao) GetUserById(ctx context.Context, id int) (gmUser model.GmUserModel, err error) {
	err = this.orm.Raw("select id,username,password,realname,is_ban from gm_user where id = ?;", id).WithContext(ctx).Scan(&gmUser).Error
	err = errors.WithStack(err)
	return
}

// 新增用户
func (this *GmUserDao) Insert(ctx context.Context, orm *gorm.DB, gmUser model.GmUserModel) (id int64, err error) {

	if gmUser.Password == "" {
		gmUser.Password = util.GetUUid()
	}

	data := &model.GmUserModel{
		Username:      gmUser.Username,
		Password:      gmUser.GetPassword(),
		Avatar:        gmUser.Avatar,
		Realname:      gmUser.Realname,
		Email:         gmUser.Email,
		WorkWechatUid: gmUser.WorkWechatUid,
		UpdateTime:    time.Now(),
		CreateTime:    time.Now(),
		LastLoginTime: time.Now(),
		IsBan:         0,
	}

	err = orm.WithContext(ctx).Create(data).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	id = int64(data.Id)
	return
}

// Update
func (this *GmUserDao) Update(ctx context.Context, gmUser model.GmUserModel) (err error) {

	err = this.orm.Exec("update gm_user set "+
		"username = ?,realname=?,avatar=?,email =?,work_wechat_uid=?  where id = ? ;",
		gmUser.Username, gmUser.Realname, gmUser.Avatar, gmUser.Email, gmUser.WorkWechatUid,
		gmUser.Id).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

type SealType int

const (
	Ban   SealType = 1
	UnBan SealType = 0
)

// Update
func (this *GmUserDao) SealUser(ctx context.Context, id int, sealType SealType) (err error) {

	err = this.orm.Exec("update gm_user set is_ban=? where id = ? ;", sealType, id).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

// Update
func (this *GmUserDao) UpdatePassById(ctx context.Context, gmUser model.GmUserModel) (err error) {

	err = this.orm.Exec("update gm_user set password=? where id = ? ;", gmUser.GetPassword(), gmUser.Id).WithContext(ctx).Error
	err = errors.WithStack(err)
	return
}

// Select
func (this *GmUserDao) Select(ctx context.Context, isAdmin bool, userName,
	realName string, isBan bool, userIds []int, page, pageSize int) (gmUser []model.GmUserModel, count int64, err error) {
	db := this.orm.WithContext(ctx).Table("gm_user")
	if !isAdmin {
		db = db.Where(" id != 1")
	}

	if len(userIds) > 0 {
		db = db.Where(" id in ? ", userIds)
	}

	if userName != "" {
		db = db.Where(" username like ? ", "%"+userName+"%")
	}

	if realName != "" {
		db = db.Where(" realname like ? ", "%"+realName+"%")
	}

	if isBan {
		db = db.Where(" is_ban = 1 ")
	} else {
		db = db.Where(" is_ban = 0 ")
	}

	err = db.Select("count(*)").Scan(&count).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = db.Select("*").Limit(pageSize).Offset(orm.CreatePage(page, pageSize)).Scan(&gmUser).Error

	err = errors.WithStack(err)

	return
}

// Delete
func (this *GmUserDao) Delete(ctx context.Context, tx *gorm.DB, id int) (err error) {
	err = tx.Where("id = ?", id).Delete(&model.GmUserModel{}).WithContext(ctx).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

func (this *GmUserDao) GetByField(field string, value interface{}) (gmUser model.GmUserModel, err error) {
	err = this.orm.Model(&model.GmUserModel{}).Where(fmt.Sprintf(" %s = ? ", field), value).Scan(&gmUser).Error
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return gmUser, nil
}

func (this *GmUserDao) GetRolesFromUser(userID int) ([]int, error) {
	type Role struct {
		RoleId int `grom:"role_id"`
	}
	var roles []Role
	roleIds := []int{}

	err := this.orm.
		Raw(`select role_id from user_role_relation where user_id = ?`,
			userID).Scan(&roles).Error
	if err != nil {
		return roleIds, errors.WithStack(err)
	}
	for _, v := range roles {
		roleIds = append(roleIds, v.RoleId)
	}
	return roleIds, nil
}

func (this *GmUserDao) RemoveUserRoles(tx *gorm.DB, userID int) (err error) {
	return errors.WithStack(tx.Where("user_id = ?", userID).Delete(model.UserRoleRelationModel{}).Error)
}

func (this *GmUserDao) GetUserByRoleId(tx *gorm.DB, roleId []int) (users []model.UserRoleRelationModel, err error) {
	err = errors.WithStack(
		tx.Where("role_id in ?", roleId).Find(&users).Error,
	)
	return
}

func (this *GmUserDao) RemoveRoles(tx *gorm.DB, roleId int) (err error) {
	return errors.WithStack(tx.Where("role_id = ?", roleId).Delete(model.UserRoleRelationModel{}).Error)
}

func (this *GmUserDao) AddRolesToUser(tx *gorm.DB, userID int, roles []int) (err error) {

	err = this.RemoveUserRoles(tx, userID)
	if err != nil {
		return errors.WithStack(err)
	}

	userRoleRelations := []model.UserRoleRelationModel{}

	for _, roleId := range roles {

		userRoleRelations = append(userRoleRelations, model.UserRoleRelationModel{
			UserId:     userID,
			RoleId:     roleId,
			UpdateTime: time.Now(),
			CreateTime: time.Now(),
		})

	}
	err = tx.Create(&userRoleRelations).Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
