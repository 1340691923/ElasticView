package util

import (
	"github.com/valyala/fasthttp"
	"strings"
)

// 获取真实的IP  1.1.1.1, 2.2.2.2, 3.3.3.3
func CtxClientIP(ctx *fasthttp.RequestCtx) string {
	clientIP := Bytes2str(ctx.Request.Header.Peek("X-Forwarded-For"))
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
		//获取最开始的一个 即 1.1.1.1
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = strings.TrimSpace(Bytes2str(ctx.Request.Header.Peek("X-Real-Ip")))
	if len(clientIP) > 0 {
		return clientIP
	}
	return ctx.RemoteIP().String()
}
