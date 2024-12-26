package web

// ES 任务 路由
func (this *WebServer) runWs() {
	this.engine.GetGinEngine().GET("/ws", this.wsController.WsAction)
}
