package web

// ES mapping 路由
func (this *WebServer) runEsMap() {

	const AbsolutePath = "/api/es_map"
	group := this.engine.Group("ES映射", "/api/es_map")
	{
		group.POST(false, "查看mapping列表", "/ListAction", this.esMappingController.ListAction)
		group.POST(false, "修改mapping", "/UpdateMappingAction", this.esMappingController.UpdateMappingAction)

	}
}
