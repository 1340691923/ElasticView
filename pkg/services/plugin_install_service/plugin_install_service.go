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
	evBackDao *dao.EvBackDao, pluginStoreService *pluginstore.PluginStoreService,
	pluginsService *updatechecker.PluginsService) *PluginInstaller {
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

func (this *PluginInstaller) Add(ctx context.Context, pluginID, version string) error {
	if ok, _ := this.installing.Load(pluginID); ok != nil {
		return nil
	}
	this.installing.Store(pluginID, true)
	defer func() {
		this.installing.Delete(pluginID)
	}()

	err := this.install(ctx, pluginID, version)
	if err != nil {
		return err
	}
	this.pluginsService.InstrumentedCheckForUpdates(ctx)

	return nil
}

func (this *PluginInstaller) AddUploadPlugin(ctx *gin.Context, f *multipart.FileHeader) (pluginName string, err error) {
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

	if ok, _ := this.installing.Load(pluginID); ok != nil {
		return "", errors.New(fmt.Sprintf("插件[%s]已在安装中...", pluginID))
	}
	this.installing.Store(pluginID, true)
	defer func() {
		this.installing.Delete(pluginID)
	}()

	if plugin, exists := this.plugin(ctx, pluginID, ""); exists {

		this.log.Sugar().Infof("开始删除之前安装的老版本插件")
		err = this.Remove(ctx, plugin.ID, plugin.Version())
		if err != nil {
			return "", errors.WithStack(err)
		}
	}
	this.log.Sugar().Infof("Installing upload plugin", "pluginId", pluginID)
	dest := filepath.Join(this.cfg.Plugin.LoadPath, fileName)
	err = ctx.SaveUploadedFile(f, dest)
	if err != nil {
		os.Remove(dest)
		return "", errors.WithStack(err)
	}
	err = this.pluginStoreService.FastInitPlugin(ctx, fileName)

	if err != nil {
		os.Remove(dest)
		return "", err
	}

	this.pluginsService.InstrumentedCheckForUpdates(ctx)

	return pluginID, nil
}

func (this *PluginInstaller) install(ctx context.Context, pluginID, version string) (err error) {
	isLinux := runtime.GOOS == "linux" //linux 才进行热更
	var pluginArchiveInfo *vo.GetPluginDownloadUrlRes
	var exists bool
	var plugin *plugin.Plugin
	var pluginPath string
	var oldPluginPath string
	if plugin, exists = this.plugin(ctx, pluginID, version); exists {

		if plugin.Version() == version {
			return errors.New("已安装该插件版本")
		}

		//获取插件远端下载地址
		pluginArchiveInfo, err = this.evBackDao.GetPluginDownloadUrl(ctx, &dto.GetPluginDownloadUrlReq{
			PluginAlias: pluginID,
			Version:     version,
			Os:          runtime.GOOS,
			Arch:        runtime.GOARCH,
		})
		if err != nil {
			return err
		}
		if isLinux {
			this.log.Sugar().Infof("开始改名之前的插件名为插件-old")

			pluginPath, oldPluginPath, err = this.Rename(ctx, plugin.ID)
			if err != nil {
				this.log.Sugar().Errorf("插件改名出现错误", oldPluginPath, pluginPath, err)
				os.Rename(oldPluginPath, pluginPath)
				return err
			}
		} else {
			this.log.Sugar().Infof("开始删除之前安装的老版本插件")

			err = this.Remove(ctx, plugin.ID, plugin.Version())
			if err != nil {
				return err
			}
		}

	} else {
		var err error
		pluginArchiveInfo, err = this.evBackDao.GetPluginDownloadUrl(ctx, &dto.GetPluginDownloadUrlReq{
			PluginAlias: pluginID,
			Version:     version,
			Os:          runtime.GOOS,
			Arch:        runtime.GOARCH,
		})
		if err != nil {
			return errors.WithStack(err)
		}

	}
	this.log.Sugar().Infof("开始下载插件", "pluginId", pluginID, "version", version, pluginArchiveInfo.DownloadUrl)

	downloadPluginName, err := util.DownloadFile(pluginArchiveInfo.DownloadUrl, this.cfg.Plugin.LoadPath)

	if err != nil {
		this.log.Sugar().Errorf("下载插件出现错误", oldPluginPath, pluginPath, err)
		os.Rename(oldPluginPath, pluginPath)
		return errors.WithStack(err)
	}

	if isLinux && exists {
		this.log.Sugar().Infof("开始插件热更新")
		err = this.pluginStoreService.FastReloadPlugin(ctx, downloadPluginName)
		if err != nil {
			this.log.Sugar().Errorf("插件热更新出现错误", oldPluginPath, pluginPath, err)
			os.Rename(oldPluginPath, pluginPath)
			return errors.WithStack(err)
		} else {
			this.log.Sugar().Infof("开始刪除旧插件文件", oldPluginPath)
			os.Remove(oldPluginPath)
		}
	} else {
		err = this.pluginStoreService.FastInitPlugin(ctx, downloadPluginName)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	this.pluginsService.InstrumentedCheckForUpdates(ctx)

	return nil
}

// 卸载插件
func (this *PluginInstaller) Remove(ctx context.Context, pluginID, version string) error {
	plugin, exists := this.plugin(ctx, pluginID, version)
	if !exists {
		return errors.New("插件不存在")
	}
	pluginFilePath := plugin.GetPluginFileName()
	err := this.pluginStoreService.FastShutdown(ctx, plugin)

	if err != nil {
		return errors.WithStack(err)
	}

	err = this.pluginStoreService.FastRemove(ctx, plugin)

	if err != nil {
		return errors.WithStack(err)
	}

	pluginFilePath = filepath.Join(this.cfg.Plugin.LoadPath, pluginFilePath)

	os.Remove(pluginFilePath)

	return err
}

// 重命名插件
func (this *PluginInstaller) Rename(ctx context.Context, pluginID string) (oldPluginFilePath, newPluginFilePath string, err error) {
	plugin, exists := this.plugin(ctx, pluginID, "")
	if !exists {
		return "", "", errors.New("插件不存在")
	}
	pluginFileName := plugin.GetPluginFileName()

	pluginFilePath := filepath.Join(this.cfg.Plugin.LoadPath, pluginFileName)

	newPluginFilePath = filepath.Join(this.cfg.Plugin.LoadPath, fmt.Sprintf("%s%s", pluginFileName, "-old"))
	this.log.Sugar().Infof("修改插件名", pluginFilePath, newPluginFilePath)
	return pluginFilePath, newPluginFilePath, os.Rename(pluginFilePath, newPluginFilePath)
}

func (this *PluginInstaller) plugin(ctx context.Context, pluginID, pluginVersion string) (*plugin.Plugin, bool) {
	p, exists := this.pluginStore.Plugin(ctx, pluginID)
	if !exists {
		return nil, false
	}

	return p, true
}
