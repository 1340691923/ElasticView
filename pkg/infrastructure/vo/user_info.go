package vo

type UserInfoV2 struct {
	UserId   int      `json:"userId"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
	Perms    []string `json:"perms"`
}
