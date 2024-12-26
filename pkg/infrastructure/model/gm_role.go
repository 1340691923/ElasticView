package model

type GmRole struct {
	Id          int     `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	RoleName    string  `gorm:"column:role_name;default:NULL" json:"name"`
	Description string  `gorm:"column:description;default:NULL" json:"description"`
	RoleList    *string `gorm:"column:role_list" json:"role_list"`
}

func (g *GmRole) TableName() string {
	return "gm_role"
}
