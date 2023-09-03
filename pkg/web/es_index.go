package web

// ES索引 路由
func (this *WebServer) runEsIndex() {

	const AbsolutePath = "/api/es_index"
	group := this.engine.Group("ES索引", AbsolutePath)
	{
		group.POST(false, "删除ES索引", "/DeleteAction", this.esIndexController.DeleteAction)
		group.POST(false, "创建ES索引", "/CreateAction", this.esIndexController.CreateAction)
		group.POST(false, "获取ES索引配置", "/GetSettingsAction", this.esIndexController.GetSettingsAction)
		group.POST(false, "得到所有的ES索引名", "/IndexNamesAction", this.esIndexController.IndexNamesAction)
		group.POST(false, "重建索引", "/ReindexAction", this.esIndexController.ReindexAction)
		group.POST(false, "获取别名", "/GetAliasAction", this.esIndexController.GetAliasAction)
		group.POST(false, "迁移别名到其他索引", "/MoveAliasToIndex", this.esIndexController.MoveAliasToIndex)
		group.POST(false, "添加别名", "/AddAliasToIndex", this.esIndexController.AddAliasToIndex)
		group.POST(false, "批量添加别名", "/BatchAddAliasToIndex", this.esIndexController.BatchAddAliasToIndex)
		group.POST(false, "移除别名", "/RemoveAlias", this.esIndexController.RemoveAlias)

		group.POST(false, "获取所有的索引配置信息", "/GetSettingsInfoAction", this.esIndexController.GetSettingsInfoAction)
		group.POST(false, "获取索引的Stats", "/StatsAction", this.esIndexController.StatsAction)

		group.POST(false, "得到所有的索引数量", "/IndexsCountAction", this.esIndexController.IndexsCountAction)

	}
}
