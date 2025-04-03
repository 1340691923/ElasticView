package web

import "github.com/gin-gonic/gin"

// ES 任务 路由
func (this *WebServer) runWs() {
	this.engine.GetGinEngine().GET("/ws", gin.WrapH(this.wsController.HttpHandle()))
}
