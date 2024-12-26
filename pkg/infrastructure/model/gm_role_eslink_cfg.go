package model

import "time"

type GmRoleEslinkCfgV2 struct {
	Id          int32     `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	RoleId      int32     `gorm:"column:role_id;default:0;NOT NULL" json:"role_id"`
	EsLinkCfgId int32     `gorm:"column:es_link_cfg_id;default:0;NOT NULL" json:"es_link_cfg_id"`
	EsLinkId    int32     `gorm:"column:es_link_id;default:0;NOT NULL" json:"es_link_id"`
	Created     time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated     time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}

func (g *GmRoleEslinkCfgV2) TableName() string {
	return "gm_role_eslink_cfg_v2"
}
