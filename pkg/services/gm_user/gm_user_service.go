// BI用户层
package gm_user

import (
	"errors"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	"github.com/1340691923/ElasticView/pkg/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/model"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

// GmUserService
type GmUserService struct {
	log  *logger.AppLogger
	sqlx *sqlstore.SqlStore
	jwt  *jwt_svr.Jwt
}

func NewGmUserService(log *logger.AppLogger, sqlx *sqlstore.SqlStore, jwt *jwt_svr.Jwt) *GmUserService {
	return &GmUserService{log: log, sqlx: sqlx, jwt: jwt}
}

func (this *GmUserService) CheckLogin(username, password string) (token string, err error) {
	model2 := model.NewGmUserModel(this.sqlx, this.log)
	model2.Password = password
	model2.Username = username

	gmUser, err := model2.GetUserByUP()

	if err != nil {
		this.log.Error("登录失败", zap.Error(err))
		err = errors.New("用户验证失败")
		return
	}

	sqlstore.SqlBuilder.
		Update("gm_user").
		SetMap(map[string]interface{}{"last_login_time": time.Now().Format(util.TimeFormat)}).
		Where(sqlstore.Eq{"id": gmUser.ID}).
		RunWith(this.sqlx).
		Exec()

	token, err = this.jwt.CreateToken(gmUser)
	if err != nil {
		return
	}
	return
}

func (this *GmUserService) GetRoleInfo(roleId int32) (gminfo model.GmRoleModel, err error) {
	var model2 model.GmRoleModel
	gminfo, err = model2.GetById(this.sqlx, int(roleId))
	if err != nil {
		return
	}
	return
}

func (this *GmUserService) IsExitUser(claims *jwt_svr.Claims) bool {
	model2 := model.NewGmUserModel(this.sqlx, this.log)
	model2.Username = cast.ToString(claims.Username)
	model2.RoleId = cast.ToInt32(claims.RoleId)
	return model2.Exsit()
}
