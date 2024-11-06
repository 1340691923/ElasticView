package api

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/response"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/services/es_link_service"
	"github.com/1340691923/ElasticView/pkg/services/gm_role"
	"github.com/1340691923/ElasticView/pkg/services/gm_user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"strconv"
)

// GM角色控制器
type ManagerRoleController struct {
	*BaseController
	log *logger.AppLogger
	cfg *config.Config

	jwtSvr        *jwt_svr.Jwt
	gmRoleService *gm_role.GmRoleService
	gmUserService *gm_user.GmUserService
	rbac          *access_control.Rbac
	orm           *sqlstore.SqlStore
	esLinkService *es_link_service.EsLinkService
}

func NewManagerRoleController(baseController *BaseController, log *logger.AppLogger, cfg *config.Config, jwtSvr *jwt_svr.Jwt, gmRoleService *gm_role.GmRoleService, gmUserService *gm_user.GmUserService, rbac *access_control.Rbac, orm *sqlstore.SqlStore, esLinkService *es_link_service.EsLinkService) *ManagerRoleController {
	return &ManagerRoleController{BaseController: baseController, log: log, cfg: cfg, jwtSvr: jwtSvr, gmRoleService: gmRoleService, gmUserService: gmUserService, rbac: rbac, orm: orm, esLinkService: esLinkService}
}

// 获取所有的EV角色
func (this *ManagerRoleController) RolesAction(ctx *gin.Context) {

	roles, err := this.gmRoleService.Select(ctx, this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)))
	if err != nil {
		this.Error(ctx, err)
		return
	}
	list, err := this.gmRoleService.GetRoles(ctx, roles)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, list)
	return
}

// 新增EV角色
func (this *ManagerRoleController) RolesAddAction(ctx *gin.Context) {

	var model2 dto.GmRoleModel

	err := ctx.Bind(&model2)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	var roleModel model.GmRole
	roleModel.Id = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = &model2.RoleList
	roleModel.Description = model2.Description
	id, err := this.gmRoleService.Add(ctx, roleModel)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	eg := errgroup.Group{}

	for _, api := range model2.Api {
		api := api
		eg.Go(func() error {
			_, err = this.rbac.AddPolicy(strconv.Itoa(int(id)), api, "*")
			if err != nil {
				this.log.Sugar().Errorf("err:%s", err.Error())
				return err
			}
			return nil
		})
	}
	err = eg.Wait()

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, map[string]interface{}{"id": id})
}

// 修改EV角色
func (this *ManagerRoleController) RolesUpdateAction(ctx *gin.Context) {
	var model2 dto.GmRoleModel
	err := ctx.Bind(&model2)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if model2.ID == 1 && !this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)) {
		this.Error(ctx, errors.New("您无权修改该角色!"))
		return
	}

	var roleModel model.GmRole
	roleModel.Id = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = &model2.RoleList
	roleModel.Description = model2.Description
	err = this.gmRoleService.Update(ctx, roleModel)

	this.rbac.RemoveFilteredPolicy(0, strconv.Itoa(model2.ID)) //先全清掉
	eg := errgroup.Group{}

	for _, api := range model2.Api {
		api := api
		roleId := model2.ID
		eg.Go(func() error {
			_, err = this.rbac.AddPolicy(strconv.Itoa(roleId), api, "*")
			if err != nil {
				this.log.Sugar().Errorf("err:%s", err.Error())
				return err
			}
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 删除EV角色
func (this *ManagerRoleController) RolesDelAction(ctx *gin.Context) {

	var reqData dto.RolesDelReq

	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	id := reqData.Id

	if id == 1 && !this.gmUserService.IsAdminUser(this.GetRoleCache(ctx)) {
		this.Error(ctx, errors.New("您无权修改该角色!"))
		return
	}
	tx := this.orm.Begin()
	err = this.gmRoleService.Delete(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	err = this.gmUserService.DeleteByRoleId(ctx, tx, id)
	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	err = this.esLinkService.DeleteRoleEslinkCfgByRoleId(ctx, tx, id)

	if err != nil {
		tx.Rollback()
		this.Error(ctx, err)
		return
	}

	tx.Commit()
	this.rbac.RemoveFilteredPolicy(0, strconv.Itoa(id)) //先全清掉

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 获取EV角色下拉选
func (this *ManagerRoleController) RoleOptionAction(ctx *gin.Context) {

	list, err := this.gmRoleService.GetOptions(ctx)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
	return
}
