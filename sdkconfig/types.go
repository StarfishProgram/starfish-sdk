package sdkconfig

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
