package middleware

import (
	"strconv"

	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/platform-basic-libs/api_config"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	"github.com/1340691923/ElasticView/platform-basic-libs/rbac"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	fiber "github.com/gofiber/fiber/v2"
)

const ADMIN_ROLE = 1

func Rbac(ctx *fiber.Ctx) error {
	var err error
	token := util.GetToken(ctx)
	var claims *jwt.Claims
	claims, err = jwt.ParseToken(token)
	if err != nil {
		logs.Logger.Sugar().Errorf("OperaterLog err:%s", err.Error())
		return err
	}
	obj := ctx.Path()

	sub := int(claims.RoleId)
	apiRouterConfig := api_config.NewApiRouterConfig()
	//最高权限用户可免接口鉴权
	if sub == ADMIN_ROLE {
		return ctx.Next()
	}
	for _, routerConfig := range apiRouterConfig.GetRouterConfigs() {
		if obj == routerConfig.Url {
			ok, err := rbac.Enforcer.EnforceSafe(strconv.Itoa(sub), obj, "*")
			if err != nil {
				res.Error(ctx, my_error.NewBusiness(TOKEN_ERROR, ERROR_RBAC_LOAD))
				return err
			}
			if !ok {
				res.Error(ctx, my_error.NewBusiness(TOKEN_ERROR, ERROR_RBAC_AUTH))
				return err
			}
		}
	}
	return ctx.Next()
}
