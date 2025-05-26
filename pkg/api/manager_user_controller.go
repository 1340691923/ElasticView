package api

import (
	"fmt"
	"time"

	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/web_engine"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// BI用户控制器
type ManagerUserController struct {
	*BaseController
	log            *logger.AppLogger
	cfg            *config.Config
	orm            *orm.Gorm
	jwtSvr         *jwt_svr.Jwt
	gmUserService  *gm_user.GmUserService
	routerEngine   *web_engine.WebEngine
	evUpdate       *updatechecker.EvUpdate
	pluginRegistry manager.Service
}

func NewManagerUserController(baseController *BaseController, log *logger.AppLogger, cfg *config.Config, orm *orm.Gorm, jwtSvr *jwt_svr.Jwt, gmUserService *gm_user.GmUserService, routerEngine *web_engine.WebEngine, evUpdate *updatechecker.EvUpdate, pluginRegistry manager.Service) *ManagerUserController {
	return &ManagerUserController{BaseController: baseController, log: log, cfg: cfg, orm: orm, jwtSvr: jwtSvr, gmUserService: gmUserService, routerEngine: routerEngine, evUpdate: evUpdate, pluginRegistry: pluginRegistry}
}

// 登录
// @Summary EV用户登录
// @Accept application/json
// @Produce application/json
// @Param X-Token header string false "用户令牌"
// @Param object body dto.User false "查询参数"
// @Success 0 {object} vo.User
// @Router /api/gm_user/login [post]
func (this *ManagerUserController) Login(ctx *gin.Context) {
	var reqData dto.User

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	var token string
	if reqData.OAuthCode != "" && reqData.State != "" {
		token, err = this.gmUserService.CheckLoginByOAuth(ctx, reqData.OAuthCode, reqData.State)
		if err != nil {
			this.Error(ctx, err)
			return
		}
	} else {
		username := reqData.Username
		password := reqData.Password

		token, err = this.gmUserService.CheckLogin(ctx, username, password)
		if err != nil {
			this.Error(ctx, err)
			return
		}
	}

	this.Success(ctx, "登录成功", vo.User{
		Token:    token,
		UnixTime: time.Now().Unix(),
	})
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

	err = this.gmUserService.UpdatePassById(ctx, claims.UserID, reqData.Password)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	util.TokenBucket.LoadOrStore(token, cast.ToInt64(claims.ExpiresAt.Unix()))

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 修改用户的密码
func (this *ManagerUserController) ModifyPwdByUserId(ctx *gin.Context) {
	type ReqData struct {
		Id       int    `json:"id"`
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

	err = this.gmUserService.UpdatePassById(ctx, reqData.Id, reqData.Password)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 用户详细信息
func (this *ManagerUserController) UserInfo(ctx *gin.Context) {
	type ReqData struct {
		BaseRoutes []*gm_user.Route `json:"baseRoutes"`
	}
	var reqData ReqData

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	token := this.GetToken(ctx)
	claims, err := this.jwtSvr.ParseToken(token)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	baseRoutes := this.gmUserService.MergePluginRoutes(ctx, reqData.BaseRoutes)

	gmUserService := this.gmUserService

	roleList, qiankunApps, err := gmUserService.GetRoleInfo(ctx, this.GetRoleCache(ctx), baseRoutes)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)) {
		roleListByte, _ := json.Marshal(baseRoutes)
		roleList = string(roleListByte)
	}

	claims.Avatar = "https://oss.youlai.tech/youlai-boot/2023/05/16/811270ef31f548af9cffc026dfc3777b.gif?imageView2/1/w/80/h/80"
	this.Success(ctx, "登录成功", map[string]interface{}{
		"roles":             []string{"admin"},
		"introduction":      "",
		"name":              claims.RealName,
		"list":              roleList,
		"qiankunMicroApps":  qiankunApps,
		"avatar":            claims.Avatar,
		"evUpdateAvailable": this.evUpdate.UpdateAvailable(),
		"evLatestVersion":   this.evUpdate.LatestVersion(),
		"evDownloadUrl":     this.evUpdate.DownloadUrl(),
	})
	return

}

func (this *ManagerUserController) GetRoutesConfig(ctx *gin.Context) {
	type ReqData struct {
		Routers []*gm_user.Route `json:"routers"`
	}
	var reqData ReqData

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	routes := this.gmUserService.MergePluginRoutes(ctx, reqData.Routers)

	this.Success(ctx, response.SearchSuccess, routes)
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

// EV用户列表
func (this *ManagerUserController) UserListAction(ctx *gin.Context) {

	var err error

	var reqData dto.SearchUserReq

	err = ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	var list []model.GmUserModel
	var count int64
	list, count, err = this.gmUserService.Select(
		ctx, this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)),
		reqData.UserName, reqData.RealName, reqData.IsBan,
		reqData.RoleIds, reqData.UserIds,
		reqData.Page, reqData.PageSize)

	if err != nil {
		this.Error(ctx, err)
		return
	}
	var users []vo.GmUsers
	for _, user := range list {
		roles := []int{}
		roles, err = this.gmUserService.GetRolesByUserID(user.Id)
		if err != nil {
			this.log.Error("GetRolesByUserID", zap.Error(err))
		}
		users = append(users, vo.GmUsers{
			Id:            user.Id,
			Username:      user.Username,
			Password:      user.Password,
			Avatar:        user.Avatar,
			Realname:      user.Realname,
			IsBan:         user.IsBan,
			Email:         user.Email,
			WorkWechatUid: user.WorkWechatUid,
			RoleIds:       roles,
			UpdateTime:    user.UpdateTime.Format(time.DateTime),
			CreateTime:    user.CreateTime.Format(time.DateTime),
			LastLoginTime: user.LastLoginTime.Format(time.DateTime),
		})
	}

	this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"list":  users,
		"count": count,
	})
	return
}

// 封禁用户
func (this *ManagerUserController) SealUserAction(ctx *gin.Context) {
	var reqData dto.SealUserReq
	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.gmUserService.SealUser(ctx, reqData.Id, true)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
}

// 解封用户
func (this *ManagerUserController) UnSealUserAction(ctx *gin.Context) {
	var reqData dto.SealUserReq
	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.gmUserService.SealUser(ctx, reqData.Id, false)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)

}

// 删除EV用户
func (this *ManagerUserController) DeleteUserAction(ctx *gin.Context) {

	var reqData dto.DeleteUserReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Id == 1 {
		this.Error(ctx, errors.New("您无权删除该用户!"))
		return
	}

	err = this.gmUserService.Delete(ctx, reqData.Id)
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

	gmUser, err := this.gmUserService.GetUserById(ctx, reqData.Id)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, gmUser)
}

// 修改EV用户信息
func (this *ManagerUserController) UserUpdateAction(ctx *gin.Context) {

	var reqData dto.UserUpdateReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Username == "" {
		this.Error(ctx, errors.New("用户名不能为空"))
	}

	if util.InArr(reqData.RoleIds, gm_user.AdminRole) && !this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)) {
		this.Error(ctx, errors.New("非管理员权限组无法修改管理员信息"))
		return
	}

	var id = reqData.Id

	err = this.gmUserService.Update(ctx, model.GmUserModel{
		Id:            int(id),
		Username:      reqData.Username,
		Email:         reqData.Email,
		Realname:      reqData.Realname,
		WorkWechatUid: reqData.WorkWechatUid,
	}, reqData.RoleIds)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 新增EV用户
func (this *ManagerUserController) UserAddAction(ctx *gin.Context) {

	var reqData dto.UserAddReq

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	if reqData.Username == "" {
		this.Error(ctx, errors.New("用户名不能为空"))
	}

	if reqData.Password == "" {
		this.Error(ctx, errors.New("密码不能为空"))
	}

	if util.InArr(reqData.RoleIds, gm_user.AdminRole) && !this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)) {
		this.Error(ctx, errors.New("非管理员权限组无法新增管理员信息"))
		return
	}

	userModel := model.GmUserModel{}

	userModel.Realname = reqData.Realname
	userModel.Password = reqData.Password
	userModel.Username = reqData.Username
	userModel.Email = reqData.Email
	userModel.WorkWechatUid = reqData.WorkWechatUid

	id, err := this.gmUserService.Insert(ctx, userModel, reqData.RoleIds)
	if err != nil {

		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, id)
}

func (this *ManagerUserController) UrlConfig(ctx *gin.Context) {

	reqData := new(dto.UrlConfigReq)

	err := ctx.Bind(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	list := []vo.RouterConfig{}
	cfg_with_module := this.routerEngine.GetRouterConfigGroups()

	for _, plugin := range this.pluginRegistry.Plugins(ctx) {
		routers := []web_engine.RouterConfig{}
		for _, backendRoute := range plugin.PluginData().PluginJsonData.BackendRoutes {
			routers = append(routers, web_engine.RouterConfig{
				Url:      fmt.Sprintf("/%s%s", plugin.PluginData().PluginJsonData.PluginAlias, backendRoute.Path),
				Remark:   backendRoute.Remark,
				NeedAuth: backendRoute.NeedAuth,
			})
		}
		cfg_with_module = append(cfg_with_module, web_engine.RouterConfigGroup{
			GroupRemark:   plugin.PluginData().PluginJsonData.PluginName,
			RouterConfigs: routers,
		})
	}

	for _, rg := range cfg_with_module {
		for _, router := range rg.RouterConfigs {
			list = append(list, vo.RouterConfig{
				Url:    router.Url,
				Remark: fmt.Sprintf("%s-%s", rg.GroupRemark, router.Remark),
			})
		}

	}

	this.Success(ctx, response.SearchSuccess, map[string]interface{}{
		"cfg":             list,
		"cfg_with_module": cfg_with_module,
	})
}

func (this *ManagerUserController) UserInfoV2(ctx *gin.Context) {
	token := this.GetToken(ctx)
	claims, err := this.jwtSvr.ParseToken(token)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res := vo.UserInfoV2{}
	res.UserId = claims.UserID
	res.Username = claims.RealName
	res.Avatar = claims.Avatar
	res.Perms = make([]string, 0)

	roleIds, err := this.gmUserService.GetRolesByUserID(res.UserId)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	res.Roles = roleIds
	this.Success(ctx, response.SearchSuccess, res)
}

func (this *ManagerUserController) GetOAuthList(ctx *gin.Context) {

	var reqData dto.GetOAuthConfigReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	cfgs, err := this.gmUserService.GetOAuthList(reqData.CallBack)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, cfgs)
}

func (this *ManagerUserController) GetOAuthConfigs(ctx *gin.Context) {

	cfgs, err := this.gmUserService.GetOAuthConfigs()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, cfgs)
}

func (this *ManagerUserController) SaveOAuthConfigs(ctx *gin.Context) {
	var reqData dto.SaveOAuthConfigReq

	err := ctx.BindJSON(&reqData)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	err = this.gmUserService.SaveOAuthConfigs(reqData.ApplicationName, reqData.Config)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.OperateSuccess, nil)
}

func (this *ManagerUserController) NoAuthRoute(ctx *gin.Context) {

	gmUserService := this.gmUserService

	roleList, qiankunApps := gmUserService.GetRoleList2C(ctx)

	avatar := "https://oss.youlai.tech/youlai-boot/2023/05/16/811270ef31f548af9cffc026dfc3777b.gif?imageView2/1/w/80/h/80"
	this.Success(ctx, "获取成功", map[string]interface{}{
		"roles":            []string{},
		"introduction":     "",
		"name":             "游客",
		"list":             roleList,
		"qiankunMicroApps": qiankunApps,
		"avatar":           avatar,
	})

	return

}
