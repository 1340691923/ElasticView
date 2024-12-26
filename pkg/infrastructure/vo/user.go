package vo

type User struct {
	Token    string `json:"token"`
	UnixTime int64  `json:"unix_time"`
}

type GmUsers struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`

	Avatar        string `json:"avatar"`
	IsBan         int32  `json:"is_ban"`
	Realname      string `json:"realname"`
	Email         string `json:"email"`
	WorkWechatUid string `json:"work_wechat_uid"`
	RoleIds       []int  `json:"role_ids"`
	UpdateTime    string `json:"update_time"`
	CreateTime    string `json:"create_time"`
	LastLoginTime string `json:"last_login_time"`
}
