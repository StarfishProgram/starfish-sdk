package sdkjwt

import (
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
	"github.com/golang-jwt/jwt/v4"
)

// Config JWT配置
type Config struct {
	Issuer      string `toml:"issuer"`      // 发行人
	SecretKey   string `toml:"secretKey"`   // 签名私钥
	ExpiresTime int64  `toml:"expiresTime"` // 失效时间(秒)
	ReissueTime int64  `toml:"reissueTime"` // 重新颁发时间(秒) : 令牌剩余时间小于该值则重新颁发新令牌
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId sdktypes.ID
	RoleId sdktypes.ID
	Pubkey string
}

type _Jwt struct {
	config *Config
}

func (j *_Jwt) NewToken(userId sdktypes.ID, roleId sdktypes.ID, pubkey string) (string, error) {
	jwtData := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Unix()+j.config.ExpiresTime, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
		},
		UserId: userId,
		RoleId: roleId,
		Pubkey: pubkey,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtData)
	tokenStr, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", sdkcodes.Internal.WithMsg(err.Error())
	}
	return tokenStr, nil
}
func (j *_Jwt) FlushToken(jwtData *UserClaims) (string, error) {
	jwtData.ExpiresAt = jwt.NewNumericDate(time.Unix(time.Now().Unix()+j.config.ExpiresTime, 0))
	jwtData.IssuedAt = jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0))
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtData)
	tokenStr, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", sdkcodes.Internal.WithMsg(err.Error())
	}
	return tokenStr, nil
}

func (j *_Jwt) ParseToken(tokenStr string) (*UserClaims, error) {
	var jwtData = new(UserClaims)
	token, err := jwt.ParseWithClaims(tokenStr, jwtData, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	})
	if err != nil {
		return nil, sdkcodes.AccessLimited.WithMsg(err.Error())
	}
	if token.Valid {
		return jwtData, nil
	}
	return nil, sdkcodes.AccessLimited
}

var ins map[string]Jwt

func init() {
	ins = make(map[string]Jwt)
}

func Init(config *Config, key ...string) {
	_ins := &_Jwt{config}
	if len(key) == 0 {
		ins[""] = _ins
	} else {
		ins[key[0]] = _ins
	}
}

func Ins(key ...string) Jwt {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
