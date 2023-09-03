package web

// ES基础操作 路由
func (this *WebServer) runEs() {

	const AbsolutePath = "/api/es"
	group := this.engine.Group("ES基础操作", AbsolutePath)
	{
		group.POST(false, "将索引手动恢复为可写状态", "/RecoverCanWrite", this.esController.RecoverCanWrite)
		group.POST(false, "Ping ES", "/PingAction", this.esController.PingAction)
		group.POST(false, "Es的CAT操作", "/CatAction", this.esController.CatAction)
		group.POST(false, "运行DSL", "/RunDslAction", this.esController.RunDslAction)

		group.POST(false, "SQL语法转ES语法", "/SqlToDslAction", this.esController.SqlToDslAction)
		group.POST(false, "一些索引的操作", "/OptimizeAction", this.esController.OptimizeAction)

	}
}
