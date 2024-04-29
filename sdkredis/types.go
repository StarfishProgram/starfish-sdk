package sdkredis

import (
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
)

type SelectKey interface {
	// Keys 选择Key
	Keys(keys ...string) SelectKey
	// OP 执行操作
	OP() Oper
}

type Oper interface {
	// Get 获取值
	Get(v any) sdkcodes.Code
	// Set 设置值
	Set(val any, expr ...time.Duration) sdkcodes.Code
	// IncrByInt 递增
	IncrByInt(val int64) sdktypes.Result[int64]
	// IncrByFloat 递增
	IncrByFloat(val float64) sdktypes.Result[float64]
	// Del 删除Key
	Del() sdkcodes.Code
}
