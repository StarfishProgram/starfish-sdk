package sdkredis

import "testing"

func TestRedis(t *testing.T) {
	Init(&Config{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "abc123ABC",
		DB:       0,
		Prefix:   "starfish",
	})

	result := Ins().Keys("aqiang").Keys("token").OP().Set(123)
	t.Log("写入Redis :", result)

	var val int64
	result = Ins().Keys("aqiang", "token").OP().Get(&val)
	t.Log("读取Redis, 结果 :", result, " 值 :", val)
}
