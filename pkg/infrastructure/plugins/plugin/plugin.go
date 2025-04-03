package plugin

import (
	"context"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"
	"github.com/1340691923/ElasticView/pkg/infrastructure/plugins/backendplugin"
	"github.com/1340691923/eve-plugin-sdk-go/backend"
	"github.com/1340691923/eve-plugin-sdk-go/build"
	"github.com/1340691923/eve-plugin-sdk-go/util"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Plugin struct {
	ID             string
	mu             sync.Mutex
	PluginDir      string
	PluginFileName string
	log            hclog.Logger
	client         backendplugin.Plugin
	orm            *orm.Gorm
	EvOrm          *orm.Gorm
	Cfg            *config.Config
	pluginData     *build.PluginInitRespData
	LogFilePath    string //日志文件存放地址
	StartTime      time.Time
	StopTime       time.Time
}

func (p *Plugin) GetPluginFileName() string {
	return p.PluginFileName
}

func (p *Plugin) PluginID() string {
	return p.ID
}

func (p *Plugin) PluginData() *build.PluginInitRespData {
	return p.pluginData
}

func (p *Plugin) Logger() hclog.Logger {
	return p.log
}

func (p *Plugin) Gorm() *orm.Gorm {
	return p.orm
}

func (p *Plugin) SetLogger(l hclog.Logger) {
	p.log = l
}

func (p *Plugin) getSqlLiteDbName() string {
	pluginID := p.ID
	if p.PluginData().PluginJsonData.BackendDebug {
		return fmt.Sprintf("%s-test", pluginID)
	}
	return pluginID
}

func (p *Plugin) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client == nil {
		return fmt.Errorf("could not start plugin %s as no plugin client exists", p.ID)
	}

	err := p.client.Start(ctx)

	if err != nil {
		return errors.WithStack(err)
	}

	checkResultRes, err := p.client.CheckHealth(ctx, &backend.CheckHealthRequest{})
	if err != nil {
		return errors.WithStack(err)
	}

	if checkResultRes == nil {
		return errors.New("插件没有返回心跳返回值")
	}

	p.pluginData = new(build.PluginInitRespData)
	err = json.Unmarshal(checkResultRes.JSONDetails, &p.pluginData)
	if err != nil {
		return errors.New("插件信息返回异常")
	}

	var pluginOrm *orm.Gorm

	pluginOrm, err = orm.NewPluginGorm(p.EvOrm, p.Cfg,
		p.GetMysqlDbPath(),
		p.GetStorePath(), p.Logger())

	if err != nil {
		return errors.WithStack(err)
	}

	p.orm = pluginOrm

	p.Migrator()

	p.StartTime = time.Now()

	return nil
}

func (p *Plugin) GetStorePath() string {
	return filepath.Join(p.Cfg.Plugin.StorePath, fmt.Sprintf("%s.db", p.getSqlLiteDbName()))
}

func (p *Plugin) GetMysqlDbPath() string {
	return strings.ToLower(strings.ReplaceAll(p.ID, "-", "_"))
}

func (p *Plugin) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.orm = nil

	if p.client == nil {
		return nil
	}

	err := p.client.Stop(ctx)

	if err != nil {
		return errors.WithStack(err)
	}

	p.StopTime = time.Now()

	return nil
}

func (p *Plugin) Decommission() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		return p.client.Decommission()
	}
	return nil
}

func (p *Plugin) DisDecommission() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		return p.client.DisDecommission()
	}
	return nil
}

func (p *Plugin) IsDecommissioned() bool {
	if p.client != nil {
		return p.client.IsDecommissioned()
	}
	return false
}

func (p *Plugin) Exited() bool {
	if p.client != nil {
		return p.client.Exited()
	}
	return false
}

func (p *Plugin) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	pluginClient, ok := p.Client()
	if !ok {
		return errors.New("插件没有实现CallResource接口")
	}
	return pluginClient.CallResource(ctx, req, sender)
}

func (p *Plugin) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, errors.New("插件没有实现CheckHealth接口")
	}
	return pluginClient.CheckHealth(ctx, req)
}

func (p *Plugin) Pub2Channel(ctx context.Context, req *backend.Pub2ChannelRequest) (*backend.Pub2ChannelResponse, error) {
	pluginClient, ok := p.Client()
	if !ok {
		return nil, errors.New("插件没有实现Live接口")
	}
	return pluginClient.Pub2Channel(ctx, req)
}

func (p *Plugin) RegisterClient(c backendplugin.Plugin) {
	p.client = c
}

func (p *Plugin) Client() (backendplugin.Plugin, bool) {
	if p.client != nil {
		return p.client, true
	}
	return nil, false
}

func (p *Plugin) Version() string {
	return p.pluginData.PluginJsonData.Version
}

func (p *Plugin) BackendDebug() bool {
	return p.pluginData.PluginJsonData.BackendDebug
}

func (p *Plugin) getMigratorTable() string {
	return "migrations"
}

func (p *Plugin) Migrator() {
	migrationsTable := p.getMigratorTable()
	currentPluginVersion := p.Version()
	type EvMigrateSql struct {
		Id  string `gorm:"primaryKey"`
		Sql string
	}

	type PluginInfo struct {
		Id           int `gorm:"primaryKey"`
		LocalVersion string
	}

	var evMigrateSql EvMigrateSql
	if !p.Gorm().Migrator().HasTable(&evMigrateSql) {
		p.orm.AutoMigrate(&evMigrateSql)
	}

	for _, v2 := range p.pluginData.Gormigrate.Migrations {
		v := v2
		js, _ := json.Marshal(v2)
		data := EvMigrateSql{Id: v.ID, Sql: string(js)}
		err := p.orm.Save(&data).Error
		if err != nil {
			p.Logger().Error("err", zap.Error(err))
		}
	}

	var evMigrateSqls []EvMigrateSql

	p.orm.Find(&evMigrateSqls)

	if len(evMigrateSqls) == 0 {
		return
	}

	sort.Slice(evMigrateSqls, func(i, j int) bool {
		return util.LessThan(evMigrateSqls[i].Id, evMigrateSqls[j].Id)
	})

	storeMaxVersion := evMigrateSqls[0].Id

	var plugininfo PluginInfo
	hasPluginInfo := p.Gorm().Migrator().HasTable(&plugininfo)
	if !hasPluginInfo {
		err := p.Gorm().AutoMigrate(&plugininfo)
		if err != nil {
			p.Logger().Error("err", zap.Error(err))
		}
	} else {
		err := p.Gorm().First(&plugininfo).Error
		if err != nil {
			p.Logger().Error("err", zap.Error(err))
		}
		storeMaxVersion = plugininfo.LocalVersion
	}

	if hasPluginInfo && (storeMaxVersion == currentPluginVersion) {
		return
	}

	defer func() {
		plugininfo.Id = 1
		plugininfo.LocalVersion = currentPluginVersion
		p.Gorm().Save(&plugininfo)
	}()

	migrations := []*gormigrate.Migration{}

	for _, v2 := range evMigrateSqls {
		v := v2
		var migration build.Migration
		err := json.Unmarshal([]byte(v.Sql), &migration)
		if err != nil {
			p.Logger().Error("err", zap.Error(err))
		}
		mig := &gormigrate.Migration{
			ID: v.Id,
			Migrate: func(tx *gorm.DB) error {
				migrateSqls := migration.SqliteMigrateSqls
				if p.Cfg.DbType == config.MysqlDbTyp {
					migrateSqls = migration.MysqlMigrateSqls
				}
				for _, migrateSql := range migrateSqls {
					err := tx.Exec(migrateSql.Sql, migrateSql.Args...).Error
					if err != nil {
						p.Logger().Error("err", zap.Error(err))
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {

				rollbackSqls := migration.SqliteRollback
				if p.Cfg.DbType == config.MysqlDbTyp {
					rollbackSqls = migration.MysqlRollback
				}

				for _, sql := range rollbackSqls {

					err := tx.Exec(sql.Sql, sql.Args...).Error
					if err != nil {
						p.Logger().Error("err", zap.Error(err))
					}
				}
				return nil
			},
		}
		migrations = append(migrations, mig)
	}

	m := gormigrate.New(p.orm.DB, gormigrate.DefaultOptions, migrations)

	if !p.orm.Migrator().HasTable(migrationsTable) {

		for _, v := range evMigrateSqls {

			if util.LessThan(currentPluginVersion, v.Id) {
				break
			}
			err := m.MigrateTo(v.Id)
			if err != nil {
				p.Logger().Error("err", zap.Error(err))
			}
		}
		return
	}

	isRollback := util.LessThan(currentPluginVersion, storeMaxVersion)

	if isRollback {
		sort.Slice(evMigrateSqls, func(i, j int) bool {
			return util.LessThan(evMigrateSqls[j].Id, evMigrateSqls[i].Id)
		})

		for _, v := range evMigrateSqls {
			if util.LessThan(storeMaxVersion, v.Id) {
				continue
			}
			if util.LessThan(v.Id, currentPluginVersion) {
				break
			}

			err := m.RollbackTo(v.Id)
			if err != nil {
				p.Logger().Error("err", zap.Error(err))
			}
		}
		return
	}

	sort.Slice(evMigrateSqls, func(i, j int) bool {
		return util.LessThan(evMigrateSqls[i].Id, evMigrateSqls[j].Id)
	})

	for _, v := range evMigrateSqls {
		if util.LessThan(v.Id, storeMaxVersion) {
			continue
		}
		if util.LessThan(currentPluginVersion, v.Id) {
			break
		}

		err := m.MigrateTo(v.Id)
		if err != nil {
			p.Logger().Error("err", zap.Error(err))
		}
	}

}
