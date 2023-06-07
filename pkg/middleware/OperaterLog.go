// 中间件层
package middleware

import (
	"encoding/json"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/jwt"
	"github.com/1340691923/ElasticView/pkg/util"
	fiber "github.com/gofiber/fiber/v2"
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

	parmas, _ := json.Marshal(parmasMap)
	body, _ := json.Marshal(bodyMap)

	gmOperaterLog := model.GmOperaterLog{
		OperaterName:   claims.Username,
		OperaterId:     int64(claims.ID),
		OperaterAction: ctx.Path(),
		Method:         ctx.Method(),
		Parmas:         string(parmas),
		Body:           string(body),
		OperaterRoleId: int(claims.RoleId),
	}

	err = gmOperaterLog.Insert()

	if err != nil {
		logs.Logger.Sugar().Errorf("OperaterLog err:%s", err.Error())
	}

	return ctx.Next()

}
