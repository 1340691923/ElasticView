package web

// BI用户 路由
func (this *WebServer) runManagerUser() {
	group := this.engine.Group("用户/角色管理模块", "/api/gm_user")
	{
		group.POST(false, "查询用户信息", "/info", this.managerUserController.UserInfo)
		group.POST(false, "GM角色列表", "/roles", this.managerRoleController.RolesAction)
		group.POST(false, "退出登录", "/logout", this.managerUserController.LogoutAction)
		group.POST(false, "查询用户列表", "/userlist", this.managerUserController.UserListAction)
		group.POST(false, "角色下拉选", "/roleOption", this.managerRoleController.RoleOptionAction)
		group.POST(false, "通过ID获取用户信息", "/getUserById", this.managerUserController.GetUserByIdAction)
		group.POST(false, "获取接口路由信息", "/UrlConfig", this.managerUserController.UrlConfig)
		group.Use(this.middleWareService.OperaterLog)

		group.POST(true, "修改自己的密码", "/ModifyPwd", this.managerUserController.ModifyPwd)
		group.POST(true, "修改角色", "/role/update", this.managerRoleController.RolesUpdateAction)
		group.POST(true, "新增角色", "role/add", this.managerRoleController.RolesAddAction)
		group.POST(true, "删除角色", "/role/delete", this.managerRoleController.RolesDelAction)
		group.POST(true, "修改用户", "/UpdateUser", this.managerUserController.UserUpdateAction)
		group.POST(true, "新增用户", "/InsertUser", this.managerUserController.UserAddAction)
		group.POST(true, "删除用户", "/DelUser", this.managerUserController.DeleteUserAction)
	}
}
