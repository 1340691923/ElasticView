package controller

import (
	"strconv"

	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/rbac"
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
	"github.com/1340691923/ElasticView/platform-basic-libs/service/gm_role"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	. "github.com/gofiber/fiber/v2"
)

// GM角色控制器
type RoleController struct {
	BaseController
}

//获取所有的GM 角色
func (this RoleController) RolesAction(ctx *Ctx) error {
	var service gm_role.GmRoleService
	roles, err := service.Select()
	if err != nil {
		return this.Error(ctx, err)
	}
	var list []request.GmRoleModel
	for _, v := range roles {
		roleRes := request.GmRoleModel{
			ID:          v.ID,
			RoleName:    v.RoleName,
			Description: v.Description,
			RoleList:    v.RoleList,
		}
		apis := []string{}

		rows, err := db.Sqlx.Query("select v1 from casbin_rule where v0 = ?;", v.ID)
		if util.FilterMysqlNilErr(err) {
			logs.Logger.Sugar().Errorf("err:", err)
			continue
		}
		defer rows.Close()
		for rows.Next() {
			api := ""
			err := rows.Scan(&api)
			if err != nil {
				logs.Logger.Sugar().Errorf("err:", err)
				continue
			}
			apis = append(apis, api)
		}
		roleRes.Api = apis
		list = append(list, roleRes)
	}

	if err != nil {
		return this.Error(ctx, err)
	}

	return this.Success(ctx, response.SearchSuccess, list)
}

//新增GM角色
func (this RoleController) RolesAddAction(ctx *Ctx) error {

	var model2 request.GmRoleModel

	err := ctx.BodyParser(&model2)
	if err != nil {
		return this.Error(ctx, err)
	}
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err := roleModel.Insert()

	if err != nil {
		return this.Error(ctx, err)
	}

	go func() {
		for _, api := range model2.Api {
			_, err = rbac.Enforcer.AddPolicySafe(strconv.Itoa(int(id)), api, "*")
			if err != nil {
				logs.Logger.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()

	return this.Success(ctx, response.OperateSuccess, map[string]interface{}{"id": id})
}

// 修改GM角色
func (this RoleController) RolesUpdateAction(ctx *Ctx) error {
	var model2 request.GmRoleModel
	err := ctx.BodyParser(&model2)
	if err != nil {
		return this.Error(ctx, err)
	}
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	err = roleModel.Update()

	rbac.Enforcer.RemoveFilteredPolicy(0, strconv.Itoa(model2.ID)) //先全清掉

	go func() {
		for _, api := range model2.Api {
			_, err = rbac.Enforcer.AddPolicySafe(strconv.Itoa(model2.ID), api, "*")
			if err != nil {
				logs.Logger.Sugar().Errorf("err:%s", err.Error())
				return
			}
		}
	}()
	return this.Success(ctx, response.OperateSuccess, nil)
}

// 删除GM角色
func (this RoleController) RolesDelAction(ctx *Ctx) error {
	id := this.FormIntDefault(ctx, "id", 0)

	var service gm_role.GmRoleService
	err := service.Delete(id)
	if err != nil {
		return this.Error(ctx, err)
	}
	rbac.Enforcer.RemoveFilteredPolicy(0, strconv.Itoa(id)) //先全清掉

	return this.Success(ctx, response.OperateSuccess, nil)
}

// 获取Gm角色下拉选
func (this RoleController) RoleOptionAction(ctx *Ctx) error {

	var model model.GmRoleModel

	list, err := model.Select()
	if err != nil {
		return this.Error(ctx, err)
	}
	return this.Success(ctx, response.SearchSuccess, list)
}
