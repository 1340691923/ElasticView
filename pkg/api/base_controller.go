package api

import (
	"bytes"
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/1340691923/ElasticView/pkg/response"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// 父控制器结构体
type BaseController struct {
	*request.Request
	*response.Response
}

func NewBaseController(request *request.Request, response *response.Response) *BaseController {
	return &BaseController{Request: request, Response: response}
}

func (this *BaseController) getPostBody(ctx *gin.Context) []byte {
	body, _ := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}
