package model

import (
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"time"
)

type GmTimedList struct {
	ID             int64  `gorm:"column:id" json:"id" db:"id"`
	Action         int    `gorm:"column:action" json:"action" db:"action"`
	ExecTime       int64  `gorm:"column:exec_time" json:"exec_time" db:"exec_time"`
	ExecTimeFormat string `gorm:"column:exec_time_format" json:"exec_time_format" db:"exec_time_format"`
	Status         int64  `gorm:"column:status" json:"status" db:"status"` //0等待执行 1正在执行 2成功 3失败 4取消
	Msg            string `gorm:"column:msg" json:"msg" db:"msg"`
	Extra          string `gorm:"column:extra" json:"extra" db:"extra"`
	Data           string `gorm:"column:data" json:"data" db:"data"`
	CreateTime     string `gorm:"column:create_time" json:"create_time" db:"created"`
	UpdateTime     string `gorm:"column:update_time" json:"update_time" db:"updated"`
	TaskId         string `gorm:"column:task_id" json:"task_id" db:"task_id"`
}

const (
	GmTimedListWait   = 0
	GmTimedListStart  = 1
	GmTimedListSucc   = 2
	GmTimedListErr    = 3
	GmTimedListCancel = 4
)

func (this *GmTimedList) TableName() string {
	return `gm_timed_list`
}

func (this *GmTimedList) UpdateStatus(status int, msg string, taskID string) (err error) {
	setMap := util.Map{
		"status":  status,
		"msg":     msg,
		"updated": time.Now().Format(util.TimeFormat),
	}
	_, err = db.SqlBuilder.
		Update(this.TableName()).
		SetMap(setMap).
		Where(db.Eq{"task_id": taskID}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *GmTimedList) GetTaskList() (gmTimedList []GmTimedList, err error) {
	sql, args, err := db.SqlBuilder.Select("*").
		From(this.TableName()).
		Where(db.Gte{"exec_time": time.Now().Unix()}).
		Where(db.Or{
			db.Eq{"status": GmTimedListWait},
			db.Eq{"status": GmTimedListStart},
		}).ToSql()
	if err != nil {
		return
	}
	err = db.Sqlx.Select(&gmTimedList, sql, args...)
	if err != nil {
		return
	}
	return
}

func (this *GmTimedList) AddTask(taskid, execTimeFormat, data string, action int, execTime int64) (err error) {
	_, err = db.SqlBuilder.
		Insert(this.TableName()).
		SetMap(util.Map{
			"task_id":          taskid,
			"action":           action,
			"exec_time_format": execTimeFormat,
			"data":             data,
			"exec_time":        execTime,
			"created":          time.Now().Format(util.TimeFormat),
			"updated":          time.Now().Format(util.TimeFormat),
		}).RunWith(db.Sqlx).Exec()
	if err != nil {
		return
	}
	return
}
