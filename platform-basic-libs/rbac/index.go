package rbac

import (
	"github.com/casbin/casbin"
	xormadapter "github.com/casbin/xorm-adapter"
)

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
	if err != nil {
		return
	}
	return
}
