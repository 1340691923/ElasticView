package web

func (this *WebServer) runNotice() {
	const AbsolutePath = "/api/notice"
	group := this.engine.Group("通知模块", AbsolutePath)
	{
		group.Use(this.middleWareService.OperaterLog)
		group.POST(false, "查询通知列表", "/GetList", this.noticeController.GetList)
		group.POST(false, "标记通知为已读", "/MarkReadNotice", this.noticeController.MarkReadNotice)
		group.POST(true, "清空通知数据", "/Truncate", this.noticeController.Truncate)
	}
}
