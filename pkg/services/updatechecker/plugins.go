package updatechecker

import (
	"context"
	"fmt"

	"github.com/1340691923/ElasticView/pkg/services/plugin_config_service"
	"github.com/1340691923/ElasticView/pkg/services/plugin_install_service"

	"github.com/1340691923/ElasticView/pkg/services/notice_service"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
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
	availableUpdates   map[string]string
	pluginStore        manager.Service
	pluginConfigServie *plugin_config_service.PluginConfigServie
	pluginInstaller    *plugin_install_service.PluginInstaller
	enabled            bool
	evBackDao          *dao.EvBackDao
	mutex              sync.RWMutex
	log                *zap.Logger
	updateCheckURL     *url.URL
	noticeService      *notice_service.NoticeService
}

func ProvidePluginsService(log *logger.AppLogger, cfg *config.Config,
	evBackDao *dao.EvBackDao, pluginStore manager.Service, pluginInstaller *plugin_install_service.PluginInstaller,
	noticeService *notice_service.NoticeService, pluginConfigServie *plugin_config_service.PluginConfigServie) (*PluginsService, error) {
	logger := log.Named("plugins.update.checker")

	return &PluginsService{
		pluginInstaller:    pluginInstaller,
		enabled:            cfg.CheckForPluginUpdates,
		noticeService:      noticeService,
		log:                logger,
		evBackDao:          evBackDao,
		pluginStore:        pluginStore,
		pluginConfigServie: pluginConfigServie,
		availableUpdates:   make(map[string]string),
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
		s.log.Sugar().Warn("Update check failed", zap.Error(err), zap.String("所花时间", time.Since(start).String()))
		return
	}

	ctxLogger.Sugar().Infof("Update check succeeded", "duration", time.Since(start).String())
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

	go func() {
		for pluginId, pluginVersion := range availableUpdates {
			p, has := s.pluginStore.Plugin(context.Background(), pluginId)
			if has {
				pName := p.PluginData().PluginJsonData.PluginName

				if s.pluginConfigServie.IsAutoUpdateEnabled(ctx, pluginId) {
					// 发送开始自动更新的通知
					s.noticeService.LiveBroadcastEvMsg2All(&dto.NoticeData{
						Title:       fmt.Sprintf("%s插件开始自动更新", pName),
						Content:     fmt.Sprintf("%s插件检测到新版本(%s)，正在自动更新中...", pName, pluginVersion),
						Type:        "插件自动更新",
						Level:       dto.NoticeLevelInfo,
						IsTask:      true,
						FromUid:     0,
						PluginAlias: "",
						Source:      "ElasticView",
						PublishTime: time.Now(),
					})

					// 执行自动更新
					err = s.pluginInstaller.Add(ctx, pluginId, pluginVersion)

					if err != nil {
						// 更新失败的通知
						s.noticeService.LiveBroadcastEvMsg2All(&dto.NoticeData{
							Title:       fmt.Sprintf("%s插件自动更新失败", pName),
							Content:     fmt.Sprintf("%s插件自动更新到版本(%s)失败：%s", pName, pluginVersion, err.Error()),
							Type:        "插件自动更新失败",
							Level:       dto.NoticeLevelDanger,
							IsTask:      true,
							FromUid:     0,
							PluginAlias: "",
							Source:      "ElasticView",
							NoticeJumpBtn: &dto.NoticeJumpBtn{
								Text:     "跳转",
								JumpUrl:  "/plugins/manager",
								JumpType: dto.NoticeBtnJumpTypeInternal,
							},
							PublishTime: time.Now(),
						})
						s.log.Error("插件自动更新失败", zap.String("pluginId", pluginId), zap.String("pluginName", pName), zap.String("version", pluginVersion), zap.Error(err))
					} else {
						// 更新成功的通知
						s.noticeService.LiveBroadcastEvMsg2All(&dto.NoticeData{
							Title:       fmt.Sprintf("%s插件自动更新成功", pName),
							Content:     fmt.Sprintf("%s插件已成功自动更新到版本(%s)", pName, pluginVersion),
							Type:        "插件自动更新成功",
							Level:       dto.NoticeLevelSuccess,
							IsTask:      true,
							FromUid:     0,
							PluginAlias: "",
							Source:      "ElasticView",
							NoticeJumpBtn: &dto.NoticeJumpBtn{
								Text:     "跳转",
								JumpUrl:  "/plugins/manager",
								JumpType: dto.NoticeBtnJumpTypeInternal,
							},
							PublishTime: time.Now(),
						})
						s.log.Info("插件自动更新成功", zap.String("pluginId", pluginId), zap.String("pluginName", pName), zap.String("version", pluginVersion))
					}

				} else {
					// 手动更新提示（原有逻辑）
					s.noticeService.LiveBroadcastEvMsg2All(&dto.NoticeData{
						Title:       fmt.Sprintf("%s插件有更新", pName),
						Content:     fmt.Sprintf("%s插件发布了新版本(%s),请升级", pName, pluginVersion),
						Type:        "插件需更新",
						Level:       dto.NoticeLevelSuccess,
						IsTask:      true,
						FromUid:     0,
						PluginAlias: "",
						Source:      "ElasticView",
						NoticeJumpBtn: &dto.NoticeJumpBtn{
							Text:     "跳转",
							JumpUrl:  "/plugins/manager",
							JumpType: dto.NoticeBtnJumpTypeInternal,
						},
						PublishTime: time.Now(),
					})
				}
			}
		}
	}()

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
