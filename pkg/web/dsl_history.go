package web

// DslHistory 路由
func (this *WebServer) runDslHistory() {

	const AbsolutePath = "/api/dslHistory"
	dslHistory := this.engine.Group("ES查询历史记录", AbsolutePath)
	{
		dslHistory.POST(false, "清空DSL查询历史记录", "/CleanAction", this.dslHistoryController.CleanAction)
		dslHistory.POST(false, "查看DSL查询历史记录", "/ListAction", this.dslHistoryController.ListAction)
	}
}
