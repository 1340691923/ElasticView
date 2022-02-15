//控制器层
package controller

import (
	"github.com/1340691923/ElasticView/platform-basic-libs/request"
	"github.com/1340691923/ElasticView/platform-basic-libs/response"
)

//父控制器结构体
type BaseController struct {
	response.Response
	request.Request
}
