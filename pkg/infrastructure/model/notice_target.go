package model

import "time"

type NoticeTarget struct {
	ID       int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NoticeID int       `gorm:"column:notice_id;not null;uniqueIndex:idx_noticeid_target" json:"notice_id"` // 通知 ID
	TargetID int       `gorm:"column:target_id;uniqueIndex:idx_noticeid_target" json:"target_id"`          // 目标 ID，可为空
	Created  time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated  time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}
