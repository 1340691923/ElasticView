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
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/backendplugin/provider"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/manager/process"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore"
	util2 "github.com/1340691923/ElasticView/pkg/util"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type PluginService struct {
	orm            *sqlstore.SqlStore
	pluginRegistry manager.Service
	log            *logger.AppLogger
	progressSvr    *process.Service
	cfg            *config.Config
	rbac           *access_control.Rbac
	gmUserDao      *dao.GmUserDao
	jwtSvr         *jwt_svr.Jwt
	evBackDao      *dao.EvBackDao
}

func NewPluginService(orm *sqlstore.SqlStore, pluginRegistry manager.Service, log *logger.AppLogger, progressSvr *process.Service, cfg *config.Config, rbac *access_control.Rbac, gmUserDao *dao.GmUserDao, jwtSvr *jwt_svr.Jwt, evBackDao *dao.EvBackDao) *PluginService {
	return &PluginService{orm: orm, pluginRegistry: pluginRegistry, log: log, progressSvr: progressSvr, cfg: cfg, rbac: rbac, gmUserDao: gmUserDao, jwtSvr: jwtSvr, evBackDao: evBackDao}
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
		return
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
	c, err := this.jwtSvr.ParseToken(ctx.GetHeader("X-Token"))
	if err != nil {
		return errors.WithStack(err)
	}

	pluginJsonData := plugin.PluginData().PluginJsonData

	roles, err := this.gmUserDao.GetRolesFromUser(c.UserID)
	if err != nil {
		return errors.WithStack(err)
	}

	contentType := ctx.GetHeader("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var b []byte
		b, _ = ctx.GetRawData()

		bodySize := len(b)

		// 判断大小，例如设置限制为1MB
		const maxBodySize = 20 * 1024 * 1024 // 20MB

		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		gmOperaterLog := &model.GmOperaterLog{
			OperaterName:   cast.ToString(c.Username),
			OperaterId:     cast.ToInt(c.UserID),
			OperaterAction: fmt.Sprintf("/%s%s", pluginID, ctx.Param("action")),
			Method:         ctx.Request.Method,
			Created:        time.Now(),
		}
		if bodySize <= maxBodySize {
			var body []byte
			body, err = util2.GzipCompress(util2.Bytes2str(b))
			if err != nil {
				return
			}

			gmOperaterLog.Body = body
		}

		err = this.orm.Create([]*model.GmOperaterLog{gmOperaterLog}).Error

		if err != nil {
			this.log.Error("OperaterLog", zap.Error(err))
		}
	}

	if this.IsAdminUser(roles) && !pluginJsonData.BackendDebug {
		for _, routerConfig := range pluginJsonData.BackendRoutes {

			if ctx.Param("action") == routerConfig.Path && routerConfig.NeedAuth {

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
	ctx.Request.Header.Set(util.EvUserID, cast.ToString(c.UserID))
	delete(ctx.Request.Header, "X-Token")
	plugins.NewDataSourcePlugin(ctx, plugin).CallPluginResource()

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
	log, closeLogCallback, err := logger.InitPluginLog(this.cfg, pluginID)
	if err != nil {
		return errors.WithStack(err)
	}
	p = provider.DefaultProvider(ctx, log, closeLogCallback, &provider.Config{
		ID:       pluginID,
		IsDebug:  true,
		TestAddr: addr,
		TestPid:  pid,
		ExecArgs: []string{"-evRpcKey=123"},
	}, this.cfg)

	err = p.Start(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	err = this.pluginRegistry.Add(ctx, p)
	if err != nil {
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

const AdminRole = 1

func (this *PluginService) IsAdminUser(roleId []int) bool {
	return util2.InArr(roleId, AdminRole)
}
