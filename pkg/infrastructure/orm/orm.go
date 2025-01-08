package orm

import (
	_ "embed"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	logger2 "github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm/sqlite"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm/svr_log"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Gorm struct {
	*gorm.DB
}

func newGorm(DB *gorm.DB) *Gorm {
	return &Gorm{DB: DB}
}

// NewMy 创建一个连接My的实体池
func NewGorm(cfg *config.Config, logger *logger2.AppLogger) (db *Gorm, err error) {

	gromCfg := &gorm.Config{}

	if cfg.DeBug {
		gromCfg.Logger = svr_log.NewGormLogI(logger)
	}

	orm, err := gorm.Open(createDbDialector(cfg.DbType, cfg.CreateDbDSN()), gromCfg)

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("%s连接失败", cfg.DbType))
		return
	}

	db = newGorm(orm)

	// 获取底层的 sql.DB 对象
	sqlDB, err := db.DB.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100)            // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(10)             // 设置最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)   // 设置连接的最大生命周期
	sqlDB.SetConnMaxIdleTime(time.Minute) // 设置空闲连接的最大存活时间

	logger.Info("sqlStore组件初始化成功")

	return
}

func createDbDialector(dbTyp string, dbSource string) gorm.Dialector {
	var dialector gorm.Dialector

	if dbTyp == config.SqliteDbTyp {
		dialector = sqlite.Open(dbSource)
	} else {
		dialector = mysql.Open(dbSource)
	}

	return dialector
}

func NewPluginGorm(evOrm *Gorm, cfg *config.Config, dbName string, pluginStorePath string, hlog hclog.Logger) (db *Gorm, err error) {
	var dsn string
	var dialector gorm.Dialector
	if cfg.DbType == config.SqliteDbTyp {
		dsn = pluginStorePath + "?_pragma=charset(utf8)&_pragma=parse_time(true)&_pragma=_busy_timeout(9999999)"
		dialector = sqlite.Open(dsn)

	} else if cfg.DbType == config.MysqlDbTyp {
		err = evOrm.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
		if err != nil {
			return
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.Mysql.Username,
			cfg.Mysql.Pwd,
			cfg.Mysql.IP,
			cfg.Mysql.Port,
			dbName,
		)
		dialector = mysql.Open(dsn)
	} else {
		errors.New("未知db类型")
		return
	}

	orm, err := gorm.Open(dialector, &gorm.Config{
		Logger: NewGormLogI(hlog),
	})

	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("插件存储系统：%s连接失败", pluginStorePath))
		return
	}

	db = newGorm(orm)

	// 获取底层的 sql.DB 对象
	sqlDB, err := db.DB.DB()
	if err != nil {
		panic("failed to get sql.DB")
	}

	// 设置连接池参数
	if cfg.DbType == config.SqliteDbTyp {
		sqlDB.SetMaxOpenConns(1)    // 设置最大打开连接数为 1
		sqlDB.SetMaxIdleConns(1)    // 设置最大空闲连接数为 1
		sqlDB.SetConnMaxLifetime(0) // 禁用连接过期时间
	} else if cfg.DbType == config.MysqlDbTyp {
		sqlDB.SetMaxOpenConns(10) // 设置最大打开连接数为 10
		sqlDB.SetMaxIdleConns(10) // 设置最大空闲连接数为 10
	}

	hlog.Debug(fmt.Sprintf("插件存储系统：%s连接成功", pluginStorePath))

	return
}
