package model

// EsLinkModel es连接信息表
type EsLinkModel struct {
	ID      int64  `gorm:"column:id" json:"id" db:"id"`
	Ip      string `gorm:"column:ip" json:"ip" db:"ip"`
	User    string `gorm:"column:user" json:"user" db:"user"`
	Pwd     string `gorm:"column:pwd" json:"pwd" db:"pwd"`
	Created string `gorm:"column:created" json:"created" db:"created"`
	Updated string `gorm:"column:updated" json:"updated" db:"updated"`
	Remark  string `gorm:"column:remark" json:"remark" db:"remark"`
	Version int    `json:"version" db:"version"`
	RootPEM string `gorm:"column:rootpem" json:"rootpem" db:"rootpem"`
	CertPEM string `gorm:"column:certpem" json:"certpem" db:"certpem"`
	KeyPEM  string `gorm:"column:keypem" json:"keypem" db:"keypem"`
}
