package api

import (
	"errors"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"time"
)

// BI用户控制器
type ManagerUserController struct {
	*BaseController
	log           *logger.AppLogger
	cfg           *config.Config
	sqlx          *sqlstore.SqlStore
	jwtSvr        *jwt_svr.Jwt
	gmUserService *gm_user.GmUserService
	routerEngine  *web_engine.WebEngine
}

func NewManagerUserController(baseController *BaseController, routerEngine *web_engine.WebEngine,
	log *logger.AppLogger, cfg *config.Config,
	sqlx *sqlstore.SqlStore, jwtSvr *jwt_svr.Jwt,
	gmUserService *gm_user.GmUserService) *ManagerUserController {
	return &ManagerUserController{
		BaseController: baseController, log: log,
		cfg: cfg, sqlx: sqlx, jwtSvr: jwtSvr, gmUserService: gmUserService,
		routerEngine: routerEngine}
}

// 登录
func (this *ManagerUserController) Login(ctx *gin.Context) {

	type ReqData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var reqData ReqData

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	username := reqData.Username
	password := reqData.Password

	token, err := this.gmUserService.CheckLogin(username, password)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "登录成功", map[string]interface{}{"token": token, "unix_time": time.Now().Unix()})
}

// 修改自己的密码
func (this *ManagerUserController) ModifyPwd(ctx *gin.Context) {
	type ReqData struct {
		Password string `json:"password"`
	}

	var reqData ReqData

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Password == "" {
		this.Error(ctx, errors.New("密码不能为空"))
		return
	}

	token := this.GetToken(ctx)
	claims, err := this.jwtSvr.ParseToken(token)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	gmUserModel := model.NewGmUserModel(this.sqlx, this.log)
	gmUserModel.ID = cast.ToInt32(claims.UserID)
	gmUserModel.Password = reqData.Password
	err = gmUserModel.UpdatePassById()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	util.TokenBucket.LoadOrStore(token, cast.ToInt64(claims.ExpiresAt.Unix()))

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 用户详细信息
func (this *ManagerUserController) UserInfo(ctx *gin.Context) {
	gmUserService := this.gmUserService
	token := this.GetToken(ctx)
	claims, err := this.jwtSvr.ParseToken(token)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	info, err := gmUserService.GetRoleInfo(claims.RoleId)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, "登录成功", map[string]interface{}{"roles": []string{"admin"}, "introduction": info.Description, "name": claims.RealName, "list": info.RoleList, "avatar": ""})
}

// 退出登录
func (this *ManagerUserController) LogoutAction(ctx *gin.Context) {
	token := this.GetToken(ctx)
	var claims *jwt_svr.Claims
	claims, err := this.jwtSvr.ParseToken(token)
	if err != nil {
		this.log.Error("LogoutAction err", zap.Error(err))
		this.Success(ctx, response.LogoutSuccess, nil)
	}
	util.TokenBucket.LoadOrStore(token, cast.ToInt64(claims.ExpiresAt.Unix()))

	this.Success(ctx, response.LogoutSuccess, nil)
}

// BI用户列表
func (this *ManagerUserController) UserListAction(ctx *gin.Context) {

	appid := gjson.GetBytes(this.getPostBody(ctx), "appid").String()
	notHaveCurrentUser := gjson.GetBytes(this.getPostBody(ctx), "notHaveCurrentUser").Bool()
	userModel := model.NewGmUserModel(this.sqlx, this.log)
	var err error
	var list []model.GmUserModel

	claims, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.log.Error("jwt.ParseToken", zap.Error(err))
		this.Error(ctx, err)
		return
	}

	if notHaveCurrentUser {
		list, err = userModel.SelectByUid(appid, int(claims.UserID))
	} else {
		list, err = userModel.Select(appid)
	}

	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
	return
}

// 删除BI用户
func (this *ManagerUserController) DeleteUserAction(ctx *gin.Context) {

	var reqData dto.DeleteUserReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	userModel := model.NewGmUserModel(this.sqlx, this.log)
	userModel.ID = reqData.Id
	if userModel.ID == 1 {
		this.Error(ctx, errors.New("您无权删除该用户!"))
		return
	}

	err = userModel.Delete()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.DeleteSuccess, nil)
}

// 用ID获取用户信息
func (this *ManagerUserController) GetUserByIdAction(ctx *gin.Context) {
	var reqData dto.GetUserByIdReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	userModel := model.NewGmUserModel(this.sqlx, this.log)
	userModel.ID = reqData.Id
	gmUser, err := userModel.GetUserById()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, gmUser)
}

// 修改BI用户信息
func (this *ManagerUserController) UserUpdateAction(ctx *gin.Context) {

	var reqData dto.UserUpdateReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	userModel := model.NewGmUserModel(this.sqlx, this.log)
	var id = reqData.Id

	userModel.ID = int32(id)
	userModel.Realname = reqData.Realname
	userModel.RoleId = reqData.RoleId
	userModel.Password = reqData.Password
	userModel.Username = reqData.Username
	spew.Dump(reqData)
	err = userModel.Update()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 新增BI用户
func (this *ManagerUserController) UserAddAction(ctx *gin.Context) {

	var reqData dto.UserAddReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	userModel := model.NewGmUserModel(this.sqlx, this.log)

	userModel.Realname = reqData.Realname
	userModel.RoleId = reqData.RoleId
	userModel.Password = reqData.Password
	userModel.Username = reqData.Username
	id, err := userModel.Insert()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, id)
}

func (this *ManagerUserController) UrlConfig(ctx *gin.Context) {
	type RouterConfig struct {
		Url    string `json:"url"`
		Remark string `json:"remark"`
	}
	list := []RouterConfig{}
	for _, rg := range this.routerEngine.GetRouterConfigGroups() {
		for _, router := range rg.RouterConfigs {
			if router.NeedAuth {
				list = append(list, RouterConfig{
					Url:    router.Url,
					Remark: fmt.Sprintf("%s-%s", rg.GroupRemark, router.Remark),
				})
			}
		}
	}
	this.Success(ctx, response.SearchSuccess, list)
}
