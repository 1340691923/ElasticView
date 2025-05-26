package web

// ES连接 路由
func (this *WebServer) runEsLink() {

	const AbsolutePath = "/api/es_link"
	group := this.engine.Group("数据源", AbsolutePath)
	{
		group.POST(false, "查看连接配置下拉选", "/OptAction", this.esLinkController.OptAction)
		group.POST(false, "查看数据源列表", "/ListAction", this.esLinkController.ListAction)
		group.POST(false, "查看鉴权列表", "/GetEsCfgList", this.esLinkController.GetEsCfgList)
		group.POST(false, "查看鉴权配置下拉选", "/GetEsCfgOpt", this.esLinkController.GetEsCfgOpt)

		group.Use(this.middleWareService.OperaterLog)

		group.POST(false, "按id删除连接鉴权配置", "/DeleteEsCfgRelation", this.esLinkController.DeleteEsCfgRelation)

		group.POST(true, "新增数据源信息", "/InsertAction", this.esLinkController.InsertAction)
		group.POST(true, "删除数据源信息", "/DeleteAction", this.esLinkController.DeleteAction)
		group.POST(true, "修改数据源信息", "/UpdateAction", this.esLinkController.UpdateAction)

		group.POST(true, "新增鉴权配置信息", "/InsertEsCfgAction", this.esLinkController.InsertEsCfgAction)
		group.POST(true, "修改鉴权配置信息", "/UpdateEsCfgAction", this.esLinkController.UpdateEsCfgAction)
		group.POST(true, "删除鉴权配置信息", "/DeleteEsCfgAction", this.esLinkController.DeleteEsCfgAction)

	}
}
