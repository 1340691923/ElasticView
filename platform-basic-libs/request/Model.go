package request

import "errors"

// GmRoleModel
type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}

type TimingModel struct {
	Page   int  `json:"page"`
	Limit  int  `json:"limit"`
	Status *int `json:"status"`
	Action *int `json:"action"`
}

type DataxInfoListReq struct {
	Remark string `json:"remark"`
	Typ    string `json:"typ"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type DataxInfoInsertReq struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	Remark   string `json:"remark"`
	Typ      string `json:"typ"`
}

type DataxInfoDelReq struct {
	ID int `json:"id"`
}

type DataxInfoTestLinkReq struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
	Remark   string `json:"remark"`
	Typ      string `json:"typ"`
}

func (this DataxInfoInsertReq) Validate() (err error) {
	if this.IP == "" {
		err = errors.New("ip 不能为空")
	}
	return
}
