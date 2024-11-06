package dto

type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}

type UserUpdateReq struct {
	Id            int    `json:"id"`
	Realname      string `json:"realname"`
	RoleIds       []int  `json:"role_ids"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	WorkWechatUid string `json:"work_wechat_uid"`
}

type UserAddReq struct {
	Realname      string `json:"realname"`
	RoleIds       []int  `json:"role_ids"`
	Password      string `json:"password"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	WorkWechatUid string `json:"work_wechat_uid"`
}

type UserListReq struct {
	Appid int `json:"appid,omitempty" `
}

type DeleteUserReq struct {
	Id int `json:"id"`
}

type GetUserByIdReq struct {
	Id int `json:"id"`
}
type RolesDelReq struct {
	Id int `json:"id"`
}

type SealUserReq struct {
	Id int `json:"id"`
}

type UnSealUserReq struct {
	Id int `json:"id"`
}
