package sqlite

import (
	_ "embed"
	"fmt"
	"github.com/1340691923/ElasticView/engine/db"
	"github.com/1340691923/ElasticView/platform-basic-libs/util"
	"os"

	"log"
	"strings"
)

//go:embed es_view.sql
var SqlByte []byte

//初始化sqlite数据
func Init() {

	if util.CheckFileIsExist("lock"){
		return
	}

	execSqlArr := strings.Split(util.Bytes2str(SqlByte), ";")

	var err error

	for _, execSql := range execSqlArr {
		_, err = db.Sqlx.Exec(execSql)
		if err != nil {
			log.Println(fmt.Sprintf("初始化 sqlite 执行建表语句sql:%v失败:%s", execSql, err.Error()))
			panic(err)
		}
	}

	log.Println("初始化sqlite数据完成！")
	os.Create("lock")
}
