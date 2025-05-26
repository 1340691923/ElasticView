package model

import "time"

// EsLinkV3 es连接信息表 (v3版本)
type EsLinkV3 struct {
	Id       int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	Ip       string    `gorm:"column:ip;default:'';NOT NULL"  json:"ip"`
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated  time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
	Remark   string    `gorm:"column:remark;default:默认连接;index:es_link_v3_remark,unique" json:"remark"`
	Version  string    `gorm:"column:version;default:'elasticsearch6.x';NOT NULL" json:"version"`
	CreateBy int       `gorm:"column:create_by;default:0;NOT NULL;index:es_link_v3_remark,unique" json:"create_by"`
}

func (e *EsLinkV3) TableName() string {
	return "es_link_v3"
}
