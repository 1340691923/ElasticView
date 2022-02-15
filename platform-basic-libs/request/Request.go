package request

import (
	"strconv"
	"strings"

	"github.com/1340691923/ElasticView/platform-basic-libs/my_error"
	fiber "github.com/gofiber/fiber/v2"
)

//自定义请求 辅助方法
type Request struct {
}

type CheckConfigStruct struct {
	Code int
	Key  string
}

//检查请求参数
func (this Request) CheckParameter(checkConfig []CheckConfigStruct, ctx *fiber.Ctx) (err error) {
	method := strings.ToUpper(ctx.Method())
	for _, config := range checkConfig {
		if method == "GET" {
			if ctx.FormValue(config.Key) == "" {
				err = my_error.NewBusiness(ParmasNullError, config.Code)
				return
			}
		} else if method == "POST" {
			if ctx.FormValue(config.Key) == "" {
				err = my_error.NewBusiness(ParmasNullError, config.Code)
				return
			}
		}
	}
	return
}

// FormIntDefault 获取Form参数 如果出错则返回默认值
func (this Request) FormIntDefault(ctx *fiber.Ctx, key string, def int) int {
	i, err := strconv.Atoi(ctx.FormValue(key))
	if err != nil {
		return def
	}
	return i
}

//获取用户token信息
func (this Request) GetToken(ctx *fiber.Ctx) (token string) {
	return ctx.Get("X-Token")
}
