package dao

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

type PluginConfigDao struct {
	log *logger.AppLogger
	orm *orm.Gorm
}

func NewPluginConfigDao(log *logger.AppLogger, orm *orm.Gorm) *PluginConfigDao {
	return &PluginConfigDao{log: log, orm: orm}
}

// GetPluginConfig 获取插件配置，如果不存在则创建默认配置（原子操作）
func (this *PluginConfigDao) GetPluginConfig(ctx context.Context, pluginID string) (*model.PluginConfig, error) {
	// 先尝试插入默认配置（如果不存在）
	defaultConfig := model.PluginConfig{
		PluginID:   pluginID,
		AutoUpdate: false, // 默认关闭自动更新
	}

	// 使用 INSERT ON DUPLICATE KEY UPDATE 确保记录存在
	err := this.orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}},
		DoNothing: true, // 如果已存在则不做任何操作
	}).Create(&defaultConfig).Error

	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 再查询获取最新配置
	var config model.PluginConfig
	err = this.orm.WithContext(ctx).Where("plugin_id = ?", pluginID).First(&config).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}

// UpdatePluginConfig 更新插件配置（原子操作）
func (this *PluginConfigDao) UpdatePluginConfig(ctx context.Context, pluginID string, autoUpdate bool) error {
	config := model.PluginConfig{
		PluginID:   pluginID,
		AutoUpdate: autoUpdate,
	}

	// 使用 INSERT ON DUPLICATE KEY UPDATE 原子性地插入或更新
	err := this.orm.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "plugin_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"auto_update"}),
	}).Create(&config).Error

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
