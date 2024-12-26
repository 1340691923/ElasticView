package access_control

import (
	"fmt"
	"github.com/1340691923/ElasticView/pkg/infrastructure/config"
	"github.com/1340691923/ElasticView/pkg/infrastructure/logger"
	"github.com/1340691923/ElasticView/pkg/infrastructure/orm"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var text = `
[request_definition]
        r = sub, obj, act

[policy_definition]
        p = sub, obj, act

[policy_effect]
        e = some(where (p.eft == allow))

[matchers]
        m = r.sub == p.sub && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*") || r.sub == "1"
	`

type Rbac struct {
	cfg      *config.Config
	log      *logger.AppLogger
	enforcer *casbin.SyncedEnforcer
	orm      *orm.Gorm
}

// 初始化项目启动任务
func NewRbac(cfg *config.Config, log *logger.AppLogger, orm *orm.Gorm) (*Rbac, error) {
	obj := &Rbac{
		cfg:      cfg,
		log:      log,
		orm:      orm,
		enforcer: nil,
	}

	var dbSource = cfg.CreateDbDSN()

	obj.log = log.Named("rbac")

	policy, err := model.NewModelFromString(text)
	if err != nil {
		return nil, err
	}

	adapter, err := gormadapter.NewAdapterByDB(orm.DB)
	if err != nil {
		return nil, err
	}
	obj.enforcer, err = casbin.NewSyncedEnforcer(policy, adapter)

	if err != nil {
		return nil, err
	}
	obj.log.Info(fmt.Sprintf("Rbac组件初始化成功！连接：%v", dbSource))
	return obj, nil
}

func (this *Rbac) LoadPolicy() error {
	return this.enforcer.LoadPolicy()
}

func (this *Rbac) Enforce(rvals ...interface{}) (bool, error) {
	return this.enforcer.Enforce(rvals...)
}

func (this *Rbac) AddPolicy(params ...interface{}) (bool, error) {
	return this.enforcer.AddPolicy(params...)
}

func (this *Rbac) RemoveFilteredPolicy(fieldIndex int, fieldValues ...string) (bool, error) {
	return this.enforcer.RemoveFilteredPolicy(fieldIndex, fieldValues...)
}

func (this *Rbac) SavePolicy() error {
	return this.enforcer.SavePolicy()
}
