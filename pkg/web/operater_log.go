package web

func (this *WebServer) runOperaterLog() {
	group := this.engine.Group("操作记录模块", "/api/operater_log")
	{
		group.POST(true, "查看用户操作记录", "/ListAction", this.gmOperaterController.ListAction)
	}
}
