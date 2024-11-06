package model

import "time"

type DataxTransferList struct {
	Id          int       `gorm:"column:id;primary_key;NOT NULL" json:"id"`
	FormData    string    `gorm:"column:form_data;default:NULL" json:"form_data"`
	Remark      string    `gorm:"column:remark;default:;NOT NULL;index:datax_transfer_list_remark,unique" json:"remark"`
	TableName2  string    `gorm:"column:table_name;default:'';NOT NULL" json:"table_name"`
	IndexName   string    `gorm:"column:index_name;default:'';NOT NULL" json:"index_name"`
	ErrorMsg    string    `gorm:"column:error_msg;default:'';NOT NULL" json:"error_msg"`
	CrontabSpec string    `gorm:"column:crontab_spec;default:'';NOT NULL" json:"crontab_spec"`
	Dbcount     int       `gorm:"column:dbcount;default:0;NOT NULL" json:"dbcount"`
	Escount     int       `gorm:"column:escount;default:0;NOT NULL" json:"escount"`
	EsConnect   int       `gorm:"column:es_connect;default:0;NOT NULL;index:datax_transfer_list_remark,unique" json:"es_connect"`
	Status      string    `gorm:"column:status;default:'任务运行中...';NOT NULL" json:"status"`
	Updated     time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
	Created     time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
}

func (d *DataxTransferList) TableName() string {
	return "datax_transfer_list"
}
