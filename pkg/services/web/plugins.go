package web

// ES基础操作 路由
func (this *WebServer) runPlugins() {

	const AbsolutePath = "/api/plugins"
	group := this.engine.Group("插件操作", AbsolutePath)
	{
		group.POST(true, "离线上传插件", "/UploadPlugin", this.pluginController.UploadPlugin)
		group.Use(this.middleWareService.OperaterLog)
		group.POST(false, "查询本地插件列表", "/GetLocalPluginList", this.pluginController.GetLocalPluginList)
		group.POST(false, "导入EvKey到配置文件", "/ImportEvKey", this.pluginController.ImportEvKey)
		group.POST(false, "查询插件动态", "/GetWxArticleList", this.pluginController.GetWxArticleList)
		group.POST(false, "查询插件市场", "/PluginMarket", this.pluginController.PluginMarket)
		group.POST(false, "查询插件详情", "/GetPluginInfo", this.pluginController.GetPluginInfo)
		group.POST(true, "安装插件", "/InstallPlugin", this.pluginController.InstallPlugin)
		group.POST(true, "卸载插件", "/UnInstallPlugin", this.pluginController.UnInstallPlugin)
		group.POST(false, "给插件star", "/StarPlugin", this.pluginController.StarPlugin)
		group.POST(false, "给插件评论", "/AddComment", this.pluginController.AddComment)
		group.POST(false, "给插件评论点赞", "/LikeComment", this.pluginController.LikeComment)
		group.POST(false, "给插件star", "/ListComments", this.pluginController.ListComments)

	}
}
