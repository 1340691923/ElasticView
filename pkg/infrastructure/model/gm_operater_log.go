package model

import "time"

type GmOperaterLog struct {
	Id             int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	OperaterId     int       `gorm:"column:operater_id;default:0;NOT NULL" json:"operater_id"`
	OperaterName   string    `gorm:"column:operater_name;default:'';NOT NULL" json:"operater_name"`
	OperaterAction string    `gorm:"column:operater_action;default:'';NOT NULL" json:"operater_action"`
	Method         string    `gorm:"column:method;default:'';NOT NULL" json:"method"`
	Body           []byte    `gorm:"column:body;NOT NULL" json:"body"`
	Created        time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	CostTime       string    `gorm:"column:cost_time;default:'';NOT NULL" json:"cost_time"`
	Status         string    `gorm:"column:status;default:'';NOT NULL" json:"status"`
}

func (g *GmOperaterLog) TableName() string {
	return "gm_operater_log"
}
