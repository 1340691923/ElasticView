package plugin_install_service

import (
	"context"
	"fmt"
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
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func (m *PluginInstaller) AddUploadPlugin(ctx *gin.Context, f *multipart.FileHeader) (pluginName string, err error) {
	fileName := f.Filename
	buildOs := runtime.GOOS
	buildArch := runtime.GOARCH

	if !strings.Contains(fileName, buildOs) || !strings.Contains(fileName, buildArch) {
		return "", errors.New("该插件不符合当前操作系统架构")
	}

	tmpArr := strings.Split(fileName, "_")

	if len(tmpArr) == 0 {
		return "", errors.New("该插件不符合ElasticView插件命名规范")
	}
	pluginID := tmpArr[0]

	if ok, _ := m.installing.Load(pluginID); ok != nil {
		return "", errors.New(fmt.Sprintf("插件[%s]已在安装中...", pluginID))
	}
	m.installing.Store(pluginID, true)
	defer func() {
		m.installing.Delete(pluginID)
	}()

	if plugin, exists := m.plugin(ctx, pluginID, ""); exists {

		m.log.Sugar().Infof("开始删除之前安装的老版本插件")
		err = m.Remove(ctx, plugin.ID, plugin.Version())
		if err != nil {
			return "", errors.WithStack(err)
		}
	}
	m.log.Sugar().Infof("Installing upload plugin", "pluginId", pluginID)
	dest := filepath.Join(m.cfg.Plugin.LoadPath, fileName)
	err = ctx.SaveUploadedFile(f, dest)
	if err != nil {
		os.Remove(dest)
		return "", errors.WithStack(err)
	}
	err = m.pluginStoreService.FastInitPlugin(ctx, fileName)

	if err != nil {
		os.Remove(dest)
		return "", err
	}

	m.pluginsService.InstrumentedCheckForUpdates(ctx)

	return pluginID, nil
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

	err = m.pluginStoreService.FastRemove(ctx, plugin)

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
