package model

import "time"

type EslinkCfgModelTmp struct {
	RelationId int       `json:"relation_id" db:"relation_id"`
	ID         int       `gorm:"column:id" json:"id" db:"id"`
	RoleCfgId  int       `json:"role_cfg_id" db:"role_cfg_id"`
	EsLinkId   int       `gorm:"column:es_link_id" json:"es_link_id" db:"es_link_id"`
	User       string    `gorm:"column:user" json:"user" db:"user"`
	Pwd        string    `gorm:"column:pwd" json:"pwd" db:"pwd"`
	Created    time.Time `gorm:"column:created" json:"created" db:"created"`
	Updated    time.Time `gorm:"column:updated" json:"updated" db:"updated"`
	Remark     string    `gorm:"column:remark" json:"remark" db:"cfg_remark"`
	RootPEM    string    `gorm:"column:rootpem" json:"rootpem" db:"rootpem"`
	CertPEM    string    `gorm:"column:certpem" json:"certpem" db:"certpem"`
	KeyPEM     string    `gorm:"column:keypem" json:"keypem" db:"keypem"`
	Header     string    `gorm:"column:header" json:"header" db:"header"`
}

type EslinkCfgV2 struct {
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

func (e *EslinkCfgV2) TableName() string {
	return "eslink_cfg_v2"
}
