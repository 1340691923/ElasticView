package sqlite

import (
	_ "embed"
	"fmt"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/1340691923/ElasticView/pkg/engine/config"
	"github.com/1340691923/ElasticView/pkg/engine/db"
	"github.com/1340691923/ElasticView/pkg/util"
	"os"
	"path/filepath"

	"log"
	"strings"
)

//go:embed es_view.sql
var SqlByte []byte

func init() {
	core.Register(core.MinLevel, "sqlite3", InitSqliteData)
}

func InitSqliteData() (fn func(), err error) {
	fn = func() {}
	driverType := config.GlobConfig.DbType
	if driverType == "sqlite3" {
		Init()
	}
	return
}

// 初始化sqlite数据
func Init() {

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
		_, err = db.Sqlx.Exec(execSql)
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
