package sqlstore

import (
	_ "embed"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	logger2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore/sqlite"
	"github.com/1340691923/ElasticView/pkg/infrastructure/sqlstore/svr_log"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"path/filepath"
)

type SqlStore struct {
	*gorm.DB
}

func newSqlStore(DB *gorm.DB) *SqlStore {
	return &SqlStore{DB: DB}
}

// NewMy 创建一个连接My的实体池
func NewSqlStore(cfg *config.Config, logger *logger2.AppLogger) (db *SqlStore, err error) {

	gromCfg := &gorm.Config{}

	if cfg.DeBug {
		gromCfg.Logger = svr_log.NewGormLogI(logger)
	}

	orm, err := gorm.Open(createDbDialector(cfg), gromCfg)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("%s连接失败", cfg.DbType))
		return
	}

	db = newSqlStore(orm)
	logger.Info("sqlStore组件初始化成功")

	return
}

func createDbDialector(cfg *config.Config) gorm.Dialector {
	var dialector gorm.Dialector
	if cfg.DbType == config.SqliteDbTyp {
		dialector = sqlite.Open(cfg.CreateDbDSN())
	} else {
		dialector = mysql.Open(cfg.CreateDbDSN())
	}

	return dialector
}

func NewPluginSqlStore(pluginStorePath, pluginName string, log hclog.Logger) (db *SqlStore, err error) {

	dsn := filepath.Join(pluginStorePath, fmt.Sprintf("%s.db", pluginName)) +
		"?_pragma=charset(utf8)&_pragma=parse_time(true)&_pragma=_busy_timeout(9999999)"

	orm, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: NewGormLogI(log),
	})

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("插件存储系统：%s连接失败", pluginName))
		return
	}

	db = newSqlStore(orm)

	log.Debug(fmt.Sprintf("插件存储系统：%s连接成功", pluginName))

	return
}
