package dto

type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}

type UserUpdateReq struct {
	Id       int    `json:"id"`
	Realname string `json:"realname"`
	RoleId   int32  `json:"role_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserAddReq struct {
	Realname string `json:"realname"`
	RoleId   int32  `json:"role_id"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserListReq struct {
	Appid int `json:"appid,omitempty" `
}

type DeleteUserReq struct {
	Id int32 `json:"id"`
}

type GetUserByIdReq struct {
	Id int32 `json:"id"`
}
type RolesDelReq struct {
	Id int `json:"id"`
}
