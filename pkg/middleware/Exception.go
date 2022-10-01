package middleware

// 内置异常
const (
	INVALID_PARAMS                 int = 40001
	ERROR_AUTH_CHECK_TOKEN_FAIL    int = 40002
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT int = 40003
	ERROR_RBAC_LOAD                int = 40004
	ERROR_RBAC_AUTH                int = 40005
)

// 内置异常表 TOKEN_ERROR
var TOKEN_ERROR = map[int]string{
	INVALID_PARAMS:                 "Token不能为空",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_RBAC_LOAD:                "读取rdbc权限列表失败",
	ERROR_RBAC_AUTH:                "您没有该资源的访问权限",
}
