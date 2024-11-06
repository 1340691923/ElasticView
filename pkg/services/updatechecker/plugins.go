package updatechecker

import (
	"context"
	"github.com/spf13/cast"

	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"go.uber.org/zap"

	"net/url"
	"sync"
	"time"

	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
)

type PluginsService struct {
	availableUpdates map[string]string
	pluginStore      manager.Service
	enabled          bool
	evBackDao        *dao.EvBackDao
	mutex            sync.RWMutex
	log              *zap.Logger
	updateCheckURL   *url.URL
}

func ProvidePluginsService(log *logger.AppLogger, cfg *config.Config, evBackDao *dao.EvBackDao, pluginStore manager.Service) (*PluginsService, error) {
	logger := log.Named("plugins.update.checker")

	return &PluginsService{
		enabled: cfg.CheckForPluginUpdates,

		log:              logger,
		evBackDao:        evBackDao,
		pluginStore:      pluginStore,
		availableUpdates: make(map[string]string),
	}, nil
}

func (s *PluginsService) IsDisabled() bool {
	return !s.enabled
}

func (s *PluginsService) Run(ctx context.Context) error {
	s.instrumentedCheckForUpdates(ctx)

	ticker := time.NewTicker(time.Minute * 10)
	run := true

	for run {
		select {
		case <-ticker.C:
			s.instrumentedCheckForUpdates(ctx)
		case <-ctx.Done():
			run = false
		}
	}

	return ctx.Err()
}

func (s *PluginsService) HasUpdate(ctx context.Context, pluginID string) (string, bool) {
	s.mutex.RLock()
	updateVers, updateAvailable := s.availableUpdates[pluginID]
	s.mutex.RUnlock()
	if updateAvailable {
		plugin, exists := s.pluginStore.Plugin(ctx, pluginID)
		if !exists {
			return "", false
		}

		if canUpdate(plugin.Version(), updateVers) {
			return updateVers, true
		}
	}

	return "", false
}

func (s *PluginsService) instrumentedCheckForUpdates(ctx context.Context) {
	start := time.Now()
	ctxLogger := s.log
	if err := s.checkForUpdates(ctx); err != nil {
		ctxLogger.Sugar().Warnf("Update check failed", "error", err, "duration", time.Since(start))
		return
	}

	ctxLogger.Sugar().Infof("Update check succeeded", "duration", time.Since(start))
}

func (s *PluginsService) InstrumentedCheckForUpdates(ctx context.Context) {
	s.instrumentedCheckForUpdates(ctx)
}

func (s *PluginsService) checkForUpdates(ctx context.Context) error {
	ctxLogger := s.log
	ctxLogger.Debug("Preparing plugins eligible for version check")
	localPlugins := s.pluginsEligibleForVersionCheck(ctx)

	ctxLogger.Sugar().Debugf("Checking for plugin updates")

	gcomPlugins, err := s.evBackDao.GetEvPluginsMaxVersion(ctx, s.pluginIDsCSV(localPlugins))

	if err != nil {
		return err
	}

	availableUpdates := map[string]string{}
	for slug, version := range gcomPlugins {
		if localP, exists := localPlugins[slug]; exists {
			if canUpdate(localP.Version(), cast.ToString(version)) {
				availableUpdates[localP.ID] = cast.ToString(version)
			}
		}
	}

	if len(availableUpdates) > 0 {
		s.mutex.Lock()
		s.availableUpdates = availableUpdates
		s.mutex.Unlock()
	}

	return nil
}

func (s *PluginsService) pluginIDsCSV(m map[string]*plugin.Plugin) []string {
	ids := make([]string, 0, len(m))
	for pluginID := range m {
		ids = append(ids, pluginID)
	}

	return ids
}

func (s *PluginsService) pluginsEligibleForVersionCheck(ctx context.Context) map[string]*plugin.Plugin {
	result := make(map[string]*plugin.Plugin)
	for _, p := range s.pluginStore.Plugins(ctx) {
		if p.BackendDebug() {
			continue
		}

		result[p.ID] = p
	}

	return result
}
