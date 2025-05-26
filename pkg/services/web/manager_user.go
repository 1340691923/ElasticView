package web

// BI用户 路由
func (this *WebServer) runManagerUser() {
	group := this.engine.Group("权限模块", "/api/gm_user")
	{
		group.Any(false, "获取前端路由配置", "/GetRoutesConfig", this.managerUserController.GetRoutesConfig)
		group.Any(false, "查询用户信息", "/info", this.managerUserController.UserInfo)
		group.Any(false, "查询用户信息V2", "/infoV2", this.managerUserController.UserInfoV2)
		group.POST(false, "权限组下拉选", "/roleOption", this.managerRoleController.RoleOptionAction)
		group.Use(this.middleWareService.OperaterLog)
		group.Any(true, "查询第三方认证配置", "/GetOAuthConfigs", this.managerUserController.GetOAuthConfigs)
		group.Any(true, "设置第三方认证配置", "/SaveOAuthConfigs", this.managerUserController.SaveOAuthConfigs)
		group.POST(false, "GM权限组列表", "/roles", this.managerRoleController.RolesAction)
		group.POST(false, "退出登录", "/logout", this.managerUserController.LogoutAction)
		group.POST(false, "查询用户列表", "/userlist", this.managerUserController.UserListAction)
		group.POST(false, "通过ID获取用户信息", "/getUserById", this.managerUserController.GetUserByIdAction)
		group.POST(false, "获取接口路由信息", "/UrlConfig", this.managerUserController.UrlConfig)
		group.POST(false, "修改自己的密码", "/ModifyPwd", this.managerUserController.ModifyPwd)
		group.POST(true, "修改用户的密码", "/ModifyPwdByUserId", this.managerUserController.ModifyPwdByUserId)
		group.POST(true, "修改权限组", "/role/update", this.managerRoleController.RolesUpdateAction)
		group.POST(true, "新增权限组", "role/add", this.managerRoleController.RolesAddAction)
		group.POST(true, "删除权限组", "/role/delete", this.managerRoleController.RolesDelAction)
		group.POST(true, "修改用户", "/UpdateUser", this.managerUserController.UserUpdateAction)
		group.POST(true, "新增用户", "/InsertUser", this.managerUserController.UserAddAction)
		group.POST(true, "删除用户", "/DelUser", this.managerUserController.DeleteUserAction)
		group.POST(true, "封禁用户", "/SealUserAction", this.managerUserController.SealUserAction)
		group.POST(true, "解封用户", "/UnSealUserAction", this.managerUserController.UnSealUserAction)

	}
}
