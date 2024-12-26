package web

func (this *WebServer) runAI() {
	group := this.engine.Group("AI模块", "/api/ai")
	{
		group.Use(this.middleWareService.OperaterLog)
		group.POST(false, "查询百炼大模型", "/SearchBigMode", this.aiController.SearchBigMode)

	}
}
