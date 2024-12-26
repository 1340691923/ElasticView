// GM角色层
package gm_role

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"gorm.io/gorm"
)

// GmRoleService
type GmRoleService struct {
	log *logger.AppLogger

	roleDao *dao.GmRoleDao
}

func NewGmRoleService(log *logger.AppLogger, roleDao *dao.GmRoleDao) *GmRoleService {
	return &GmRoleService{log: log, roleDao: roleDao}
}

func (this *GmRoleService) Select(
	ctx context.Context, isAdmin bool,
	roleName string,
	page int, pageSize int,
) (list []model.GmRole, count int64, err error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	list, count, err = this.roleDao.SelectByPage(ctx, isAdmin, roleName, page, pageSize)
	if err != nil {
		return
	}
	return
}

func (this *GmRoleService) GetOptions(ctx context.Context) (list []model.GmRole, err error) {
	list, err = this.roleDao.Select(ctx, true)
	if err != nil {
		return
	}
	return
}

func (this *GmRoleService) Add(ctx context.Context, model2 model.GmRole) (id int64, err error) {

	id, err = this.roleDao.Insert(ctx, &model.GmRole{
		RoleName:    model2.RoleName,
		Description: model2.Description,
		RoleList:    model2.RoleList,
	})
	return
}

func (this *GmRoleService) Update(ctx context.Context, model2 model.GmRole) (err error) {
	err = this.roleDao.Update(ctx, model2)
	return
}

func (this *GmRoleService) Delete(ctx context.Context, orm *gorm.DB, id int) (err error) {
	err = this.roleDao.Delete(ctx, orm, id)
	return
}

func (this *GmRoleService) GetRoles(ctx context.Context, roles []model.GmRole) (list []dto.GmRoleModel, err error) {
	for _, v := range roles {
		roleRes := dto.GmRoleModel{
			ID:          int(v.Id),
			RoleName:    v.RoleName,
			Description: v.Description,
			RoleList:    *v.RoleList,
		}

		roleRes.Api = this.roleDao.GetApis(ctx, int(v.Id))
		list = append(list, roleRes)
	}
	return
}
