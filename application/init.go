package application

import (
	"fmt"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/engine/logs"
	"github.com/1340691923/ElasticView/model"
	"github.com/1340691923/ElasticView/platform-basic-libs/rbac"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"log"
	"strconv"
)


// 初始化日志
func  InitLogs() (err error) {
	logger := logs.NewLog(
		logs.WithLogPath(GlobConfig.Log.LogDir),
		logs.WithStorageDays(GlobConfig.Log.StorageDays),
	)
	logs.Logger, err = logger.InitLog()
	if err!=nil{
		return
	}
	log.Println(fmt.Sprintf("日志组件初始化成功！日志所在目录：%v，保存天数为：%v",GlobConfig.Log.LogDir,GlobConfig.Log.StorageDays))
	return
}

// 初始化mysql连接
func InitMysql() (err error) {
	config := GlobConfig.Mysql
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Pwd,
		config.IP,
		config.Port,
		config.DbName)
	db.Sqlx, err = db.NewSQLX(
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
	)
	if err!=nil{
		return
	}
	log.Println(fmt.Sprintf("Mysql组件初始化成功！连接：%v，最大打开连接数：%v，最大等待连接数:%v",
		dbSource,
		config.MaxOpenConns,
		config.MaxIdleConns,
		))
	return
}

// 初始化项目启动任务
func InitTask() (err error) {
	esLinkModel := model.EsLinkModel{}
	if err = esLinkModel.FlushEsLinkList(); err != nil {
		return err
	}
	return err
}

// 初始化项目启动任务
func  InitRbac() (err error) {
	config := GlobConfig.Mysql
	dbSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Pwd,
		config.IP,
		config.Port,
		config.DbName)
	err = rbac.Run("mysql",dbSource)
	if err!=nil{
		return
	}
	log.Println(fmt.Sprintf("Rbac组件初始化成功！连接：%v",
		dbSource,
	))
	return
}

func InitOpenWinBrowser()(err error){
	config := GlobConfig
	if !config.DeBug{
		port := ":" + strconv.Itoa(config.Port)
		uri := fmt.Sprintf("%s%s", "http://127.0.0.1", port)
		util.OpenWinBrowser(uri)
		log.Println(fmt.Sprintf("将打开浏览器！地址为：%v",
			uri,
		))
	}
	return
}
