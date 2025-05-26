package model

import "time"

type Notice struct {
	ID          int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"column:title" json:"title"`                                     //标题
	Content     string    `gorm:"column:content;type:text" json:"content"`                       //消息内容
	Type        string    `gorm:"column:type;default:'system'" json:"type"`                      // system / alert / announcement
	Level       string    `gorm:"column:level;default:'info'" json:"level"`                      // info / warn / error
	IsTask      int       `gorm:"column:is_task;default:0" json:"is_task"`                       // 是否任务
	FromUid     int       `gorm:"column:from_uid;default:0" json:"from_uid"`                     // 发布者用户 ID
	PluginAlias string    `gorm:"column:plugin_alias;default:''" json:"plugin_alias"`            // 插件 ID
	Source      string    `gorm:"column:source;default:''" json:"source"`                        // 来源信息
	BtnDesc     string    `gorm:"column:btn_desc;default:''" json:"btn_desc"`                    // 按钮文案
	BtnJumpUrl  string    `gorm:"column:btn_jump_url;default:''" json:"btn_jump_url"`            // 跳转链接
	BtnJumpType string    `gorm:"column:btn_jump_type;default:'internal'" json:"btn_jump_type"`  // internal / external
	TargetType  string    `gorm:"column:target_type;index:idx_notice_target" json:"target_type"` // 目标类型：all / role / user
	Created     time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated     time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
	IsRead      bool      `gorm:"-" json:"is_read"`
}
