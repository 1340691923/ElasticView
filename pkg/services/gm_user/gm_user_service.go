// BI用户层
package gm_user

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/infrastructure/vo"
	"github.com/1340691923/ElasticView/pkg/services/oauth"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// GmUserService
type GmUserService struct {
	log                  *logger.AppLogger
	gmUserDao            *dao.GmUserDao
	gmRoleDao            *dao.GmRoleDao
	jwt                  *jwt_svr.Jwt
	pluginRegistry       manager.Service
	cfg                  *config.Config
	orm                  *sqlstore.SqlStore
	oAuthServiceRegistry *oauth.OAuthServiceRegistry
}

func NewGmUserService(log *logger.AppLogger, gmUserDao *dao.GmUserDao, gmRoleDao *dao.GmRoleDao, jwt *jwt_svr.Jwt, pluginRegistry manager.Service, cfg *config.Config, orm *sqlstore.SqlStore, oAuthServiceRegistry *oauth.OAuthServiceRegistry) *GmUserService {
	return &GmUserService{log: log, gmUserDao: gmUserDao, gmRoleDao: gmRoleDao, jwt: jwt, pluginRegistry: pluginRegistry, cfg: cfg, orm: orm, oAuthServiceRegistry: oAuthServiceRegistry}
}

func (this *GmUserService) CheckLogin(ctx context.Context, username, password string) (token string, err error) {

	gmUser, err := this.gmUserDao.GetUserByUP(ctx, model.GmUserModel{
		Username: username,
		Password: password,
	})

	if err != nil {
		this.log.Error("登录失败", zap.Error(err))
		err = errors.New("用户验证失败")
		return
	}

	if gmUser.Id == 0 {
		err = errors.New("用户验证失败")
		return
	}

	if gmUser.IsBan == 1 {
		err = errors.New("用户已被封禁")
		return
	}

	token, err = this.jwt.CreateToken(gmUser)
	if err != nil {
		return
	}
	return
}

func (this *GmUserService) CheckLoginByOAuth(ctx context.Context, code string, state string) (token string, err error) {

	application := gjson.Get(state, "application").String()

	svr, err := this.oAuthServiceRegistry.FindServiceByName(application)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if !svr.Enable() {
		return "", fmt.Errorf("已经关闭%s认证", svr.GetAppliactionName())
	}

	oAuthToken, err := svr.GetToken(code)
	if err != nil {
		return "", errors.WithStack(err)
	}
	ui, err := svr.GetUserInfo(oAuthToken)

	if err != nil {
		return "", errors.WithStack(err)
	}
	this.log.Info("第三方登录接口返回", zap.Reflect("用户数据", ui))
	var gmUser model.GmUserModel
	gmUser, err = this.gmUserDao.GetByField(svr.GetUserField(), ui.Id)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if gmUser.Id == 0 {
		gmUser, err = this.gmUserDao.GetByField("email", ui.Email)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	if gmUser.Id == 0 {
		gmUser = model.GmUserModel{
			Username:      ui.Username,
			Password:      util.GetUUid(),
			Avatar:        ui.AvatarUrl,
			Realname:      ui.DisplayName,
			Email:         ui.Email,
			WorkWechatUid: ui.Id,
			UpdateTime:    time.Now(),
			CreateTime:    time.Now(),
			LastLoginTime: time.Now(),
		}
		tx := this.orm.Begin()
		//开始插入数据
		userID, err := this.gmUserDao.Insert(ctx, tx, gmUser)
		if err != nil {
			tx.Rollback()
			return "", errors.WithStack(err)
		}
		gmUser.Id = int(userID)

		err = this.gmUserDao.AddRolesToUser(tx, gmUser.Id, []int{2})
		if err != nil {
			tx.Rollback()
			return "", errors.WithStack(err)
		}
		tx.Commit()

	}

	if gmUser.IsBan == 1 {
		return "", errors.New("用户已被封禁")
	}

	token, err = this.jwt.CreateToken(gmUser)
	if err != nil {
		return
	}

	return token, nil
}

type QiankunMicroApp struct {
	Name       string `json:"name"`
	Entry      string `json:"entry"`
	ActiveRule string `json:"activeRule"`
}

func (this *GmUserService) GetRoleInfo(ctx context.Context, roleIds []int, baseRoutes []*Route) (
	roleMenu string, roleQiankunMicroApps []*QiankunMicroApp, err error) {
	menus := [][]*Route{}
	microApps := map[string]struct{}{}

	for _, roleId := range roleIds {
		//todo...
		gminfo, err := this.gmRoleDao.GetById(ctx, roleId)
		if err != nil {
			return "", nil, errors.WithStack(err)
		}
		if gminfo.RoleList != nil {
			outputMenus, qiankunMicroApps := this.GetRoleList(ctx, *gminfo.RoleList)
			menus = append(menus, outputMenus)

			for _, v := range qiankunMicroApps {
				if _, ok := microApps[v.Name]; !ok {
					roleQiankunMicroApps = append(roleQiankunMicroApps, v)
					microApps[v.Name] = struct{}{}
				}
			}
		}
	}

	mergeMenuTree := this.mergeMenuTrees(menus...)

	mergeMenuTree = this.findIntersection(baseRoutes, mergeMenuTree)

	menusBytes, _ := json.Marshal(mergeMenuTree)

	roleMenu = string(menusBytes)

	return
}

type Route struct {
	Path       string     `json:"path"`
	Name       *string    `json:"name,omitempty"`
	Component  *string    `json:"component,omitempty"`
	Redirect   *string    `json:"redirect,omitempty"`
	AlwaysShow *bool      `json:"alwaysShow,omitempty"`
	Meta       *RouteMeta `json:"meta,omitempty"`
	Children   Children   `json:"children"`
	Hidden     *bool      `json:"hidden,omitempty"`
}

type RouteMeta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type Children []*Route

func (i Children) MarshalJSON() ([]byte, error) {
	if len(i) == 0 {
		return json.Marshal([]interface{}{})
	} else {
		var tempValue []*Route // 这里需要重新定义一个变量，再序列化，否则会造成循环调用
		for _, item := range i {
			tempValue = append(tempValue, item)
		}
		return json.Marshal(tempValue)
	}
}

func addRoutes(destRoutes *Route, srcRoutes []*build.Route) {
	for _, srcRoute := range srcRoutes {
		destRoutes.Children = append(destRoutes.Children, &Route{
			Path: srcRoute.Path,
			Name: util.StringPtr(srcRoute.Path),
			Meta: &RouteMeta{
				Title: srcRoute.Meta.Title,
				Icon:  srcRoute.Meta.Icon,
			},
		})
		if len(srcRoute.Children) > 0 {
			addRoutes(destRoutes, srcRoute.Children)
		}
	}
}

// 合并菜单路由
func (this *GmUserService) mergeMenuTrees(menus ...[]*Route) []*Route {
	menuMap := make(map[string]*Route)

	for _, menu := range menus {
		for _, item := range menu {
			this.mergeRoute(item, menuMap)
		}
	}

	merged := make([]*Route, 0, len(menuMap))
	for _, item := range menuMap {
		merged = append(merged, item)
	}

	return merged
}

func (this *GmUserService) mergeRoute(item *Route, menuMap map[string]*Route) {
	if existing, found := menuMap[cast.ToString(item.Path)]; found {
		// 合并子菜单
		for _, child := range item.Children {
			this.mergeChildRoute(child, existing)
		}
	} else {
		// 创建新的菜单项
		newItem := &Route{
			Path:      item.Path,
			Name:      item.Name,
			Component: item.Component,
			Meta:      item.Meta,
			Children:  []*Route{},
		}
		menuMap[cast.ToString(item.Path)] = newItem

		// 添加子菜单
		for _, child := range item.Children {
			this.mergeChildRoute(child, newItem)
		}
	}
}

func (this *GmUserService) mergeChildRoute(child *Route, parent *Route) {
	if existingChild, found := parent.ChildrenMap()[cast.ToString(child.Path)]; found {
		// 如果子菜单已存在，继续合并其子项
		for _, subChild := range child.Children {
			this.mergeChildRoute(subChild, existingChild)
		}
	} else {
		// 如果子菜单不存在，则直接添加
		parent.Children = append(parent.Children, child)
	}
}

func (m *Route) ChildrenMap() map[string]*Route {
	childMap := make(map[string]*Route)
	for i := range m.Children {
		childMap[cast.ToString(m.Children[i].Path)] = m.Children[i]
	}
	return childMap
}

var EvDefaultMenu = []string{"/connect-tree", "/permission", "/plugins"}

func (this *GmUserService) GetRoleList(ctx context.Context, js string) (outputRoute []*Route, qiankunMicroApps []*QiankunMicroApp) {
	var routes []*Route
	json.Unmarshal([]byte(js), &routes)
	pluginIds := []string{}
	for _, plugin := range this.pluginRegistry.Plugins(ctx) {
		pluginIds = append(pluginIds, fmt.Sprintf("/%s", plugin.PluginID()))
		pluginJson := plugin.PluginData().PluginJsonData
		frontendRoutes := pluginJson.FrontendRoutes

		if len(frontendRoutes) == 0 {
			continue
		}
		entry := fmt.Sprintf("%sapi/call_plugin_views/%s/", this.cfg.RootUrl, plugin.ID)

		if pluginJson.FrontendDebug {
			entry = fmt.Sprintf("http://localhost:%d/", pluginJson.FrontendDevPort)
		}

		_, subUrl, err := this.cfg.ParseAppUrlAndSubUrl()
		if err != nil {
			subUrl = ""
		} else {
			subUrl = subUrl + "/"
		}

		qiankunMicroApps = append(qiankunMicroApps, &QiankunMicroApp{
			Name:       plugin.ID,
			Entry:      entry,
			ActiveRule: fmt.Sprintf("%s#/%s", subUrl, plugin.ID),
		})
		pluginRoutes := &Route{
			Path:      fmt.Sprintf("/%s", plugin.ID),
			Component: util.StringPtr("layout"),
			Name:      util.StringPtr(plugin.ID),
			Meta: &RouteMeta{
				Title: pluginJson.PluginName,
				Icon:  "table",
			},
		}

		addRoutes(pluginRoutes, frontendRoutes)
	}

	for _, route := range routes {
		//ev 默认菜单
		if util.InstrArr(EvDefaultMenu, route.Path) {
			outputRoute = append(outputRoute, route)
		}
		// 插件主菜单
		if util.InstrArr(pluginIds, route.Path) {
			outputRoute = append(outputRoute, route)
		}
	}

	return
}

func (this *GmUserService) MergePluginRoutes(ctx context.Context, routes []*Route) []*Route {
	for _, plugin := range this.pluginRegistry.Plugins(ctx) {

		pluginJson := plugin.PluginData().PluginJsonData
		frontendRoutes := pluginJson.FrontendRoutes

		if len(frontendRoutes) == 0 {
			continue
		}

		pluginRoutes := &Route{
			Path:      fmt.Sprintf("/%s", plugin.ID),
			Component: util.StringPtr("layout"),
			Name:      util.StringPtr(plugin.ID),
			Meta: &RouteMeta{
				Title: pluginJson.PluginName,
				Icon:  "table",
			},
		}

		addRoutes(pluginRoutes, frontendRoutes)

		routes = append(routes, pluginRoutes)
	}
	return routes
}

func (this *GmUserService) findIntersection(routesA, routesB []*Route) []*Route {
	routeMap := make(map[string]*Route)

	// 将第二棵树的路径存储到映射中
	for _, routeB := range routesB {
		routeMap[routeB.Path] = routeB
	}

	var intersection []*Route

	// 遍历第一棵树，构建交集
	for _, routeA := range routesA {

		path := routeA.Path

		if len(routeA.Children) == 1 {

			childrenPath := routeA.Children[0].Path

			routeA.Children[0].Path = fmt.Sprintf("%s/%s", path, childrenPath)
		}

		if routeB, exists := routeMap[path]; exists {

			newRoute := &Route{
				Path:       routeA.Path,
				Name:       routeA.Name,
				Component:  routeA.Component,
				Redirect:   routeA.Redirect,
				AlwaysShow: routeA.AlwaysShow,
				Meta:       routeA.Meta,
				Hidden:     routeA.Hidden,
			}

			// 递归查找子节点的交集
			newRoute.Children = this.findIntersection(routeA.Children, routeB.Children)
			intersection = append(intersection, newRoute)
		}
	}

	return intersection
}

func (this *GmUserService) DeleteByRoleId(ctx context.Context, orm *gorm.DB, roleId int) (err error) {
	err = this.gmUserDao.RemoveRoles(orm, roleId)
	return
}

func (this *GmUserService) IsExsitUser(ctx context.Context, claims *jwt_svr.Claims) (bool, error) {
	//todo...
	return this.gmUserDao.Exsit(ctx, model.GmUserModel{
		Username: cast.ToString(claims.Username),
		//RoleId:   cast.ToInt(claims.RoleId),
	})
}

func (this *GmUserService) UpdatePassById(ctx context.Context, id int, pwd string) error {
	return this.gmUserDao.UpdatePassById(ctx, model.GmUserModel{
		Id:       id,
		Password: pwd,
	})
}

func (this *GmUserService) Select(ctx context.Context, isAdmin bool) (gmUserModel []model.GmUserModel, err error) {
	return this.gmUserDao.Select(ctx, isAdmin)
}

func (this *GmUserService) SealUser(ctx context.Context, id int, isBan bool) (err error) {

	sealUserType := dao.Ban

	if !isBan {
		sealUserType = dao.UnBan
	}

	if isBan && id == AdminRole {
		return errors.New("不可对超管进行封禁操作")
	}

	err = this.gmUserDao.SealUser(ctx, id, sealUserType)
	if err != nil {
		return errors.WithStack(err)
	}

	return
}

func (this *GmUserService) GetUserById(ctx context.Context, id int) (gmUserModel model.GmUserModel, err error) {
	return this.gmUserDao.GetUserById(ctx, id)
}

func (this *GmUserService) Update(ctx context.Context, gmUser model.GmUserModel, roleIds []int) (err error) {
	err = this.gmUserDao.Update(ctx, gmUser)
	if err != nil {
		return errors.WithStack(err)
	}

	err = this.gmUserDao.AddRolesToUser(this.orm.DB, gmUser.Id, roleIds)
	if err != nil {
		return errors.WithStack(err)
	}
	return
}

// 谁起来了事务，则由谁回滚
func (this *GmUserService) Insert(ctx context.Context, gmUser model.GmUserModel, roleIds []int) (id int64, err error) {

	orm := this.orm.Begin()
	userId, err := this.gmUserDao.Insert(ctx, orm, gmUser)
	if err != nil {
		orm.Rollback()
		return 0, errors.WithStack(err)
	}
	err = this.gmUserDao.AddRolesToUser(orm, int(userId), roleIds)
	if err != nil {
		orm.Rollback()
		return 0, errors.WithStack(err)
	}
	orm.Commit()

	return userId, nil
}

func (this *GmUserService) Delete(ctx context.Context, id int) (err error) {
	tx := this.orm.Begin()
	err = this.gmUserDao.RemoveUserRoles(tx, id)
	if err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	err = this.gmUserDao.Delete(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	tx.Commit()

	return nil
}

func (this *GmUserService) GetOAuthList(callback string) (cfgs []vo.OAuthConfig, err error) {
	svrs := this.oAuthServiceRegistry.GetServices()

	authCallback := this.cfg.RootUrl + "api/callback"

	for _, v := range svrs {

		stateMap := map[string]interface{}{
			"application":    v.GetAppliactionName(),
			"login_callback": callback,
		}

		cfgs = append(cfgs, vo.OAuthConfig{
			OauthUrl: v.GetOAuthUrl(authCallback, stateMap),
			Name:     v.GetAppliactionName(),
			Enable:   v.Enable(),
			Img:      v.GetImg(),
		})
	}
	return
}

func (this *GmUserService) GetOAuthConfigs() (cfgs map[string]map[string]interface{}, err error) {
	svrs := this.oAuthServiceRegistry.GetServices()
	cfgs = map[string]map[string]interface{}{}
	for _, v := range svrs {
		cfgs[v.GetAppliactionName()] = v.GetConfig()
	}
	return
}

func (this *GmUserService) SaveOAuthConfigs(applicationName string, data map[string]interface{}) (err error) {
	svr, err := this.oAuthServiceRegistry.FindServiceByName(applicationName)
	if err != nil {
		return errors.WithStack(err)
	}
	svr.SetConfig(data)
	return
}

const AdminRole = 1

func (this *GmUserService) IsAdminUser(roleId []int) bool {
	return util.InArr(roleId, AdminRole)
}

func (this *GmUserService) GetRolesByUserID(userId int) ([]int, error) {
	return this.gmUserDao.GetRolesFromUser(userId)
}
