package model

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
)

type GmOperaterLog struct {
	ID             int64  `gorm:"column:id" json:"id" form:"id"`
	OperaterName   string `gorm:"column:operater_name" json:"operater_name" form:"operater_name"`
	OperaterId     int64  `gorm:"column:operater_id" json:"operater_id" form:"operater_id"`
	OperaterAction string `gorm:"column:operater_action" json:"operater_action" form:"operater_action"`
	Created        string `gorm:"column:created" json:"created" form:"created"`
	Method         string `gorm:"column:method" json:"method" form:"method"`
	Parmas         string `gorm:"column:parmas" json:"parmas" form:"parmas"`
	Body           string `gorm:"column:body" json:"body" form:"body"`
	OperaterRoleId int    `gorm:"column:operater_role_id" json:"operater_role_id" form:"operater_role_id"`
}

func (this *GmOperaterLog) TableName() string {
	return "gm_operater_log"
}

func (this *GmOperaterLog) Insert() (err error) {
	body, err := util.GzipCompress(this.Body)
	if err != nil {
		return
	}
	_, err = db.SqlBuilder.Insert(this.TableName()).SetMap(map[string]interface{}{
		"operater_name":    this.OperaterName,
		"operater_id":      this.OperaterId,
		"operater_action":  this.OperaterAction,
		"method":           this.Method,
		"parmas":           this.Parmas,
		"body":             body,
		"operater_role_id": this.OperaterRoleId,
	}).RunWith(db.Sqlx).Exec()
	if err != nil {
		return
	}
	return
}
