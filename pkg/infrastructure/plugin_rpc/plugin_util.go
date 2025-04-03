package plugin_rpc

// 插件工具api
func (this *PluginRpcServer) runPluginUtil() {

	const AbsolutePath = "/api/plugin_util"

	engine := this.engine.GetGinEngine()

	group := engine.Group(AbsolutePath)
	{
		group.Any("/CallPlugin/:plugin_id/*action",
			this.pluginUtilController.CallPlugin)

		group.POST("/LoadDebugPlugin", this.pluginUtilController.LoadDebugPlugin)
		group.POST("/StopDebugPlugin", this.pluginUtilController.StopDebugPlugin)

		group.POST("/LiveBroadcast", this.pluginUtilController.LiveBroadcast)

		group.POST("/BatchLiveBroadcast", this.pluginUtilController.BatchLiveBroadcast)

		group.POST("/ExecSql", this.pluginUtilController.ExecSql)

		group.POST("/ExecMoreSql", this.pluginUtilController.ExecMoreSql)

		group.POST("/SelectSql", this.pluginUtilController.SelectSql)

		group.POST("/GetRoles4UserID", this.pluginUtilController.GetRoles4UserID)

		group.POST("/FirstSql", this.pluginUtilController.FirstSql)
		group.POST("/SaveDb", this.pluginUtilController.SaveDb)
		group.POST("/UpdateDb", this.pluginUtilController.UpdateDb)
		group.POST("/InsertOrUpdateDb", this.pluginUtilController.InsertOrUpdate)
		group.POST("/DeleteDb", this.pluginUtilController.DeleteDb)

		group.POST("/EsVersion", this.pluginUtilController.EsVersion)
		group.POST("/Ping", this.pluginUtilController.Ping)
		group.POST("/EsCatNodes", this.pluginUtilController.EsCatNodes)
		group.POST("/EsClusterStats", this.pluginUtilController.EsClusterStats)
		group.POST("/EsIndicesSegmentsRequest", this.pluginUtilController.EsIndicesSegmentsRequest)
		group.POST("/EsRefresh", this.pluginUtilController.EsRefresh)
		group.POST("/EsOpen", this.pluginUtilController.EsOpen)
		group.POST("/EsFlush", this.pluginUtilController.EsFlush)
		group.POST("/EsIndicesClearCache", this.pluginUtilController.EsIndicesClearCache)
		group.POST("/EsIndicesClose", this.pluginUtilController.EsIndicesClose)
		group.POST("/EsIndicesForcemerge", this.pluginUtilController.EsIndicesForcemerge)
		group.POST("/EsDeleteByQuery", this.pluginUtilController.EsDeleteByQuery)
		group.POST("/EsSnapshotCreate", this.pluginUtilController.EsSnapshotCreate)
		group.POST("/EsSnapshotDelete", this.pluginUtilController.EsSnapshotDelete)
		group.POST("/EsRestoreSnapshot", this.pluginUtilController.EsRestoreSnapshot)
		group.POST("/EsSnapshotStatus", this.pluginUtilController.EsSnapshotStatus)
		group.POST("/EsSnapshotGetRepository", this.pluginUtilController.EsSnapshotGetRepository)
		group.POST("/EsSnapshotCreateRepository", this.pluginUtilController.EsSnapshotCreateRepository)
		group.POST("/EsSnapshotDeleteRepository", this.pluginUtilController.EsSnapshotDeleteRepository)
		group.POST("/EsGetIndices", this.pluginUtilController.EsGetIndices)
		group.POST("/EsCatHealth", this.pluginUtilController.EsCatHealth)
		group.POST("/EsCatShards", this.pluginUtilController.EsCatShards)
		group.POST("/EsCatCount", this.pluginUtilController.EsCatCount)
		group.POST("/EsCatAllocationRequest", this.pluginUtilController.EsCatAllocationRequest)
		group.POST("/EsCatAliases", this.pluginUtilController.EsCatAliases)

		group.POST("/EsPerformRequest", this.pluginUtilController.EsPerformRequest)

		group.POST("/EsDelete", this.pluginUtilController.EsDelete)
		group.POST("/EsUpdate", this.pluginUtilController.EsUpdate)
		group.POST("/EsCreate", this.pluginUtilController.EsCreate)
		group.POST("/EsSearch", this.pluginUtilController.EsSearch)
		group.POST("/EsIndicesPutSettingsRequest", this.pluginUtilController.EsIndicesPutSettingsRequest)
		group.POST("/EsCreateIndex", this.pluginUtilController.EsCreateIndex)
		group.POST("/EsDeleteIndex", this.pluginUtilController.EsDeleteIndex)
		group.POST("/EsReindex", this.pluginUtilController.EsReindex)
		group.POST("/EsIndicesGetSettingsRequest", this.pluginUtilController.EsIndicesGetSettingsRequest)
		group.POST("/EsPutMapping", this.pluginUtilController.EsPutMapping)
		group.POST("/EsGetMapping", this.pluginUtilController.EsGetMapping)
		group.POST("/EsGetAliases", this.pluginUtilController.EsGetAliases)
		group.POST("/EsAddAliases", this.pluginUtilController.EsAddAliases)
		group.POST("/EsRemoveAliases", this.pluginUtilController.EsRemoveAliases)
		group.POST("/EsMoveToAnotherIndexAliases", this.pluginUtilController.EsMoveToAnotherIndexAliases)
		group.POST("/EsTaskList", this.pluginUtilController.EsTaskList)
		group.POST("/EsTasksCancel", this.pluginUtilController.EsTasksCancel)
		group.POST("/EsRunDsl", this.pluginUtilController.EsRunDsl)
		group.POST("/MysqlExecSql", this.pluginUtilController.MysqlExecSql)
		group.POST("/MysqlSelectSql", this.pluginUtilController.MysqlSelectSql)
		group.POST("/MysqlFirstSql", this.pluginUtilController.MysqlFirstSql)
		group.POST("/RedisExecCommand", this.pluginUtilController.RedisExecCommand)
		group.POST("/MongoExecCommand", this.pluginUtilController.MongoExecCommand)
		group.POST("/ShowMongoDbs", this.pluginUtilController.ShowMongoDbs)
		group.POST("/GetEveToken", this.pluginUtilController.GetEveToken)

	}
}
