// 中间件层
package middleware

import (
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/jwt"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	fiber "github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
)

func OperaterLog(ctx *fiber.Ctx) error {
	var err error
	token := util.GetToken(ctx)
	var claims *jwt.Claims
	claims, err = jwt.ParseToken(token)
	if err != nil {
		logs.Logger.Sugar().Errorf("OperaterLog err:%s", err.Error())
		return err
	}

	parmasMap := util.Map{}
	bodyMap := util.Map{}

	err = ctx.BodyParser(bodyMap)
	if err != nil {
		logs.Logger.Sugar().Errorf("ctx.BodyParser err:%s", err.Error())
	}
	err = ctx.QueryParser(parmasMap)
	if err != nil {
		logs.Logger.Sugar().Errorf("ctx.QueryParser err:%s", err.Error())
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	parmas, _ := json.MarshalToString(parmasMap)
	body, _ := json.MarshalToString(bodyMap)

	gmOperaterLog := model.GmOperaterLog{
		OperaterName:   claims.Username,
		OperaterId:     int64(claims.ID),
		OperaterAction: ctx.Path(),
		Method:         ctx.Method(),
		Parmas:         parmas,
		Body:           body,
		OperaterRoleId: int(claims.RoleId),
	}

	err = gmOperaterLog.Insert()

	if err != nil {
		logs.Logger.Sugar().Errorf("OperaterLog err:%s", err.Error())
	}

	return ctx.Next()

}
