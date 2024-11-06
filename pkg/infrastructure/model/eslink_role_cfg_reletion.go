package model

import "time"

type EslinkRoleCfgReletion struct {
	Id        int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	EsLinkId  int       `gorm:"column:es_link_id;default:0;NOT NULL" json:"es_link_id"`
	RoleCfgId int       `gorm:"column:role_cfg_id;default:0;NOT NULL" json:"role_cfg_id"`
	Created   time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated   time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}

func (e *EslinkRoleCfgReletion) TableName() string {
	return "eslink_role_cfg_reletion"
}
