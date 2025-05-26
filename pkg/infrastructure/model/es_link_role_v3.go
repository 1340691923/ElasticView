package model

import "time"

// EsLinkRoleV3 es连接与角色的关联表 (v3版本)
type EsLinkRoleV3 struct {
	Id       int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	EsLinkId int       `gorm:"column:es_link_id;default:0;NOT NULL" json:"es_link_id"`
	RoleId   int       `gorm:"column:role_id;default:0;NOT NULL" json:"role_id"`
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated  time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}

func (e *EsLinkRoleV3) TableName() string {
	return "es_link_role_v3"
}
