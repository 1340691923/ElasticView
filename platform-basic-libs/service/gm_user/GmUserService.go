//GM用户层
package gm_user

import (
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
)

// GmUserService
type GmUserService struct {
}

func (this GmUserService) CheckLogin(username, password string) (token string, err error) {
	var model2 model.GmUserModel
	model2.Password = password
	model2.Username = username
	gmUser, err := model2.GetUserByUP()
	if err != nil {
		return
	}
	token, err = jwt.GenerateToken(gmUser)
	if err != nil {
		return
	}
	return
}

func (this GmUserService) GetRoleInfo(roleId int32) (gminfo model.GmRoleModel, err error) {
	var model2 model.GmRoleModel
	gminfo, err = model2.GetById(int(roleId))
	if err != nil {
		return
	}
	return
}

func (this GmUserService) IsExitUser(claims *jwt.Claims) bool {
	var model2 model.GmUserModel
	model2.Username = claims.Username
	model2.RoleId = claims.RoleId
	return model2.Exsit()
}
