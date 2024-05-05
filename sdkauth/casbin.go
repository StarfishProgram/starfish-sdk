package sdkauth

import (
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
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

	auth := &Auth{enforcer}
	go SyncRules(db, auth)

	if len(key) == 0 {
		ins[""] = auth
	} else {
		ins[key[0]] = auth
	}
}

func loadSettings(db *gorm.DB) *int64 {
	var json sdktypes.JSON
	query := db.Table("sys_settings")
	query.Where("`k1` = 'sys_authority_rule'")
	query.Select("`value`").Limit(1)
	if err := query.Scan(&json).Error; err != nil {
		sdklog.Warn("查询规则配置失败:", err)
		return nil
	}
	var dbId int64
	if err := json.To(&dbId); err != nil {
		sdklog.Warn("读取规则配置失败:", err)
		return nil
	}
	return &dbId
}

func SyncRules(db *gorm.DB, auth *Auth) {
	var localId int64
	if dbId := loadSettings(db); dbId != nil {
		localId = *dbId
	}
	t := time.NewTicker(10 * time.Second)
	for {
		<-t.C
		dbId := loadSettings(db)
		if dbId == nil || localId == *dbId {
			continue
		}
		sdklog.Infof("规则配置更新中[%d].....", dbId)
		if err := auth.LoadPolicy(); err != nil {
			sdklog.Warn("规则配置更新失败:", err)
			continue
		}
		localId = *dbId
		sdklog.Info("规则配置已更新:", localId)
	}
}

func Ins(key ...string) *Auth {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
