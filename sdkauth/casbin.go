package sdkauth

import (
	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Auth struct {
	*casbin.Enforcer
	Domain string
}

var authIns map[string]*Auth

func init() {
	authIns = make(map[string]*Auth)
}

func Init(db *gorm.DB, domain string) {
	sdk.Check(domain != "", sdkcodes.Internal.WithMsg("domain不能为空"))

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

	authIns[domain] = &Auth{enforcer, domain}
}

func Ins(domain string) *Auth {
	return authIns[domain]
}
