package core

import (
	"fmt"
	"github.com/1340691923/ElasticView/engine/crontab"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"

	"github.com/1340691923/ElasticView/platform-basic-libs/rbac"

	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	sql2 "github.com/1340691923/ElasticView/sqlite"

	"log"
	"strconv"
)

// 初始化日志
func InitLogs() (fn func(), err error) {
	fn = func() {}
	logger := logs.NewLog(
		logs.WithLogPath(GlobConfig.Log.LogDir),
		logs.WithStorageDays(GlobConfig.Log.StorageDays),
	)
	logs.Logger, err = logger.InitLog()
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("日志组件初始化成功！日志所在目录：%v，保存天数为：%v", GlobConfig.Log.LogDir, GlobConfig.Log.StorageDays))
	return
}

// 初始化mysql连接
func InitSqlx() (fn func(), err error) {
	fn = func() {}
	config := GlobConfig.Mysql
	driverType := GlobConfig.DbType
	var dbSource string
	if driverType == SqliteDbTyp {
		dbSource = GlobConfig.Sqlite.DbPath + "?_loc=Local&_busy_timeout=9999999"
	} else {
		dbSource = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.Username,
			config.Pwd,
			config.IP,
			config.Port,
			config.DbName)
	}

	db.Sqlx, err = db.NewSQLX(
		driverType,
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("%v组件初始化成功！连接：%v，最大打开连接数：%v，最大等待连接数:%v",
		driverType,
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	))
	return
}

// 初始化mysql连接
func InitSqliteData() (fn func(), err error) {
	fn = func() {}
	driverType := GlobConfig.DbType
	if driverType == "sqlite3" {
		sql2.Init()
	}

	return
}

// 初始化项目启动任务
func InitTask() (fn func(), err error) {
	fn = func() {}

	esLinkModel := model.EsLinkModel{}
	if err = esLinkModel.FlushEsLinkList(); err != nil {
		return fn, err
	}

	crontab.Crontab, err = crontab.InitCrontab()
	if err != nil {
		return fn, err
	}
	fn = func() {
		crontab.Crontab.Stop()
	}
	return fn, err
}

// 初始化项目启动任务
func InitRbac() (fn func(), err error) {
	fn = func() {}
	config := GlobConfig.Mysql
	driverType := GlobConfig.DbType
	var dbSource string
	if driverType == SqliteDbTyp {
		dbSource = GlobConfig.Sqlite.DbPath + "?_loc=Local&_busy_timeout=9999999"
	} else {
		dbSource = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			config.Username,
			config.Pwd,
			config.IP,
			config.Port,
			config.DbName)
	}

	err = rbac.Run(driverType, dbSource)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Rbac组件初始化成功！连接：%v",
		dbSource,
	))
	return
}

func InitOpenWinBrowser() (fn func(), err error) {
	fn = func() {}
	config := GlobConfig
	if !config.DeBug {
		port := ":" + strconv.Itoa(config.Port)
		uri := fmt.Sprintf("%s%s", "http://127.0.0.1", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v",
			uri,
		))
	}
	return
}

