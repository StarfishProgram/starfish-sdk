package sdkauth

import (
	"fmt"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/StarfishProgram/starfish-sdk/sdkredis"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Config struct {
	AutoSync        bool             `toml:"autoSync"`
	SyncIgnoreSelf  bool             `toml:"syncIgnoreSelf"`
	RedisSyncConfig *sdkredis.Config `toml:"redisSyncConfig"`
}

type Auth struct {
	*casbin.Enforcer
}

var ins map[string]*Auth

func init() {
	ins = make(map[string]*Auth)
}

func Init(db *gorm.DB, config *Config, key ...string) {
	casbinConfig, err := model.NewModelFromString(`
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

	enforcer, err := casbin.NewEnforcer(casbinConfig, adapter)
	sdk.CheckError(err)

	if config.AutoSync {
		watcher, err := rediswatcher.NewWatcher(
			fmt.Sprintf("%s:%v", config.RedisSyncConfig.Host, sdk.IfNil(config.RedisSyncConfig.Port, 6379)),
			rediswatcher.WatcherOptions{
				Options: redis.Options{
					Password: config.RedisSyncConfig.Password,
				},
				Channel:    "/casbin_auth_rule",
				IgnoreSelf: config.SyncIgnoreSelf,
			},
		)
		sdk.CheckError(err)

		err = watcher.SetUpdateCallback(syncCallback)
		sdk.CheckError(err)

		err = enforcer.SetWatcher(watcher)
		sdk.CheckError(err)
	}

	if len(key) == 0 {
		ins[""] = &Auth{enforcer}
	} else {
		ins[key[0]] = &Auth{enforcer}
	}
}

func syncCallback(msg string) {
	sdklog.Debug(msg)
}

func Ins(key ...string) *Auth {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
