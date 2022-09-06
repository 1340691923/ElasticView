//自定义请求辅助方法层
package request

//自定义业务异常
const (
	IdNullError = 100002
)

var ParmasNullError = map[int]string{
	IdNullError: "id不能为空！",
}
