package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `toml:"host"`     // 主机
	Port     int    `toml:"port"`     // 端口
	Password string `toml:"password"` // 密码
	DB       int    `toml:"db"`       // 数据库
	Prefix   string `toml:"prefix"`   // 前缀
}

type IRedisSelectKey interface {
	// Keys 选择Key
	Keys(keys ...string) IRedisSelectKey
	// OP 执行操作
	OP() IRedisOP
}

type redisSelectKey struct {
	ins  *_redis
	keys []string
}

func (r *redisSelectKey) Keys(keys ...string) IRedisSelectKey {
	r.keys = append(r.keys, keys...)
	return r
}

func (r *redisSelectKey) OP() IRedisOP {
	return &redisOP{ins: r}
}

type IRedisOP interface {
	// Get 获取值
	Get(v any) ICode
	// Set 设置值
	Set(val any, expr ...time.Duration) ICode
	// IncrByInt 递增
	IncrByInt(val int64) Result[int64]
	// IncrByFloat 递增
	IncrByFloat(val float64) Result[float64]
	// Del 删除Key
	Del() ICode
}
type redisOP struct {
	ins *redisSelectKey
}

func (op *redisOP) buildKey() string {
	opKey := make([]string, len(op.ins.keys)+1)
	if op.ins.ins.prefix != "" {
		opKey = append(opKey, op.ins.ins.prefix)
	}
	if len(op.ins.keys) > 0 {
		opKey = append(opKey, op.ins.keys...)
	}
	return strings.Join(opKey, ":")
}

func (op *redisOP) typeIsJson(v interface{}) Result[bool] {
	switch v.(type) {
	case *string, *[]byte, *int,
		*int8, *int16, *int32,
		*int64, *uint, *uint8,
		*uint16, *uint32, *uint64,
		*float32, *float64, *bool:
		return Result[bool]{Data: false}
	case nil:
		return Result[bool]{Data: false, Code: CodeServerError.WithMsg("类型为nil")}
	default:
		return Result[bool]{Data: true}
	}
}

func (op *redisOP) Get(v any) ICode {
	opKey := op.buildKey()
	isJsonResult := op.typeIsJson(&v)
	if isJsonResult.IsFaild() {
		return isJsonResult.Code
	}
	result := op.ins.ins.ins.Get(opKey)
	if result.Err() != nil {
		return CodeServerError.WithMsg(result.Err().Error())
	}
	if isJsonResult.Data {
		if err := json.Unmarshal([]byte(result.Val()), &v); err != nil {
			return CodeServerError.WithMsg(err.Error())
		}
	} else {
		if err := result.Scan(&v); err != nil {
			return CodeServerError.WithMsg(err.Error())
		}
	}
	return nil
}

func (op *redisOP) Set(val any, expr ...time.Duration) ICode {
	opKey := op.buildKey()
	isJsonResult := op.typeIsJson(&val)
	if isJsonResult.IsFaild() {
		return isJsonResult.Code
	}
	if isJsonResult.Data {
		jsonData, err := json.Marshal(val)
		if err != nil {
			return CodeServerError.WithMsg(err.Error())
		}
		if err := op.ins.ins.ins.Set(
			opKey,
			string(jsonData),
			IfCall(
				len(expr) > 0,
				func() time.Duration { return expr[0] },
				func() time.Duration { return time.Duration(0) },
			),
		).Err(); err != nil {
			return CodeServerError.WithMsg(err.Error())
		}
	} else {
		if err := op.ins.ins.ins.Set(
			opKey,
			val,
			IfCall(
				len(expr) > 0,
				func() time.Duration { return expr[0] },
				func() time.Duration { return time.Duration(0) },
			),
		).Err(); err != nil {
			return CodeServerError.WithMsg(err.Error())
		}
	}
	return nil
}

func (op *redisOP) IncrByInt(val int64) Result[int64] {
	opKey := op.buildKey()
	v, err := op.ins.ins.ins.IncrBy(opKey, val).Result()
	if err != nil {
		return Result[int64]{Code: CodeServerError.WithMsg(err.Error())}
	}
	return Result[int64]{Data: v}
}

func (op *redisOP) IncrByFloat(val float64) Result[float64] {
	opKey := op.buildKey()
	v, err := op.ins.ins.ins.IncrByFloat(opKey, val).Result()
	if err != nil {
		return Result[float64]{Code: CodeServerError.WithMsg(err.Error())}
	}
	return Result[float64]{Data: v}
}

func (op *redisOP) Del() ICode {
	opKey := op.buildKey()
	err := op.ins.ins.ins.Del(opKey).Err()
	if err != nil {
		return CodeServerError.WithMsg(err.Error())
	}
	return nil
}

type _redis struct {
	ins    *redis.Client
	prefix string
}

var redisIns map[string]*_redis

func init() {
	redisIns = make(map[string]*_redis)
}

func Redis(key ...string) IRedisSelectKey {
	return &redisSelectKey{
		ins: IfCall(
			len(key) == 0,
			func() *_redis { return redisIns[""] },
			func() *_redis { return redisIns[key[0]] },
		),
		keys: []string{},
	}
}

// InitRedis Redis初始化
func InitRedis(config *RedisConfig, key ...string) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	r := _redis{
		ins:    client,
		prefix: config.Prefix,
	}
	if len(key) == 0 {
		redisIns[""] = &r
	} else {
		redisIns[key[0]] = &r
	}
}
