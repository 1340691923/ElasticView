package plugin_service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/jwt_svr"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/model"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/backendplugin/provider"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager/process"
	"github.com/1340691923/ElasticView/pkg/infrastructure/pluginstore"
	"github.com/1340691923/ElasticView/pkg/services/gm_operater_log"
	"github.com/1340691923/ElasticView/pkg/services/updatechecker"
	util2 "github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"time"
)

type PluginService struct {
	orm                 *orm.Gorm
	pluginRegistry      manager.Service
	log                 *logger.AppLogger
	progressSvr         *process.Service
	cfg                 *config.Config
	rbac                *access_control.Rbac
	gmUserDao           *dao.GmUserDao
	jwtSvr              *jwt_svr.Jwt
	evBackDao           *dao.EvBackDao
	pluginStoreService  *pluginstore.PluginStoreService
	pluginUpdateChecker *updatechecker.PluginsService
	gmOperaterLogSvr    *gm_operater_log.GmOperaterLogService
}

func NewPluginService(orm *orm.Gorm, pluginRegistry manager.Service, log *logger.AppLogger, progressSvr *process.Service, cfg *config.Config, rbac *access_control.Rbac, gmUserDao *dao.GmUserDao, jwtSvr *jwt_svr.Jwt, evBackDao *dao.EvBackDao, pluginStoreService *pluginstore.PluginStoreService, pluginUpdateChecker *updatechecker.PluginsService, gmOperaterLogSvr *gm_operater_log.GmOperaterLogService) *PluginService {
	return &PluginService{orm: orm, pluginRegistry: pluginRegistry, log: log, progressSvr: progressSvr, cfg: cfg, rbac: rbac, gmUserDao: gmUserDao, jwtSvr: jwtSvr, evBackDao: evBackDao, pluginStoreService: pluginStoreService, pluginUpdateChecker: pluginUpdateChecker, gmOperaterLogSvr: gmOperaterLogSvr}
}

func (this *PluginService) ExecSql(ctx context.Context, pluginID string, sql string, args []interface{}) (rowsAffected int64, err error) {
	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return 0, errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}
	p.DbLock()
	defer p.DbUnlock()

	result := p.Gorm().WithContext(ctx).Exec(sql, args...)
	if result.Error != nil {
		err = result.Error

		if strings.Contains(err.Error(), "busy") {
			err = util2.Retry(3, 5*time.Second, func() error {
				result := p.Gorm().WithContext(ctx).Exec(sql, args...)
				if result.Error != nil {
					return result.Error
				}
				rowsAffected = result.RowsAffected
				return nil
			})
			return
		}
	}

	rowsAffected = result.RowsAffected
	return
}

// todo... 表名鉴权
func (this *PluginService) SelectSql(ctx context.Context, pluginID string, sql string, args []interface{}) (list []map[string]interface{}, err error) {

	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return nil, errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}

	err = p.Gorm().WithContext(ctx).Raw(sql, args...).Scan(&list).Error

	if err != nil {
		return nil, err
	}

	return
}

func (this *PluginService) CallPlugin(ctx *gin.Context, pluginID string) (err error) {
	plugin, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}
	c, hasTokenErr := this.jwtSvr.ParseToken(ctx.GetHeader("X-Token"))

	pluginJsonData := plugin.PluginData().PluginJsonData

	contentType := ctx.GetHeader("Content-Type")
	var body []byte
	startT := time.Now()
	if strings.Contains(contentType, "application/json") {

		if hasTokenErr == nil {
			var b []byte
			b, _ = ctx.GetRawData()
			bodySize := len(b)
			// 判断大小，例如设置限制为1MB
			const maxBodySize = 20 * 1024 * 1024 // 20MB
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			if bodySize <= maxBodySize {
				body, err = util2.GzipCompress(util2.Bytes2str(b))
				if err != nil {
					return
				}

			}
		}
	}

	if !pluginJsonData.BackendDebug {
		var roles []int
		if hasTokenErr == nil {
			roles, err = this.gmUserDao.GetRolesFromUser(c.UserID)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		for _, routerConfig := range pluginJsonData.BackendRoutes {
			if this.IsAdminUser(roles) {
				break
			}

			if ctx.Param("action") == routerConfig.Path && routerConfig.NeedAuth {
				if len(roles) == 0 {
					return errors.WithStack(errors.New("请检查该用户的角色分配"))
				}
				for _, roleId := range roles {
					ok, err := this.rbac.Enforce(strconv.Itoa(roleId),
						fmt.Sprintf("/%s%s", pluginID, routerConfig.Path), "*")
					if err != nil {
						return errors.WithStack(err)
					}
					if !ok {
						return errors.New(fmt.Sprintf("您没有操作该资源的权限:%s[%s]", pluginJsonData.PluginName, routerConfig.Remark))
					}
					break
				}
			}
		}

	}

	if hasTokenErr == nil {
		ctx.Request.Header.Set(util.EvUserID, cast.ToString(c.UserID))
		delete(ctx.Request.Header, "X-Token")
	}

	plugins.NewDataSourcePlugin(ctx, plugin).CallPluginResource()

	if strings.Contains(contentType, "application/json") {

		if hasTokenErr == nil {
			costT := time.Now().Sub(startT).String()
			resStatus := "失败"
			if ctx.Writer.Status() == 200 {
				resStatus = "正常"
			}

			gmOperaterLog := &model.GmOperaterLog{
				OperaterName:   cast.ToString(c.Username),
				OperaterId:     cast.ToInt(c.UserID),
				OperaterAction: fmt.Sprintf("/%s%s", pluginID, ctx.Param("action")),
				Method:         ctx.Request.Method,
				Created:        time.Now(),
				Status:         resStatus,
				CostTime:       costT,
				Body:           body,
			}

			this.gmOperaterLogSvr.Save(gmOperaterLog)

		}

	}

	return nil
}

func (this *PluginService) CallPluginViews(ctx *gin.Context, pluginID string) (err error) {
	plugin, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}
	plugins.NewDataSourcePlugin(ctx, plugin).CallPluginResource()
	return nil
}

func (this *PluginService) LoadDebugPlugin(ctx context.Context, pluginID string, addr string, pid int) (err error) {

	p, ok := this.pluginRegistry.Plugin(ctx, pluginID)
	if ok {
		err = this.progressSvr.Stop(ctx, p)
		if err != nil {
			return errors.WithStack(err)
		}
		err = this.pluginRegistry.Remove(ctx, p.ID)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	log, logPath, closeLogCallback, err := logger.InitPluginLog(this.cfg, pluginID)
	if err != nil {
		return errors.WithStack(err)
	}
	p = provider.DefaultProvider(ctx, log, logPath, closeLogCallback, &provider.Config{
		ID:             pluginID,
		IsDebug:        true,
		TestAddr:       addr,
		TestPid:        pid,
		PluginFileName: "调试插件",
		ExecArgs:       []string{},
	}, this.cfg)

	err = p.Start(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = this.pluginRegistry.Add(ctx, p)
	if err != nil {
		this.log.Sugar().Errorf("err", err)
		return errors.WithStack(err)
	}
	return nil
}

func (this *PluginService) StopDebugPlugin(ctx context.Context, pluginID string) (err error) {

	p, ok := this.pluginRegistry.Plugin(ctx, pluginID)
	if ok {
		err = this.progressSvr.Stop(ctx, p)
		if err != nil {
			return errors.WithStack(err)
		}
		err = this.pluginRegistry.Remove(ctx, p.ID)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (this *PluginService) StarPlugin(ctx context.Context, pluginId int64) (err error) {
	err = this.evBackDao.StarPlugin(ctx, &dto.StarPlugin{PluginId: pluginId})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type PluginVo struct {
	Developer        string `json:"developer"`
	PluginId         string `json:"plugin_id"`
	PluginFileName   string `json:"plugin_file_name"`
	PluginName       string `json:"plugin_name"`
	PluginAlias      string `json:"plugin_alias"`
	Version          string `json:"version"`
	Frontend2c       bool   `json:"frontend_2c"`
	BackendDebug     bool   `json:"backend_debug"`
	FrontendDebug    bool   `json:"frontend_debug"`
	FrontendDevPort  int    `json:"frontend_dev_port"`
	Pid              int    `json:"pid"`
	StorePath        string `json:"store_path"`
	CpuPercentStr    string `json:"cpu_percent_str"`
	MemoryPercentStr string `json:"memory_percent_str"`
	StartTime        string `json:"start_time"`
	StopTime         string `json:"stop_time"`
	IsExited         bool   `json:"is_exited"`
	LogFilePath      string `json:"log_file_path"`
	HasUpdate        bool   `json:"has_update"`
	UpdateVersion    string `json:"update_version"`
}

func (this *PluginService) PluginList(ctx context.Context) (res []*PluginVo) {
	ps := this.pluginRegistry.Plugins(ctx)
	for _, v := range ps {
		pluginID := v.PluginID()
		pluginFileName := v.GetPluginFileName()
		client, hasClient := v.Client()
		if !hasClient {
			this.log.Sugar().Errorf("get pluginClient err", v.PluginID())
			continue
		}

		developer := v.PluginData().PluginJsonData.Developer
		pluginData := v.PluginData().PluginJsonData
		pluginName := pluginData.PluginName
		pluginAlias := pluginData.PluginAlias
		version := pluginData.Version
		frontend2c := pluginData.Frontend2c
		logFilePath := v.LogFilePath
		backendDebug := pluginData.BackendDebug
		frontendDebug := pluginData.FrontendDebug
		frontendDevPort := pluginData.FrontendDevPort
		storePath := v.GetStorePath()
		pid := client.GetPid()
		cpuPercentStr := ""
		memoryPercentStr := ""
		procUtil, err := client.GetProcessUtil()
		if err == nil {
			cpuPercent, err := procUtil.CPUPercent()
			if err != nil {
				this.log.Error("err", zap.Error(err))
				continue
			}
			cpuPercentStr = fmt.Sprintf("%.2f", cpuPercent) + "%"
			memoryPercent, err := procUtil.MemoryPercent()
			if err != nil {
				this.log.Error("err", zap.Error(err))
				continue
			}
			memoryPercentStr = fmt.Sprintf("%.2f", memoryPercent) + "%"
		}
		startTime := v.StartTime.Format(time.DateTime)
		stopTime := v.StopTime.Format(time.DateTime)
		isExited := v.Exited()
		updateVersion, hasUpdate := this.pluginUpdateChecker.HasUpdate(ctx, v.PluginID())
		res = append(res, &PluginVo{
			UpdateVersion:    updateVersion,
			HasUpdate:        hasUpdate,
			PluginId:         pluginID,
			PluginFileName:   pluginFileName,
			PluginName:       pluginName,
			PluginAlias:      pluginAlias,
			Version:          version,
			Frontend2c:       frontend2c,
			BackendDebug:     backendDebug,
			FrontendDebug:    frontendDebug,
			FrontendDevPort:  frontendDevPort,
			Pid:              pid,
			StorePath:        storePath,
			LogFilePath:      logFilePath,
			CpuPercentStr:    cpuPercentStr,
			MemoryPercentStr: memoryPercentStr,
			StartTime:        startTime,
			StopTime:         stopTime,
			IsExited:         isExited,
			Developer:        developer,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Pid > res[j].Pid
	})

	return
}

const AdminRole = 1

func (this *PluginService) IsAdminUser(roleId []int) bool {
	return util2.InArr(roleId, AdminRole)
}
