package sdkauth

import (
	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Auth struct {
	*casbin.Enforcer
}

var ins map[string]*Auth

func init() {
	ins = make(map[string]*Auth)
}

func Init(db *gorm.DB, key ...string) {
	casbinConfig, err := model.NewModelFromString(`
[request_definition]
r = sub, obj

[policy_definition]
p = sub, obj, status

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, '1') && r.obj == p.obj && p.status == '1'
	`)
	sdk.AssertError(err)

	gormadapter.TurnOffAutoMigrate(db)
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "sys", "authority_rule")
	sdk.AssertError(err)

	enforcer, err := casbin.NewEnforcer(casbinConfig, adapter)
	sdk.AssertError(err)

	if len(key) == 0 {
		ins[""] = &Auth{enforcer}
	} else {
		ins[key[0]] = &Auth{enforcer}
	}
}

func Ins(key ...string) *Auth {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
