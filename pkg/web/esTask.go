package web

// ES 任务 路由
func (this *WebServer) runEsTask() {

	const AbsolutePath = "/api/es_task"
	group := this.engine.Group("ES映射", AbsolutePath)
	{
		group.POST(false, "查看任务列表", "/ListAction", this.esTaskController.ListAction)
		group.POST(false, "取消任务", "/CancelAction", this.esTaskController.CancelAction)
	}
}
