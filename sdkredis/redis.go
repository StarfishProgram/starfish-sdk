package sdkredis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config Redis配置
type Config struct {
	Host     string `toml:"host"`     // 主机
	Port     int    `toml:"port"`     // 端口
	Password string `toml:"password"` // 密码
	DB       int    `toml:"db"`       // 数据库
	Prefix   string `toml:"prefix"`   // 前缀
}

type Redis struct {
	*redis.Client
	Prefix string
}

var ins map[string]*Redis

func init() {
	ins = make(map[string]*Redis)
}

// Init Redis初始化
func Init(config *Config, key ...string) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})
	r := Redis{
		Client: client,
		Prefix: config.Prefix,
	}
	if len(key) == 0 {
		ins[""] = &r
	} else {
		ins[key[0]] = &r
	}
}

// Ins 获取数据源
func Ins(key ...string) *Redis {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
