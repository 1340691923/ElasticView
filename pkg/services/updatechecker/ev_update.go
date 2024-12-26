package updatechecker

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/hashicorp/go-version"
	"go.uber.org/zap"
	"sync"
	"time"
)

type EvUpdate struct {
	hasUpdate     bool
	latestVersion string

	enabled     bool
	evVersion   string
	mutex       sync.RWMutex
	log         *zap.Logger
	evBackDao   *dao.EvBackDao
	downloadUrl string
}

func ProvideEvUpdate(log *logger.AppLogger, cfg *config.Config, evBackDao *dao.EvBackDao) (*EvUpdate, error) {
	log = log.Named("ev check update")
	return &EvUpdate{
		evBackDao: evBackDao,
		enabled:   cfg.CheckForevUpdates,
		evVersion: config.GetVersion(),
		log:       logger.ZapLog2AppLog(log),
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
		s.log.Sugar().Errorf("Update check failed", "error", err, "duration", time.Since(start))
		return
	}
	s.log.Sugar().Infof("Update check succeeded", "duration", time.Since(start))
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
