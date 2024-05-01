package sdkjwt

import (
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
)

type Jwt interface {
	NewToken(userId sdktypes.ID, roleId sdktypes.ID, pubkey string) (string, error)
	FlushToken(userClaims *UserClaims) (string, error)
	ParseToken(tokenStr string) (*UserClaims, error)
	NeedFlush(userClaims *UserClaims) bool
}
