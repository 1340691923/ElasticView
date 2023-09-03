package api

import (
	"errors"
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/1340691923/ElasticView/pkg/services/gm_role"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GM角色控制器
type ManagerRoleController struct {
	*BaseController
	log  *logger.AppLogger
	cfg  *config.Config
	sqlx *sqlstore.SqlStore

	jwtSvr        *jwt_svr.Jwt
	gmRoleService *gm_role.GmRoleService
	rbac          *access_control.Rbac
}

func NewManagerRoleController(baseController *BaseController, rbac *access_control.Rbac, log *logger.AppLogger, cfg *config.Config, sqlx *sqlstore.SqlStore, jwtSvr *jwt_svr.Jwt, gmRoleService *gm_role.GmRoleService) *ManagerRoleController {
	return &ManagerRoleController{BaseController: baseController, rbac: rbac, log: log, cfg: cfg, sqlx: sqlx, jwtSvr: jwtSvr, gmRoleService: gmRoleService}
}

// 获取所有的GM 角色
func (this *ManagerRoleController) RolesAction(ctx *gin.Context) {

	roles, err := this.gmRoleService.Select()
	if err != nil {
		this.Error(ctx, err)
		return
	}
	list, err := this.gmRoleService.GetRoles(roles)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	this.Success(ctx, response.SearchSuccess, list)
	return
}

// 新增GM角色
func (this *ManagerRoleController) RolesAddAction(ctx *gin.Context) {

	var model2 dto.GmRoleModel

	err := ctx.Bind(&model2)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err := roleModel.Insert(this.sqlx)

	if err != nil {
		this.Error(ctx, err)
		return
	}

	go func() {
		for _, api := range model2.Api {
			_, err = this.rbac.AddPolicy(strconv.Itoa(int(id)), api, "*")
			if err != nil {
				this.log.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()

	this.Success(ctx, response.OperateSuccess, map[string]interface{}{"id": id})
}

// 修改GM角色
func (this *ManagerRoleController) RolesUpdateAction(ctx *gin.Context) {
	var model2 dto.GmRoleModel
	err := ctx.Bind(&model2)
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

	if model2.ID == 1 && claims.RoleId != 1 {
		this.Error(ctx, errors.New("您无权修改该角色!"))
		return
	}

	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	err = roleModel.Update(this.sqlx)

	this.rbac.RemoveFilteredPolicy(0, strconv.Itoa(model2.ID)) //先全清掉

	go func() {
		for _, api := range model2.Api {
			_, err = this.rbac.AddPolicy(strconv.Itoa(model2.ID), api, "*")
			if err != nil {
				this.log.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()
	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 删除GM角色
func (this *ManagerRoleController) RolesDelAction(ctx *gin.Context) {

	var reqData dto.RolesDelReq

	err := ctx.Bind(&reqData)
	if err != nil {
		this.Error(ctx, err)
		return
	}

	id := reqData.Id

	claims, err := this.jwtSvr.ParseToken(this.GetToken(ctx))
	if err != nil {
		this.Error(ctx, err)
		return
	}

	if id == 1 && claims.RoleId != 1 {
		this.Error(ctx, errors.New("您无权修改该角色!"))
		return
	}

	err = this.gmRoleService.Delete(id)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.rbac.RemoveFilteredPolicy(0, strconv.Itoa(id)) //先全清掉

	this.Success(ctx, response.OperateSuccess, nil)
	return
}

// 获取Gm角色下拉选
func (this *ManagerRoleController) RoleOptionAction(ctx *gin.Context) {

	var model model.GmRoleModel

	list, err := model.Select(this.sqlx)
	if err != nil {
		this.Error(ctx, err)
		return
	}
	this.Success(ctx, response.SearchSuccess, list)
	return
}
