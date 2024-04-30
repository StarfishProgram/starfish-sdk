package sdkdb

import (
	"testing"

	"github.com/StarfishProgram/starfish-sdk/sdklog"
)

func TestDB(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	Init(&Config{
		Host:        "127.0.0.1",
		Port:        3307,
		User:        "root",
		Password:    "husky123456",
		Database:    "starfish",
		Config:      "charset=utf8mb4&parseTime=True&loc=Local",
		MaxIdleConn: 0,
		MaxOpenConn: 0,
		MaxLifetime: 0,
		ShowSql:     true,
		SlowTime:    0,
	})

	var version string
	result := Ins().Raw("select versoin()").Take(&version)
	t.Log("result :", result)
}
