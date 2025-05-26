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

func (this *GmUserService) CheckOAuthLogin(provider, providerUserId, username, realname string) (token string, err error) {
	gmUser, err := this.GetUserByOAuthProvider(provider, providerUserId)
	
	if err != nil {
		gmUser = model.GmUserModel{
			Username: username,
			Realname: realname,
			RoleId:   2, // Default role ID
		}
		
		switch provider {
		case "wechat":
			gmUser.WechatOpenId = providerUserId
		case "dingtalk":
			gmUser.DingtalkId = providerUserId
		case "feishu":
			gmUser.FeishuOpenId = providerUserId
		}
		
		gmUser.SetSqlx(this.sqlx)
		gmUser.SetLog(this.log)
		err = gmUser.Insert()
		if err != nil {
			this.log.Error("Failed to create user from OAuth", zap.Error(err), zap.String("provider", provider))
			return "", errors.New("用户创建失败")
		}
	} else {
		sqlstore.SqlBuilder.
			Update("gm_user").
			SetMap(map[string]interface{}{"last_login_time": time.Now().Format(util.TimeFormat)}).
			Where(sqlstore.Eq{"id": gmUser.ID}).
			RunWith(this.sqlx).
			Exec()
	}
	
	token, err = this.jwt.CreateToken(gmUser)
	if err != nil {
		return "", err
	}
	
	return token, nil
}

func (this *GmUserService) GetUserByOAuthProvider(provider, providerUserId string) (gmUser model.GmUserModel, err error) {
	gmUser = model.NewGmUserModel(this.sqlx, this.log)
	
	var query sqlstore.Sqlizer
	
	switch provider {
	case "wechat":
		query = sqlstore.Eq{"wechat_open_id": providerUserId}
	case "dingtalk":
		query = sqlstore.Eq{"dingtalk_id": providerUserId}
	case "feishu":
		query = sqlstore.Eq{"feishu_open_id": providerUserId}
	default:
		return gmUser, errors.New("不支持的认证提供商")
	}
	
	err = sqlstore.SqlBuilder.
		Select("*").
		From("gm_user").
		Where(query).
		RunWith(this.sqlx).
		QueryRow().
		StructScan(&gmUser)
	
	if err != nil {
		this.log.Error("Failed to find user by OAuth provider", 
			zap.Error(err), 
			zap.String("provider", provider),
			zap.String("providerUserId", providerUserId))
		return gmUser, errors.New("用户不存在")
	}
	
	return gmUser, nil
}
