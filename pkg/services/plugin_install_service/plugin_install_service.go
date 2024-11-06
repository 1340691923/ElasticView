package plugin_install_service

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/vo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/plugin"
	"github.com/1340691923/ElasticView/pkg/infrastructure/pluginstore"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type PluginInstaller struct {
	installing         sync.Map
	log                *logger.AppLogger
	cfg                *config.Config
	pluginStore        manager.Service
	evBackDao          *dao.EvBackDao
	pluginStoreService *pluginstore.PluginStoreService
	pluginsService     *updatechecker.PluginsService
}

func ProvideInstaller(cfg *config.Config, log *logger.AppLogger, pluginStore manager.Service,
	evBackDao *dao.EvBackDao, pluginStoreService *pluginstore.PluginStoreService, pluginsService *updatechecker.PluginsService) *PluginInstaller {
	return New(log, cfg, pluginStore, evBackDao, pluginStoreService, pluginsService)
}

func New(log *logger.AppLogger, cfg *config.Config, pluginStore manager.Service,
	evBackDao *dao.EvBackDao, pluginStoreService *pluginstore.PluginStoreService,
	pluginsService *updatechecker.PluginsService) *PluginInstaller {
	return &PluginInstaller{
		installing:         sync.Map{},
		log:                log.Named("plugin.installer"),
		cfg:                cfg,
		pluginStore:        pluginStore,
		evBackDao:          evBackDao,
		pluginStoreService: pluginStoreService,
		pluginsService:     pluginsService,
	}
}

func (m *PluginInstaller) Add(ctx context.Context, pluginID, version string) error {
	if ok, _ := m.installing.Load(pluginID); ok != nil {
		return nil
	}
	m.installing.Store(pluginID, true)
	defer func() {
		m.installing.Delete(pluginID)
	}()

	err := m.install(ctx, pluginID, version)
	if err != nil {
		return err
	}
	m.pluginsService.InstrumentedCheckForUpdates(ctx)

	return nil
}

func (m *PluginInstaller) install(ctx context.Context, pluginID, version string) (err error) {

	var pluginArchiveInfo *vo.GetPluginDownloadUrlRes

	if plugin, exists := m.plugin(ctx, pluginID, version); exists {

		if plugin.Version() == version {
			return errors.New("已安装该插件版本")
		}

		pluginArchiveInfo, err = m.evBackDao.GetPluginDownloadUrl(ctx, &dto.GetPluginDownloadUrlReq{
			PluginAlias: pluginID,
			Version:     version,
			Os:          runtime.GOOS,
			Arch:        runtime.GOARCH,
		})
		if err != nil {
			return err
		}

		m.log.Sugar().Infof("开始删除之前安装的老版本插件")
		err = m.Remove(ctx, plugin.ID, plugin.Version())
		if err != nil {
			return err
		}

	} else {
		var err error
		pluginArchiveInfo, err = m.evBackDao.GetPluginDownloadUrl(ctx, &dto.GetPluginDownloadUrlReq{
			PluginAlias: pluginID,
			Version:     version,
			Os:          runtime.GOOS,
			Arch:        runtime.GOARCH,
		})
		if err != nil {
			return errors.WithStack(err)
		}

	}
	m.log.Sugar().Infof("Installing plugin", "pluginId", pluginID, "version", version, pluginArchiveInfo.DownloadUrl)

	downloadPluginName, err := util.DownloadFile(pluginArchiveInfo.DownloadUrl, m.cfg.Plugin.LoadPath)

	if err != nil {
		return errors.WithStack(err)
	}

	/*crc32, err := util.FileCRC32(filepath.Join(m.cfg.Plugin.LoadPath, downloadPluginName))

	if err != nil {
		return errors.WithStack(err)
	}

	if cast.ToString(crc32) != pluginArchiveInfo.DonwloadCrc {
		m.log.Sugar().Errorf("crc文件完整性校验不通过", crc32, pluginArchiveInfo.DonwloadCrc)
		os.Remove(filepath.Join(m.cfg.Plugin.LoadPath, downloadPluginName))
		return errors.New("crc文件完整性校验不通过")
	}*/

	err = m.pluginStoreService.FastInitPlugin(ctx, downloadPluginName)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// 卸载插件
func (m *PluginInstaller) Remove(ctx context.Context, pluginID, version string) error {
	plugin, exists := m.plugin(ctx, pluginID, version)
	if !exists {
		return errors.New("插件不存在")
	}
	pluginFilePath := plugin.GetPluginFileName()
	err := m.pluginStoreService.FastShutdown(ctx, plugin)

	if err != nil {
		return errors.WithStack(err)
	}

	pluginFilePath = filepath.Join(m.cfg.Plugin.LoadPath, pluginFilePath)

	os.Remove(pluginFilePath)

	return err
}

func (m *PluginInstaller) plugin(ctx context.Context, pluginID, pluginVersion string) (*plugin.Plugin, bool) {
	p, exists := m.pluginStore.Plugin(ctx, pluginID)
	if !exists {
		return nil, false
	}

	return p, true
}
