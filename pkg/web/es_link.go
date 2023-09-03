package web

// ES连接 路由
func (this *WebServer) runEsLink() {

	const AbsolutePath = "/api/es_link"
	group := this.engine.Group("ES连接树", AbsolutePath)
	{
		group.POST(false, "新增ES连接树信息", "/InsertAction", this.esLinkController.InsertAction)
		group.POST(false, "删除ES连接树信息", "/DeleteAction", this.esLinkController.DeleteAction)
		group.POST(false, "修改ES连接树信息", "/UpdateAction", this.esLinkController.UpdateAction)
		group.POST(false, "查看ES连接树列表", "/ListAction", this.esLinkController.ListAction)
		group.POST(false, "查看ES连接配置下拉选", "/OptAction", this.esLinkController.OptAction)

	}
}
