package process

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
	"time"
)

var (
	keepPluginAliveTickerDuration = time.Second * 1
)

type Service struct {
	log *logger.AppLogger
}

func ProvideService(log *logger.AppLogger) *Service {
	return &Service{log: log}
}

func (this *Service) Start(ctx context.Context, p *plugin.Plugin) error {

	if err := startPluginAndKeepItAlive(ctx, p); err != nil {
		return err
	}

	p.Logger().Debug("成功启动" + p.PluginData().PluginJsonData.PluginName + "插件进程")
	return nil
}

func (this *Service) Stop(ctx context.Context, p *plugin.Plugin) error {
	p.Logger().Debug("Stopping plugin process")
	if err := p.Decommission(); err != nil {
		return err
	}

	if err := p.Stop(ctx); err != nil {
		return err
	}

	return nil
}

func startPluginAndKeepItAlive(ctx context.Context, p *plugin.Plugin) error {
	if err := p.Start(ctx); err != nil {
		return err
	}

	go func(p *plugin.Plugin) {
		if err := keepPluginAlive(p); err != nil {
			p.Logger().Error("Attempt to restart killed plugin process failed", "error", err)
		}
	}(p)

	return nil
}

// keepPluginAlive will restart the plugin if the process is killed or exits
func keepPluginAlive(p *plugin.Plugin) error {
	ticker := time.NewTicker(keepPluginAliveTickerDuration)

	for {
		<-ticker.C
		if p.IsDecommissioned() {
			p.Logger().Debug("Plugin decommissioned")
			return nil
		}

		if !p.Exited() {
			continue
		}

		p.Logger().Debug("Restarting plugin")
		if err := p.Start(context.Background()); err != nil {
			p.Logger().Error("Failed to restart plugin", "error", err)
			continue
		}
		p.Logger().Debug("Plugin restarted")
	}
}
