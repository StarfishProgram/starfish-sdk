package sdkredis

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
	"github.com/go-redis/redis"
)

// Config Redis配置
type Config struct {
	Host     string `toml:"host"`     // 主机
	Port     int    `toml:"port"`     // 端口
	Password string `toml:"password"` // 密码
	DB       int    `toml:"db"`       // 数据库
	Prefix   string `toml:"prefix"`   // 前缀
}

type _RedisSelectKey struct {
	ins  *_Redis
	keys []string
}

func (r *_RedisSelectKey) Keys(keys ...string) SelectKey {
	r.keys = append(r.keys, keys...)
	return r
}

func (r *_RedisSelectKey) OP() Oper {
	return &_RedisOper{ins: r}
}

type _RedisOper struct {
	ins *_RedisSelectKey
}

func (op *_RedisOper) buildKey() string {
	opKey := make([]string, len(op.ins.keys)+1)
	if op.ins.ins.prefix != "" {
		opKey = append(opKey, op.ins.ins.prefix)
	}
	if len(op.ins.keys) > 0 {
		opKey = append(opKey, op.ins.keys...)
	}
	return strings.Join(opKey, ":")
}

func (op *_RedisOper) typeIsJson(v interface{}) sdktypes.Result[bool] {
	switch v.(type) {
	case *string, *[]byte, *int,
		*int8, *int16, *int32,
		*int64, *uint, *uint8,
		*uint16, *uint32, *uint64,
		*float32, *float64, *bool:
		return sdktypes.Result[bool]{Data: false}
	case nil:
		return sdktypes.Result[bool]{Data: false, Code: sdkcodes.Internal.WithMsg("类型为nil")}
	default:
		return sdktypes.Result[bool]{Data: true}
	}
}

func (op *_RedisOper) Get(v any) sdkcodes.Code {
	opKey := op.buildKey()
	isJsonResult := op.typeIsJson(&v)
	if isJsonResult.IsError() {
		return isJsonResult.Code
	}
	result := op.ins.ins.ins.Get(opKey)
	if result.Err() != nil {
		return sdkcodes.Internal.WithMsg("%v", result.Err().Error())
	}
	if isJsonResult.Data {
		if err := json.Unmarshal([]byte(result.Val()), &v); err != nil {
			return sdkcodes.Internal.WithMsg("%v", err.Error())
		}
	} else {
		if err := result.Scan(&v); err != nil {
			return sdkcodes.Internal.WithMsg("%v", err.Error())
		}
	}
	return nil
}

func (op *_RedisOper) Set(val any, expr ...time.Duration) sdkcodes.Code {
	opKey := op.buildKey()
	isJsonResult := op.typeIsJson(&val)
	if isJsonResult.IsError() {
		return isJsonResult.Code
	}
	if isJsonResult.Data {
		jsonData, err := json.Marshal(val)
		if err != nil {
			return sdkcodes.Internal.WithMsg("%s", err.Error())
		}
		if err := op.ins.ins.ins.Set(
			opKey,
			string(jsonData),
			sdk.IfCall(
				len(expr) > 0,
				func() time.Duration { return expr[0] },
				func() time.Duration { return time.Duration(0) },
			),
		).Err(); err != nil {
			return sdkcodes.Internal.WithMsg("%s", err.Error())
		}
	} else {
		if err := op.ins.ins.ins.Set(
			opKey,
			val,
			sdk.IfCall(
				len(expr) > 0,
				func() time.Duration { return expr[0] },
				func() time.Duration { return time.Duration(0) },
			),
		).Err(); err != nil {
			return sdkcodes.Internal.WithMsg("%s", err.Error())
		}
	}
	return nil
}

func (op *_RedisOper) IncrByInt(val int64) sdktypes.Result[int64] {
	opKey := op.buildKey()
	v, err := op.ins.ins.ins.IncrBy(opKey, val).Result()
	if err != nil {
		return sdktypes.Result[int64]{Code: sdkcodes.Internal.WithMsg("%s", err.Error())}
	}
	return sdktypes.Result[int64]{Data: v}
}

func (op *_RedisOper) IncrByFloat(val float64) sdktypes.Result[float64] {
	opKey := op.buildKey()
	v, err := op.ins.ins.ins.IncrByFloat(opKey, val).Result()
	if err != nil {
		return sdktypes.Result[float64]{Code: sdkcodes.Internal.WithMsg(err.Error())}
	}
	return sdktypes.Result[float64]{Data: v}
}

func (op *_RedisOper) Del() sdkcodes.Code {
	opKey := op.buildKey()
	err := op.ins.ins.ins.Del(opKey).Err()
	if err != nil {
		return sdkcodes.Internal.WithMsg(err.Error())
	}
	return nil
}

type _Redis struct {
	ins    *redis.Client
	prefix string
}

var ins map[string]*_Redis

func init() {
	ins = make(map[string]*_Redis)
}

// Init Redis初始化
func Init(config *Config, key ...string) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	r := _Redis{
		ins:    client,
		prefix: config.Prefix,
	}
	if len(key) == 0 {
		ins[""] = &r
	} else {
		ins[key[0]] = &r
	}
}

func Ins(key ...string) SelectKey {
	return &_RedisSelectKey{
		ins: sdk.IfCall(
			len(key) == 0,
			func() *_Redis { return ins[""] },
			func() *_Redis { return ins[key[0]] },
		),
		keys: []string{},
	}
}
