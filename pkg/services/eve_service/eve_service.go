package eve_service

import (
	"context"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/vo"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	"github.com/pkg/errors"
	"time"
)

type EvEService struct {
	log            *logger.AppLogger
	evBackDao      *dao.EvBackDao
	cfg            *config.Config
	pluginStore    manager.Service
	pluginsService *updatechecker.PluginsService
}

func NewEvEService(log *logger.AppLogger, evBackDao *dao.EvBackDao, cfg *config.Config, pluginStore manager.Service, pluginsService *updatechecker.PluginsService) *EvEService {
	return &EvEService{log: log, evBackDao: evBackDao, cfg: cfg, pluginStore: pluginStore, pluginsService: pluginsService}
}

func (this *EvEService) FlushAccessToken(ctx context.Context) (err error) {

	evKey := this.GetEvKey()
	if evKey == "" {
		return nil
	}
	_, err = this.evBackDao.GetEvAccessToken(ctx, evKey)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *EvEService) GetRemotePlugins(ctx context.Context, req *dto.FromEvPluginReq) (*vo.PluginListRes, error) {

	ps := this.pluginStore.Plugins(ctx)
	installPlugins := map[string]struct{}{}
	installPluginArr := []string{}
	for _, v := range ps {
		installPlugins[v.PluginID()] = struct{}{}
		installPluginArr = append(installPluginArr, v.PluginID())
	}

	if req.HasDownloadType != nil {
		req.HasDownloadPlugins = installPluginArr
	}

	pluginList, err := this.evBackDao.GetPluginList(ctx, req)
	if err != nil {
		return pluginList, errors.WithStack(err)
	}

	for index := range pluginList.List {
		if _, ok := installPlugins[pluginList.List[index].PluginAlias]; ok {
			pluginList.List[index].HasDownload = true
		}
		if pluginList.List[index].StarState == nil {
			unStar := 2
			pluginList.List[index].StarState = &unStar
		}
	}

	return pluginList, nil
}

func (this *EvEService) GetRemotePluginInfo(ctx context.Context, req *dto.FormEvPluginInfoReq) (*vo.PublishRes, error) {
	res, err := this.evBackDao.GetPluginInfo(ctx, req)
	if err != nil {
		return res, errors.WithStack(err)
	}
	return res, nil
}

func (this *EvEService) GetWxArticleList(ctx context.Context) (*vo.WxArticleList, error) {
	return this.evBackDao.GetWxArticleList(ctx)
}

func (this *EvEService) GetLocalPlugins(ctx context.Context) ([]vo.GetLocalPlugin, error) {
	ps := this.pluginStore.Plugins(ctx)
	res := []vo.GetLocalPlugin{}
	for _, v := range ps {
		updateVersion, hasUpdate := this.pluginsService.HasUpdate(ctx, v.PluginID())
		res = append(res, vo.GetLocalPlugin{
			PluginID:      v.PluginData().PluginJsonData.PluginAlias,
			PluginName:    v.PluginData().PluginJsonData.PluginName,
			Version:       v.Version(),
			HasUpdate:     hasUpdate,
			UpdateVersion: updateVersion,
		})
	}
	return res, nil
}

func (this *EvEService) Run(ctx context.Context) error {
	var err error
	ticker := time.NewTicker(time.Minute * 10)
	run := true

	for run {
		select {
		case <-ticker.C:
			err = this.FlushAccessToken(ctx)
			if err != nil {
				return errors.WithStack(err)
			}
		case <-ctx.Done():
			run = false
		}
	}

	return ctx.Err()
}

func (this *EvEService) GetEvKey() string {
	return this.cfg.EvKey
}

func (this *EvEService) SaveEvKey(evKey string) error {
	return this.cfg.SetEvKey(evKey).GetViperInstance().WriteConfig()
}
