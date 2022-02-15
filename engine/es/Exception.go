package es

import "errors"

//默认异常
var ReqParmasValid = errors.New("请求参数无效")

const (
	IndexNameNullError = 200001
	AliasNameNullError = 200002
)

var ParmasNullError = map[int]string{
	IndexNameNullError: "索引名不能为空",
	AliasNameNullError: "别名不能为空",
}
