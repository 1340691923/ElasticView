package rbac

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/core"
	"github.com/1340691923/ElasticView/pkg/engine/config"
	"github.com/1340691923/ElasticView/pkg/util"
	"github.com/casbin/casbin"
	xormadapter "github.com/casbin/xorm-adapter"
	"log"
	"os"
	"path/filepath"
)

func init() {
	core.Register(core.LastLevel, "rbac", InitRbac)
}

// 初始化项目启动任务
func InitRbac() (fn func(), err error) {
	fn = func() {}
	mysql := config.GlobConfig.Mysql
	driverType := config.GlobConfig.DbType
	var dbSource string
	if driverType == config.SqliteDbTyp {
		currDir := util.GetCurrentDirectory()
		dataDir := filepath.Join(currDir, "data")
		if !util.CheckFileIsExist(dataDir) {
			os.MkdirAll(dataDir, os.ModePerm)
		}
		dbSource = filepath.Join(dataDir, config.GlobConfig.Sqlite.DbPath) + "?_loc=Local&_busy_timeout=9999999"
	} else {
		dbSource = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			mysql.Username,
			mysql.Pwd,
			mysql.IP,
			mysql.Port,
			mysql.DbName)
	}

	err = Run(driverType, dbSource)
	if err != nil {
		return
	}
	log.Println(fmt.Sprintf("Rbac组件初始化成功！连接：%v", dbSource))
	return
}

var Enforcer *casbin.Enforcer

func Run(driverName string, datasource string) (err error) {
	text := `
[request_definition]
        r = sub, obj, act

[policy_definition]
        p = sub, obj, act

[policy_effect]
        e = some(where (p.eft == allow))

[matchers]
        m = r.sub == p.sub && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*") || r.sub == "1"
	`
	policy := casbin.NewModel(text)

	model := xormadapter.NewAdapter(driverName, datasource, true)
	Enforcer = casbin.NewEnforcer(policy, model)
	err = Enforcer.LoadPolicy()
	return
}
