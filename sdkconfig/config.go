package sdkconfig

import (
	"strings"

	"github.com/spf13/viper"
)

type _Config[C any] struct {
	ins *viper.Viper
}

func (c *_Config[C]) Path(path string) Config[C] {
	c.ins.AddConfigPath(path)
	return c
}

func (c *_Config[C]) Name(name string) Config[C] {
	c.ins.SetConfigName(name)
	return c
}

func (c *_Config[C]) Type(type_ string) Config[C] {
	c.ins.SetConfigType(type_)
	return c
}

func (c *_Config[C]) EnvPrefix(prefix string) Config[C] {
	c.ins.SetEnvPrefix(prefix)
	return c
}

func (c *_Config[C]) Load() *C {
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

// New 创建配置
func New[C any]() Config[C] {
	ins := viper.New()
	ins.AddConfigPath("./")
	ins.SetConfigName("config")
	ins.SetConfigType("toml")
	ins.SetEnvPrefix("APP")
	return &_Config[C]{ins}
}
