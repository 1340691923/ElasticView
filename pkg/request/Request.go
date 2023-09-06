package request

import (
	"github.com/gin-gonic/gin"
)

// 自定义请求 辅助方法
type Request struct {
}

func NewRequest() *Request {
	return &Request{}
}

// 获取用户token信息
func (this Request) GetToken(ctx *gin.Context) (token string) {
	return ctx.GetHeader("X-Token")
}
