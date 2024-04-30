package sdkwebmiddleware

import (
	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var authIns map[string]*casbin.Enforcer

func init() {
	authIns = make(map[string]*casbin.Enforcer)
}

func Auth(db *gorm.DB, key ...string) {
	config, err := model.NewModelFromString(`
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
	`)
	sdk.CheckError(err)

	gormadapter.TurnOffAutoMigrate(db)
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "sys", "auth_rule")
	sdk.CheckError(err)

	enforcer, err := casbin.NewEnforcer(config, adapter)
	sdk.CheckError(err)

	if len(key) == 0 {
		authIns[""] = enforcer
	} else {
		authIns[key[0]] = enforcer
	}
}

func AuthIns(key ...string) *casbin.Enforcer {
	if len(key) == 0 {
		return authIns[""]
	} else {
		return authIns[key[0]]
	}
}
