package model

import "time"

type UserRoleRelationModel struct {
	Id     int `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	UserId int `gorm:"column:user_id;default:0;index:user_role_id,unique" json:"user_id"`
	RoleId int `gorm:"column:role_id;default:0;index:user_role_id,unique" json:"role_id"`

	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_time"`
}

func (g *UserRoleRelationModel) TableName() string {
	return "user_role_relation"
}
