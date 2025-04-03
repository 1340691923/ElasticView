package pluginstore

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/backendplugin/provider"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager/process"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type PluginStoreService struct {
	pluginRegistry manager.Service
	cfg            *config.Config
	progressSvr    *process.Service
	log            *logger.AppLogger
	orm            *orm.Gorm
}

func NewPluginStoreService(pluginRegistry manager.Service, cfg *config.Config, progressSvr *process.Service, log *logger.AppLogger, orm *orm.Gorm) *PluginStoreService {
	return &PluginStoreService{pluginRegistry: pluginRegistry, cfg: cfg, progressSvr: progressSvr, orm: orm, log: log.Named("pluginStore")}
}

func (this *PluginStoreService) Run(ctx context.Context) error {

	err := this.init(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	<-ctx.Done()
	this.shutdown(ctx)
	return ctx.Err()
}

func (this *PluginStoreService) init(ctx context.Context) (err error) {

	if err = this.cfg.Plugin.Error(); err != nil {
		return errors.WithStack(err)
	}

	err = filepath.Walk(this.cfg.Plugin.LoadPath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		return this.FastInitPlugin(ctx, info.Name())
	})

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *PluginStoreService) shutdown(ctx context.Context) {
	var wg sync.WaitGroup
	for _, v := range s.pluginRegistry.Plugins(ctx) {
		wg.Add(1)
		go func(ctx context.Context, plugin *plugin.Plugin) {
			defer wg.Done()
			s.FastShutdown(ctx, plugin)
		}(ctx, v)
	}
	wg.Wait()
}

func (this *PluginStoreService) FastInitPlugin(ctx context.Context, fileName string) error {
	buildOs := runtime.GOOS
	buildArch := runtime.GOARCH

	if !strings.Contains(fileName, buildOs) || !strings.Contains(fileName, buildArch) {
		return nil
	}

	tmpArr := strings.Split(fileName, "_")

	if len(tmpArr) == 0 {
		return nil
	}
	pluginAlias := tmpArr[0]
	pluginLog, pluginLogPath, closeLogWrite, err := logger.InitPluginLog(this.cfg, pluginAlias)
	if err != nil {
		return errors.WithStack(err)
	}

	p := provider.DefaultProvider(ctx, pluginLog, pluginLogPath, closeLogWrite, &provider.Config{
		ID:             pluginAlias,
		PluginDir:      this.cfg.Plugin.LoadPath,
		PluginFileName: fileName,
		ExecArgs: []string{
			fmt.Sprintf("-tmpFileStorePath=%s", this.cfg.GetStorePath(pluginAlias)),
			fmt.Sprintf("-evRpcPort=%d", this.cfg.PluginRpcPort),
		},
	}, this.cfg, this.orm)
	err = this.progressSvr.Start(ctx, p)
	if err != nil {
		this.log.Error("插件启动失败", zap.Error(err))
		return errors.WithStack(err)
	}

	err = this.pluginRegistry.AddPlugin(ctx, p)

	if err != nil {
		this.log.Error("插件添加失败", zap.Error(err))
		return err
	}

	return nil

}

func (this *PluginStoreService) FastReloadPlugin(ctx context.Context, fileName string) error {
	buildOs := runtime.GOOS
	buildArch := runtime.GOARCH

	if !strings.Contains(fileName, buildOs) || !strings.Contains(fileName, buildArch) {
		return nil
	}

	tmpArr := strings.Split(fileName, "_")

	if len(tmpArr) == 0 {
		return nil
	}
	pluginAlias := tmpArr[0]
	pluginLog, pluginLogPath, closeLogWrite, err := logger.InitPluginLog(this.cfg, pluginAlias)
	if err != nil {
		return errors.WithStack(err)
	}

	p := provider.DefaultProvider(ctx, pluginLog, pluginLogPath, closeLogWrite, &provider.Config{
		ID:             pluginAlias,
		PluginDir:      this.cfg.Plugin.LoadPath,
		PluginFileName: fileName,
		ExecArgs: []string{
			fmt.Sprintf("-tmpFileStorePath=%s", this.cfg.GetStorePath(pluginAlias)),
			fmt.Sprintf("-evRpcPort=%d", this.cfg.PluginRpcPort),
		},
	}, this.cfg, this.orm)

	err = this.pluginRegistry.Reload(ctx, p)

	if err != nil {
		this.log.Error("插件添加失败", zap.Error(err))
		return err
	}

	return nil
}

func (this *PluginStoreService) FastShutdown(ctx context.Context, plugin *plugin.Plugin) (err error) {

	plugin.Logger().Debug("Stopping plugin", plugin.ID)

	if err = this.progressSvr.Stop(ctx, plugin); err != nil {
		plugin.Logger().Error("Failed to stop plugin", "error", err)
	}

	return err
}

func (this *PluginStoreService) FastRemove(ctx context.Context, plugin *plugin.Plugin) (err error) {

	plugin.Logger().Debug("Plugin remove start", plugin.ID)
	if err = this.pluginRegistry.Remove(ctx, plugin.ID); err != nil {
		plugin.Logger().Error("Failed to remove plugin", "error", err)
	}
	plugin.Logger().Debug("Plugin remove end", plugin.ID)

	return err
}
