package starfish_sdk

import (
	"strings"

	"github.com/spf13/viper"
)

// Config 配置
type Config[C any] interface {
	// Path 设置路径
	Path(path string) Config[C]
	// Name 设置文件名
	Name(name string) Config[C]
	// 设置文件类型
	Type(type_ string) Config[C]
	// 设置环境变量前缀
	EnvPrefix(prefix string) Config[C]
	// 载入配置
	Load() *C
}

type config[C any] struct {
	ins *viper.Viper
}

func (c *config[C]) Path(path string) Config[C] {
	c.ins.AddConfigPath(path)
	return c
}

func (c *config[C]) Name(name string) Config[C] {
	c.ins.SetConfigName(name)
	return c
}

func (c *config[C]) Type(type_ string) Config[C] {
	c.ins.SetConfigType(type_)
	return c
}

func (c *config[C]) EnvPrefix(prefix string) Config[C] {
	c.ins.SetEnvPrefix(prefix)
	return c
}

func (c *config[C]) Load() *C {
	c.ins.AutomaticEnv()
	c.ins.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := c.ins.ReadInConfig(); err != nil {
		panic(err)
	}
	var r C
	if err := c.ins.Unmarshal(&r); err != nil {
		panic(err)
	}
	return &r
}

// NewConfig 创建配置
func NewConfig[C any]() Config[C] {
	ins := viper.New()
	ins.AddConfigPath("./")
	ins.SetConfigName("config")
	ins.SetConfigType("toml")
	ins.SetEnvPrefix("APP")
	return &config[C]{ins}
}
