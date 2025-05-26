package model

import "time"

// EsLinkCfgRelV3 es连接与连接配置的关联表 (v3版本)
type EsLinkCfgRelV3 struct {
	Id        int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	EsLinkId  int       `gorm:"column:es_link_id;default:0;NOT NULL" json:"es_link_id"`
	EsCfgId   int       `gorm:"column:es_cfg_id;default:0;NOT NULL" json:"es_cfg_id"`
	IsDefault bool      `gorm:"column:is_default;default:false;NOT NULL" json:"is_default"`
	Created   time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated   time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}

func (e *EsLinkCfgRelV3) TableName() string {
	return "es_link_cfg_rel_v3"
}

// 保持 EsLinkCfgV3 与 EslinkCfgV2 结构一致
type EsLinkCfgV3 struct {
	Id       int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	User     string    `gorm:"column:user;default:'';NOT NULL" json:"user"`
	Pwd      string    `gorm:"column:pwd;default:'';NOT NULL" json:"pwd"`
	Rootpem  *string   `gorm:"column:rootpem" json:"rootpem"`
	Certpem  *string   `gorm:"column:certpem" json:"certpem"`
	Keypem   *string   `gorm:"column:keypem" json:"keypem"`
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated  time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
	CreateBy int       `gorm:"column:create_by;NOT NULL" json:"create_by"`
	Remark   string    `gorm:"column:remark;default:默认配置" json:"remark"`
	Header   string    `gorm:"column:header;default:'[]';NOT NULL" json:"header"`
}

func (e *EsLinkCfgV3) TableName() string {
	return "es_link_cfg_v3"
} 