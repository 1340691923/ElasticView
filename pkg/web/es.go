package web

// ES基础操作 路由
func (this *WebServer) runEs() {

	const AbsolutePath = "/api/es"
	group := this.engine.Group("ES基础操作", AbsolutePath)
	{

		group.POST(false, "Ping连接", "/PingAction", this.esController.PingAction)
		group.POST(false, "获取ES索引数", "/IndexsCountAction", this.esController.IndexsCountAction)
		group.POST(false, "进行ES的Cat操作", "/CatAction", this.esController.CatAction)

	}
}
