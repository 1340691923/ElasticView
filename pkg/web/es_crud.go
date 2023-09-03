package web

// ES基础操作 路由
func (this *WebServer) runEsCrud() {

	const AbsolutePath = "/api/es_crud"
	group := this.engine.Group("Navicat", AbsolutePath)
	{
		group.POST(false, "数据筛选", "/GetList", this.esCrudController.GetList)
		group.POST(false, "获取查询语句", "/GetDSL", this.esCrudController.GetDSL)
		group.POST(false, "下载数据", "/Download", this.esCrudController.Download)
	}
}
