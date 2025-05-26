package model

import "time"

type NoticeReadLog struct {
	ID       int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID   int       `gorm:"column:user_id;not null;uniqueIndex:uniq_user_notice" json:"user_id"`
	NoticeID int       `gorm:"column:notice_id;not null;uniqueIndex:uniq_user_notice" json:"notice_id"`
	ReadAt   time.Time `gorm:"column:read_at;autoCreateTime" json:"read_at"`
}
