package sdkjwt

import (
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
)

type Jwt interface {
	// NewToken 颁发一个新的Token
	NewToken(data map[string]any) sdktypes.Result[string]
	// FlushToken 根据旧的令牌信息颁发一个新的Token
	FlushToken(jwtData *Data) sdktypes.Result[string]
	// ParseToken 解析Token
	ParseToken(tokenStr string) sdktypes.Result[*Data]
}
