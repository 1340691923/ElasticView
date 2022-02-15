package controller

import (
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/gm_user"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	. "github.com/gofiber/fiber/v2"
)

//GM用户控制器
type UserController struct {
	BaseController
}

// 登录
func (this UserController) Login(ctx *Ctx) error {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var user User
	err := ctx.BodyParser(&user)
	if err != nil {
		logs.Logger.Sugar().Errorf("登录失败", err)
		err = my_error.NewBusiness(gm_user.AUTH_ERROR, gm_user.ERROR_AUTH)
		return this.Error(ctx, err)
	}
	username := user.Username
	password := user.Password

	var gmUserService gm_user.GmUserService
	token, err := gmUserService.CheckLogin(username, password)
	if err != nil {
		logs.Logger.Sugar().Errorf("登录失败", err)
		err = my_error.NewBusiness(gm_user.AUTH_ERROR, gm_user.ERROR_AUTH)
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "登录成功", map[string]interface{}{"token": token})
}

// 用户详细信息
func (this UserController) UserInfo(ctx *Ctx) error {
	var gmUserService gm_user.GmUserService
	token := this.GetToken(ctx)
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return this.Error(ctx, err)
	}
	info, err := gmUserService.GetRoleInfo(claims.RoleId)
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, "登录成功", map[string]interface{}{"roles": []string{"admin"}, "introduction": info.Description, "name": info.RoleName, "list": info.RoleList, "avatar": ""})
}

//退出登录
func (this UserController) LogoutAction(ctx *Ctx) error {
	token := this.GetToken(ctx)
	var claims *jwt.Claims
	claims, err := jwt.ParseToken(token)
	if err != nil {
		logs.Logger.Sugar().Errorf("LogoutAction err", err)
		return this.Success(ctx, response.LogoutSuccess, nil)
	}
	util.TokenBucket.LoadOrStore(token, claims.ExpiresAt)

	return this.Success(ctx, response.LogoutSuccess, nil)
}

//GM 用户列表
func (this UserController) UserListAction(ctx *Ctx) error {
	var userModel model.GmUserModel
	list, err := userModel.Select()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}

// 删除GM用户
func (this UserController) DeleteUserAction(ctx *Ctx) error {
	var userModel model.GmUserModel
	userModel.ID = int32(this.FormIntDefault(ctx, "id", 0))
	err := userModel.Delete()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.DeleteSuccess, nil)
}

// 用ID获取用户信息
func (this UserController) GetUserByIdAction(ctx *Ctx) error {
	var userModel model.GmUserModel
	var id = int32(this.FormIntDefault(ctx, "id", 0))
	userModel.ID = id
	gmUser, err := userModel.GetUserById()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, gmUser)
}

// 修改GM用户信息
func (this UserController) UserUpdateAction(ctx *Ctx) error {
	var userModel model.GmUserModel
	var id = int32(this.FormIntDefault(ctx, "id", 0))
	userModel.ID = id
	userModel.Realname = ctx.FormValue("realname")
	userModel.RoleId = int32(this.FormIntDefault(ctx, "role_id", 0))
	userModel.Password = ctx.FormValue("password")
	userModel.Username = ctx.FormValue("username")

	err := userModel.Update()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.OperateSuccess, nil)
}

//新增GM用户
func (this UserController) UserAddAction(ctx *Ctx) error {
	var userModel model.GmUserModel

	userModel.Realname = ctx.FormValue("realname")
	userModel.RoleId = int32(this.FormIntDefault(ctx, "role_id", 0))
	userModel.Password = ctx.FormValue("password")
	userModel.Username = ctx.FormValue("username")

	id, err := userModel.Insert()
	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.OperateSuccess, id)
}
