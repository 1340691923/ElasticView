package updatechecker

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/services/notice_service"
	"github.com/1340691923/eve-plugin-sdk-go/ev_api/dto"
	"github.com/hashicorp/go-version"
	"go.uber.org/zap"
	"sync"
	"time"
)

type EvUpdate struct {
	hasUpdate     bool
	latestVersion string
	enabled       bool
	evVersion     string
	mutex         sync.RWMutex
	log           *zap.Logger
	evBackDao     *dao.EvBackDao
	noticeService *notice_service.NoticeService
	downloadUrl   string
}

func ProvideEvUpdate(log *logger.AppLogger, cfg *config.Config,
	evBackDao *dao.EvBackDao, noticeService *notice_service.NoticeService) (*EvUpdate, error) {
	log = log.Named("ev check update")
	return &EvUpdate{
		evBackDao:     evBackDao,
		enabled:       cfg.CheckForevUpdates,
		evVersion:     config.GetVersion(),
		log:           logger.ZapLog2AppLog(log),
		noticeService: noticeService,
	}, nil
}

func (s *EvUpdate) IsDisabled() bool {
	return !s.enabled
}

func (s *EvUpdate) Run(ctx context.Context) error {
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

func (s *EvUpdate) instrumentedCheckForUpdates(ctx context.Context) {
	start := time.Now()

	if err := s.checkForUpdates(ctx); err != nil {
		s.log.Sugar().Error("Update check failed", zap.Error(err), zap.String("所花时间", time.Since(start).String()))
		return
	}
	s.log.Sugar().Infof("Update check succeeded", "duration", time.Since(start).String())
}

func (s *EvUpdate) checkForUpdates(ctx context.Context) error {
	s.log.Debug("Checking for updates")

	versionRes, err := s.evBackDao.GetEvMaxVersion(ctx)

	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.latestVersion = versionRes.Version
	s.downloadUrl = versionRes.DownloadUrl

	s.hasUpdate = canUpdate(s.evVersion, s.latestVersion)

	go func() {
		if s.hasUpdate {
			s.noticeService.LiveBroadcastEvMsg2All(&dto.NoticeData{
				Title:       "Ev更新通知",
				Content:     fmt.Sprintf("您的ElasticView落后于官网最新版本(%s),请升级", s.latestVersion),
				Type:        "ElasticView更新通知",
				Level:       dto.NoticeLevelSuccess,
				IsTask:      true,
				FromUid:     0,
				PluginAlias: "",
				Source:      "ElasticView",
				NoticeJumpBtn: &dto.NoticeJumpBtn{
					Text:     "下载",
					JumpUrl:  s.downloadUrl,
					JumpType: dto.NoticeBtnJumpTypeRemote,
				},
				PublishTime: time.Now(),
			})
		}
	}()

	return nil
}

func (s *EvUpdate) UpdateAvailable() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.hasUpdate
}

func (s *EvUpdate) LatestVersion() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.latestVersion
}

func (s *EvUpdate) DownloadUrl() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.downloadUrl
}

func canUpdate(v1, v2 string) bool {
	ver1, err1 := version.NewVersion(v1)
	if err1 != nil {
		return false
	}
	ver2, err2 := version.NewVersion(v2)
	if err2 != nil {
		return false
	}

	return ver1.LessThan(ver2)
}
