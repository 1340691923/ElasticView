package model

import (
	"github.com/1340691923/ElasticView/engine/db"
)

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
}

var EsLinkList = []EsLinkModel{}

//刷新eslink表数据到内存
func (this *EsLinkModel) FlushEsLinkList() (err error) {
	list, err := this.GetListAction()
	if err != nil {
		return
	}
	EsLinkList = list
	return
}

//获取列表信息
func (this *EsLinkModel) GetListAction() (esLinkList []EsLinkModel, err error) {
	sql, args, err := db.SqlBuilder.
		Select("*").
		From("es_link").ToSql()
	if err != nil {
		return
	}

	err = db.Sqlx.Select(&esLinkList, sql, args...)
	if err != nil {
		return
	}
	return
}
