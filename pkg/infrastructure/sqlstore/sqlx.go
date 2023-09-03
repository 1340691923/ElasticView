package sqlstore

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/logoove/sqlite"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SqlStore struct {
	DB     *sqlx.DB
	logger *logger.AppLogger
}

func (this *SqlStore) Exec(query string, args ...any) (sql.Result, error) {
	return this.DB.Exec(query, args...)
}

func (this *SqlStore) Get(dest interface{}, query string, args ...interface{}) error {
	return this.DB.Get(dest, query, args...)
}

func (this *SqlStore) Select(dest interface{}, query string, args ...interface{}) error {
	return this.DB.Select(dest, query, args...)
}

func (this *SqlStore) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return this.DB.Query(query, args...)
}

func (this *SqlStore) Close() error {
	return this.DB.Close()
}

func newSqlStore(DB *sqlx.DB, logger *logger.AppLogger) *SqlStore {
	return &SqlStore{DB: DB, logger: logger}
}

// NewMy 创建一个连接My的实体池
func NewSqlStore(cfg *config.Config, logger *logger.AppLogger) (db *SqlStore, err error) {

	sqlStore, err := sqlx.Open(cfg.DbType, cfg.CreateDbSource())
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("%s连接失败:连接信息:%s", cfg.DbType, cfg.CreateDbSource()))
		return
	}

	if cfg.Mysql.MaxOpenConns > 0 {
		sqlStore.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)
	}

	if cfg.Mysql.MaxIdleConns > 0 {
		sqlStore.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	}

	err = sqlStore.Ping()
	if err != nil {
		log.Println("err", zap.Error(err))
		return
	}
	go func() {
		for {
			err = sqlStore.Ping()
			if err != nil {
				log.Println(" db can't connect!")
			}
			time.Sleep(time.Minute)
		}
	}()
	db = newSqlStore(sqlStore, logger)
	if cfg.DbType == config.SqliteDbTyp {
		db.initSqliteDb()
	}

	log.Println("sqlStore组件初始化成功", cfg.CreateDbSource())

	return
}

//go:embed es_view.sql
var SqlByte []byte

// 初始化sqlite数据
func (this *SqlStore) initSqliteDb() {

	currDir := util.GetCurrentDirectory()

	dataDir := filepath.Join(currDir, "data")

	lockFile := filepath.Join(dataDir, "lock")

	if util.CheckFileIsExist(lockFile) {
		return
	}

	execSqlArr := strings.Split(util.Bytes2str(SqlByte), ";")

	var err error

	for _, execSql := range execSqlArr {
		log.Println("insert sql", execSql)
		_, err = this.Exec(execSql)
		if err != nil {
			log.Println(fmt.Sprintf("初始化 sqlite 执行建表语句sql:%v失败:%s", execSql, err.Error()))
			panic(err)
		}
	}

	log.Println("初始化sqlite数据完成！")
	if !util.CheckFileIsExist(dataDir) {
		os.MkdirAll(dataDir, os.ModePerm)
	}
	os.Create(lockFile)
}
