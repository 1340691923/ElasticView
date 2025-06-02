package plugin_config_service

import (
	"context"

	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"go.uber.org/zap"
)

type PluginConfigServie struct {
	log             *logger.AppLogger
	pluginConfigDao *dao.PluginConfigDao
}

func NewPluginConfigServie(log *logger.AppLogger, pluginConfigDao *dao.PluginConfigDao) *PluginConfigServie {
	return &PluginConfigServie{log: log, pluginConfigDao: pluginConfigDao}
}

func (this *PluginConfigServie) GetPluginConfig(ctx context.Context, pluginID string) (*model.PluginConfig, error) {
	return this.pluginConfigDao.GetPluginConfig(ctx, pluginID)
}

// UpdatePluginConfig 更新插件配置
func (this *PluginConfigServie) UpdatePluginConfig(ctx context.Context, pluginID string, autoUpdate bool) error {
	return this.pluginConfigDao.UpdatePluginConfig(ctx, pluginID, autoUpdate)
}

// IsAutoUpdateEnabled 检查插件是否启用自动更新
func (this *PluginConfigServie) IsAutoUpdateEnabled(ctx context.Context, pluginID string) bool {
	config, err := this.pluginConfigDao.GetPluginConfig(ctx, pluginID)
	if err != nil {
		this.log.Error("获取插件配置失败", zap.Error(err), zap.String("pluginID", pluginID))
		return false
	}
	return config.AutoUpdate
}
