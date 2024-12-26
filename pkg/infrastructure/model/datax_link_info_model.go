package model

import "time"

type DataxLinkInfo struct {
	Id       int       `gorm:"column:id;primaryKey;autoIncrement;NOT NULL" json:"id"`
	Ip       string    `gorm:"column:ip;default:'';NOT NULL" json:"ip"`
	Port     int       `gorm:"column:port;default:0;NOT NULL" json:"port"`
	DbName   string    `gorm:"column:db_name;default:'';NOT NULL" json:"db_name"`
	Username string    `gorm:"column:username;default:'';NOT NULL" json:"username"`
	Pwd      string    `gorm:"column:pwd;default:'';NOT NULL" json:"pwd"`
	Remark   string    `gorm:"column:remark;default:;NOT NULL;index:link_remark_uniq,unique" json:"remark"`
	Typ      string    `gorm:"column:typ;default:;NOT NULL;index:link_remark_uniq,unique" json:"typ"`
	Updated  time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
}

func (d *DataxLinkInfo) TableName() string {
	return "datax_link_info"
}
