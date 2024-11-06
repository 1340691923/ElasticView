// 自定义响应 辅助方法层
package response

import (
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	. "github.com/1340691923/ElasticView/pkg/infrastructure/my_error"
	proto2 "github.com/1340691923/eve-plugin-sdk-go/ev_api/proto"
	"github.com/1340691923/eve-plugin-sdk-go/genproto/pluginv2"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

type ResponseData struct {
	Code int         `json:"code"` //消息码
	Msg  string      `json:"msg"`  //消息提示
	Data interface{} `json:"data"` //附加信息
}

// 自定义响应方法
type Response struct {
	log *logger.AppLogger
}

func NewResponse(log *logger.AppLogger) *Response {
	return &Response{log: log.Named("response")}
}

const (
	SUCCESS = 0
	ERROR   = 500
)

const (
	SearchSuccess       = "查询成功"
	DeleteSuccess       = "删除成功"
	OperateSuccess      = "操作成功"
	LogoutSuccess       = "注销成功"
	ChangeLayoutSuccess = "修改布局成功"
)

// 正确信息
func (this *Response) Success(ctx *gin.Context, msg string, data interface{}) error {
	responseData := new(ResponseData)
	responseData.Msg = msg
	responseData.Data = data
	responseData.send(ctx, SUCCESS)
	return nil
}

func (this *Response) SuccessProtobuf(ctx *gin.Context, msg string, data *proto2.Response) error {
	responseData := new(pluginv2.CallResourceResponse)
	responseData.Code = int32(data.StatusCode())
	responseData.Headers = map[string]*pluginv2.StringList{
		"EV-MSG": {Values: []string{msg}},
	}
	responseData.Body = data.ResByte()

	b, err := proto.Marshal(responseData)

	if err != nil {
		return err
	}

	ctx.Writer.Write(b)

	return nil
}

func (this *Response) SuccessProtobufByAny(ctx *gin.Context, msg string, data interface{}) error {
	responseData := new(pluginv2.CallResourceResponse)
	responseData.Code = 200
	responseData.Headers = map[string]*pluginv2.StringList{
		"EV-MSG": {Values: []string{msg}},
	}

	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	responseData.Body = dataBytes

	b, err := proto.Marshal(responseData)

	if err != nil {
		return err
	}

	ctx.Writer.Write(b)

	return nil
}

func (this *Response) ErrorProtobuf(ctx *gin.Context, err error) error {

	req, _ := ctx.GetRawData()

	this.log.Sugar().Errorf("\n请求接口地址:%s\n请求Body:%s\n异常堆栈\n:%+v", ctx.Request.URL.Path, string(req), err)

	responseData := new(pluginv2.CallResourceResponse)
	responseData.Code = 202
	responseData.Headers = map[string]*pluginv2.StringList{
		"EV-MSG": {Values: []string{err.Error()}},
	}
	b, err := proto.Marshal(responseData)

	if err != nil {
		return err
	}

	ctx.Writer.Write(b)

	return nil
}

// 错误信息
func (this *Response) FastError(write io.Writer, err error) error {
	myErr := ErrorToErrorCode(err)

	this.Output(write, map[string]interface{}{
		"code": myErr.Code(),
		"msg":  myErr.Error(),
	})
	return nil
}

// 错误信息
func (this *Response) Error(ctx *gin.Context, err error) error {

	myErr := ErrorToErrorCode(err)

	var b []byte
	b, _ = ctx.GetRawData()

	this.log.Sugar().Errorf("\n请求接口地址:%s\n请求Body:%s\n异常堆栈\n:%+v", ctx.Request.URL.Path, string(b), err)
	responseData := new(ResponseData)
	responseData.Msg = myErr.Error()
	responseData.send(ctx, myErr.Code())
	return nil
}

// 输出
func (this *ResponseData) send(ctx *gin.Context, code int) error {
	this.Code = code

	b, err := json.Marshal(this)

	if err != nil {
		ctx.JSON(http.StatusAccepted, map[string]interface{}{"msg": err, "code": 500})
		return nil
	}
	if this.Code != 0 {
		ctx.Status(http.StatusAccepted)
	} else {
		ctx.Status(http.StatusOK)
	}
	ctx.Writer.Write(b)
	return nil
}

// 输出
func (this *Response) Output(write io.Writer, data map[string]interface{}) error {
	b, _ := json.Marshal(data)
	write.Write(b)
	return nil
}

// 处理异常（业务异常和默认异常）
func ErrorToErrorCode(err error) *MyError {
	if err == nil {
		return nil
	}

	errorCode, ok := err.(*MyError)

	if ok {
		return errorCode
	}
	return NewError(err.Error(), ERROR).(*MyError)
}
