package web

// ES 新手引导 路由
func (this *WebServer) runGmGuid() {

	const AbsolutePath = "/api/gm_guid"
	group := this.engine.Group("Ev新手引导", AbsolutePath)
	{
		group.POST(false, "完成新手引导", "/Finish", this.guidController.Finish)
		group.POST(false, "判断是否完成新手引导", "/IsFinish", this.guidController.IsFinish)
	}
}
