package plugin_service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/access_control"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/dao"
	dto2 "github.com/1340691923/ElasticView/pkg/infrastructure/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/dto"
	"github.com/1340691923/ElasticView/pkg/infrastructure/eve_api/vo"
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
	"github.com/1340691923/eve-plugin-sdk-go/enum"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (this *PluginService) ExecMoreSql(ctx context.Context, pluginID string, sqls []dto2.ExecSql) (err error) {
	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		err = errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
		return
	}

	tx := p.Gorm().Begin()

	for _, v := range sqls {
		result := tx.WithContext(ctx).Exec(v.Sql, v.Args...)
		if result.Error != nil {
			err = result.Error
			tx.Rollback()
			return
		}
	}
	tx.Commit()

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

func (this *PluginService) SaveDb(ctx context.Context, pluginID, table string, data map[string]interface{}) (err error) {

	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}

	if _, ok := data["id"]; ok && cast.ToInt(data["id"]) > 0 {
		err = this.InsertOrUpdateDb(ctx, pluginID, table, data, []string{"id"})
		return
	}

	err = p.Gorm().Table(table).WithContext(ctx).Create(data).Error

	if err != nil {
		return err
	}

	return
}

func (this *PluginService) UpdateDb(ctx context.Context, pluginID, table string, updateSql string, updateArgs []interface{}, data map[string]interface{}) (rowsAffected int64, err error) {

	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return 0, errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}
	query := p.Gorm().Table(table).Where(updateSql, updateArgs...).WithContext(ctx).Updates(data)
	err = query.Error

	if err != nil {
		return 0, err
	}

	rowsAffected = query.RowsAffected

	return rowsAffected, nil
}

func (this *PluginService) DeleteDb(ctx context.Context, pluginID,
	table string, whereSql string, whereArgs []interface{}) (rowsAffected int64, err error) {

	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return 0, errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}

	query := p.Gorm().Table(table).Where(whereSql, whereArgs...).WithContext(ctx).Delete(nil)

	err = query.Error
	if err != nil {
		return 0, err
	}

	rowsAffected = query.RowsAffected

	return rowsAffected, nil
}

func (this *PluginService) InsertOrUpdateDb(ctx context.Context, pluginID, table string, upsertData map[string]interface{}, uniqueKeys []string) (err error) {

	p, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}
	g := p.Gorm()

	query := g.Debug().Table(table)

	err = applyUpsert(ctx, query, QueryDSL{
		UpsertData: upsertData,
		UniqueKeys: uniqueKeys,
	}).Error

	if err != nil {
		return err
	}

	return
}

type QueryDSL struct {
	UpsertData map[string]interface{} // 没有则新增，有则更新
	UniqueKeys []string               // 冲突检查的唯一键
}

func applyUpsert(ctx context.Context, query *gorm.DB, dsl QueryDSL) *gorm.DB {
	if dsl.UpsertData != nil && len(dsl.UniqueKeys) > 0 {
		query = query.Clauses(clause.OnConflict{
			Columns:   convertToClauseColumns(dsl.UniqueKeys),
			DoUpdates: clause.Assignments(dsl.UpsertData),
		}).WithContext(ctx).Create(dsl.UpsertData)
	}
	return query
}

// 转换唯一键为 GORM 结构
func convertToClauseColumns(keys []string) []clause.Column {
	var cols []clause.Column
	for _, key := range keys {
		cols = append(cols, clause.Column{Name: key})
	}
	return cols
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
					return errors.WithStack(errors.New("请检查该用户的权限组分配"))
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
		ctx.Request.Header.Set(enum.EvUserID, cast.ToString(c.UserID))
		delete(ctx.Request.Header, "X-Token")
	}

	plugins.NewDataSourcePlugin(ctx, plugin, this.log).CallPluginResource()

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
	plugins.NewDataSourcePlugin(ctx, plugin, this.log).CallPluginResource()
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
		ExecArgs: []string{
			fmt.Sprintf("-dbType=%s", this.cfg.DbType),
		},
	}, this.cfg, this.orm)

	err = p.Start(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	err = this.pluginRegistry.AddPlugin(ctx, p)
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

func (this *PluginService) AddComment(ctx context.Context, pluginId int, content string, parentId int) (err error) {
	err = this.evBackDao.AddComment(ctx, &dto.AddCommentRequest{
		PluginID: pluginId,
		Content:  content,
		ParentID: parentId,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *PluginService) LikeComment(ctx context.Context, commentId int, state int) (err error) {
	err = this.evBackDao.LikeComment(ctx, &dto.LikeCommentRequest{
		CommentID: commentId,
		State:     state,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (this *PluginService) ListComments(ctx context.Context, pluginId int) (list *[]*vo.Comment, err error) {
	list, err = this.evBackDao.ListComments(ctx, &dto.ListCommentsRequest{
		PluginID: pluginId,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
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

		storePath := fmt.Sprintf("SQLITE3_DB:%s", v.GetStorePath())

		if this.cfg.DbType == config.MysqlDbTyp {
			storePath = fmt.Sprintf("MYSQL_DB:%s", v.GetMysqlDbPath())
		}
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

func (this *PluginService) CallPluginNoAuth(ctx *gin.Context, pluginID string) (err error) {
	plugin, b := this.pluginRegistry.Plugin(ctx, pluginID)
	if !b {
		return errors.New(fmt.Sprintf("没有找到该插件信息:%s", pluginID))
	}

	userId := util.GetEvUserID(ctx)

	pluginJsonData := plugin.PluginData().PluginJsonData

	if !pluginJsonData.BackendDebug {
		var roles []int
		if userId > 0 {
			roles, err = this.gmUserDao.GetRolesFromUser(userId)
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
					return errors.WithStack(errors.New("请检查该用户的权限组分配"))
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

	this.log.Sugar().Info("CallPluginNoAuth", zap.Int("userId", userId), zap.String("action", ctx.Param("action")))

	plugins.NewDataSourcePlugin(ctx, plugin, this.log).CallPluginResource()

	return nil
}

func (this *PluginService) GetPluginName(id string) (pluginName string) {
	p, has := this.pluginRegistry.Plugin(context.Background(), id)
	if !has {
		return ""
	}
	return p.PluginData().PluginJsonData.PluginName
}
