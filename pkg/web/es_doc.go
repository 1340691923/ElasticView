package web

// ES文档 路由
func (this *WebServer) runEsDoc() {

	const AbsolutePath = "/api/es_doc"
	group := this.engine.Group("ES文档操作", AbsolutePath)
	{
		group.POST(false, "删除文档", "/DeleteRowByIDAction", this.esDocController.DeleteRowByIDAction)
		group.POST(false, "修改文档", "/UpdateByIDAction", this.esDocController.UpdateByIDAction)
		group.POST(false, "新增文档", "/InsertAction", this.esDocController.InsertAction)

	}
}
