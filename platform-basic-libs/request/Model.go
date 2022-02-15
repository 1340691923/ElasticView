package request

// GmRoleModel
type GmRoleModel struct {
	ID          int      `json:"id" db:"id"`
	RoleName    string   `json:"name" db:"role_name"`
	Description string   `json:"description" db:"description"`
	RoleList    string   `json:"routes" db:"role_list"`
	Api         []string `json:"api"`
}
