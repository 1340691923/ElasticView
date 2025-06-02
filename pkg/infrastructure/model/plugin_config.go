package model

import "time"

type PluginConfig struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PluginID   string    `gorm:"column:plugin_id;uniqueIndex;not null" json:"plugin_id"` // 插件ID
	AutoUpdate bool      `gorm:"column:auto_update;default:false" json:"auto_update"`    // 自动更新开关，默认关闭
	Created    time.Time `gorm:"column:created;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"created"`
	Updated    time.Time `gorm:"column:updated;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"updated"`
}

func (PluginConfig) TableName() string {
	return "plugin_config"
}
