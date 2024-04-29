package sdktypes

import "github.com/StarfishProgram/starfish-sdk/sdkcodes"

// Result 结果
type Result[D any] struct {
	// 状态码
	Code sdkcodes.Code
	// 数据
	Data D
}

// IsOk 是否成功
func (r *Result[D]) IsOk() bool {
	return r.Code == nil
}

// IsError 是否错误
func (r *Result[D]) IsError() bool {
	return !r.IsOk()
}
