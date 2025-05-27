package web

func (this *WebServer) runAI() {
	group := this.engine.Group("AI模块", "/api/ai")
	{
		group.Use(this.middleWareService.OperaterLog)
		group.POST(false, "查询百炼大模型", "/SearchBigMode", this.aiController.SearchBigMode)
		group.POST(false, "获取AI配置", "/GetAIConfig", this.aiController.GetAIConfig)
		group.POST(false, "保存AI配置", "/SaveAIConfig", this.aiController.SaveAIConfig)
		group.POST(false, "测试AI连接", "/TestAIConnection", this.aiController.TestAIConnection)
	}
}
