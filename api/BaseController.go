// 控制器层
package api

import (
	"github.com/1340691923/ElasticView/pkg/request"
	"github.com/1340691923/ElasticView/pkg/response"
)

// 父控制器结构体
type BaseController struct {
	response.Response
	request.Request
}
