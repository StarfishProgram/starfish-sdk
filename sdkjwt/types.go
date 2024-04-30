package sdkjwt

import (
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
)

type Jwt interface {
	// NewToken 颁发一个新的Token
	NewToken(userId sdktypes.ID, roleId sdktypes.ID, pubkey string) (string, error)
	// FlushToken 根据旧的令牌信息颁发一个新的Token
	FlushToken(jwtData *UserClaims) (string, error)
	// ParseToken 解析Token
	ParseToken(tokenStr string) (*UserClaims, error)
}
