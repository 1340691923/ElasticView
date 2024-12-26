package model

import "time"

type GmGuid struct {
	Id       int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	Uid      int       `gorm:"column:uid;NOT NULL;index:guid_name,unique" json:"uid"`
	GuidName string    `gorm:"column:guid_name;NOT NULL;index:guid_name,unique" json:"guid_name"`
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
}

func (g *GmGuid) TableName() string {
	return "gm_guid"
}
