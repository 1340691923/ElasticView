// GM角色层
package gm_role

import (
	"github.com/1340691923/ElasticView/pkg/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/util"
)

// GmRoleService
type GmRoleService struct {
	log  *logger.AppLogger
	sqlx *sqlstore.SqlStore
}

func NewGmRoleService(log *logger.AppLogger, sqlx *sqlstore.SqlStore) *GmRoleService {
	return &GmRoleService{log: log, sqlx: sqlx}
}

func (this *GmRoleService) Select() (list []model.GmRoleModel, err error) {
	var roleModel model.GmRoleModel
	list, err = roleModel.Select(this.sqlx)
	if err != nil {
		return
	}
	return
}

func (this *GmRoleService) Add(model2 model.GmRoleModel) (id int64, err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err = roleModel.Insert(this.sqlx)
	return
}

func (this *GmRoleService) Update(model2 model.GmRoleModel) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	err = roleModel.Update(this.sqlx)
	return
}

func (this *GmRoleService) Delete(id int) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = id
	err = roleModel.Delete(this.sqlx)
	return
}

func (this *GmRoleService) GetRoles(roles []model.GmRoleModel) (list []dto.GmRoleModel, err error) {
	for _, v := range roles {
		roleRes := dto.GmRoleModel{
			ID:          v.ID,
			RoleName:    v.RoleName,
			Description: v.Description,
			RoleList:    v.RoleList,
		}
		apis := []string{}

		rows, err := this.sqlx.Query("select v1 from casbin_rule where v0 = ?;", v.ID)
		if util.FilterMysqlNilErr(err) {
			this.log.Sugar().Errorf("err:", err)
			continue
		}
		defer rows.Close()
		for rows.Next() {
			api := ""
			err := rows.Scan(&api)
			if err != nil {
				this.log.Sugar().Errorf("err:", err)
				continue
			}
			apis = append(apis, api)
		}
		roleRes.Api = apis
		list = append(list, roleRes)
	}
	return
}
