package starfish_sdk

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DatabaseConfig 数据源配置
type DatabaseConfig struct {
	Host        string `toml:"host"`        // 主机
	Port        int    `toml:"port"`        // 端口
	User        string `toml:"user"`        // 用户名
	Password    string `toml:"password"`    // 密码
	Database    string `toml:"database"`    // 数据库
	Config      string `toml:"config"`      // 连接属性
	MaxIdleConn int    `toml:"maxIdleConn"` // 最大空闲连接数
	MaxOpenConn int    `toml:"maxOpenConn"` // 最大连接数
	MaxLifetime int64  `toml:"maxLifetime"` // 最大超时时间(秒)
	ShowSql     bool   `toml:"showSql"`     // 显示执行SQL
	SlowTime    int64  `toml:"slowTime"`    // 慢查询时间(毫秒)
}

type databaseWriter struct{}

func (*databaseWriter) Printf(format string, v ...interface{}) {
	Log().Infof(format, v...)
}

var databaseIns map[string]*gorm.DB

func init() {
	databaseIns = make(map[string]*gorm.DB)
}

// InitDatabase 数据源初始化
func InitDatabase(config *DatabaseConfig, key ...string) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Config,
	)

	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         1000,
		SkipInitializeWithVersion: false,
	}

	dbLog := logger.New(
		&databaseWriter{},
		logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(config.SlowTime),
			LogLevel:                  If(config.ShowSql, logger.Info, logger.Error),
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		CreateBatchSize:        1000,
		Logger:                 dbLog,
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
	})

	if err != nil {
		panic(CodeServerError.WithMsgf("数据源连接创建失败 : %s", err.Error()))
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(CodeServerError.WithMsgf("数据源连接创建失败 : %s", err.Error()))
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.MaxLifetime) * time.Second)

	if len(key) == 0 {
		databaseIns[""] = db
	} else {
		databaseIns[key[0]] = db
	}
}

// Database 获取数据源
func Database(key ...string) *gorm.DB {
	if len(key) == 0 {
		return databaseIns[""]
	} else {
		return databaseIns[key[0]]
	}
}
