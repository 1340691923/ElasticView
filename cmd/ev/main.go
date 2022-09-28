package main

import (
	"github.com/1340691923/ElasticView/pkg/core"
	_ "github.com/1340691923/ElasticView/pkg/engine/config"
	_ "github.com/1340691923/ElasticView/pkg/engine/crontab"
	_ "github.com/1340691923/ElasticView/pkg/engine/db"
	_ "github.com/1340691923/ElasticView/pkg/engine/router"
	_ "github.com/1340691923/ElasticView/pkg/engine/sqlite"

	"github.com/1340691923/ElasticView/pkg/engine/logs"
	"github.com/1340691923/ElasticView/pkg/util"
)



// By 肖文龙
func main() {
	core.Run()
	defer core.Stop()
	util.WaitQuit(func() {
		logs.Logger.Info("ElasticView http服务停止成功...")
	})
}
