//GM角色层
package gm_role

import "github.com/1340691923/ElasticView/model"

// GmRoleService
type GmRoleService struct {
}

func (this GmRoleService) Select() (list []model.GmRoleModel, err error) {
	var roleModel model.GmRoleModel
	list, err = roleModel.Select()
	if err != nil {
		return
	}
	return
}

func (this GmRoleService) Add(model2 model.GmRoleModel) (id int64, err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	id, err = roleModel.Insert()
	return
}

func (this GmRoleService) Update(model2 model.GmRoleModel) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = model2.ID
	roleModel.RoleName = model2.RoleName
	roleModel.RoleList = model2.RoleList
	roleModel.Description = model2.Description
	roleModel.ID = model2.ID
	err = roleModel.Update()
	return
}

func (this GmRoleService) Delete(id int) (err error) {
	var roleModel model.GmRoleModel
	roleModel.ID = id
	err = roleModel.Delete()
	return
}
