package starfish_sdk

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Issuer      string `toml:"issuer"`      // 发行人
	SecretKey   string `toml:"secretKey"`   // 签名私钥
	ExpiresTime int64  `toml:"expiresTime"` // 失效时间(秒)
	ReissueTime int64  `toml:"reissueTime"` // 重新颁发时间(秒) : 令牌剩余时间小于该值则重新颁发新令牌
}

type JWTData struct {
	jwt.RegisteredClaims
	Data map[string]any
}

// IJWT JWT
type IJWT interface {
	// NewToken 颁发一个新的Token
	NewToken(data map[string]any) Result[string]
	// FlushToken 根据旧的令牌信息颁发一个新的Token
	FlushToken(jwtData *JWTData) Result[string]
	// ParseToken 解析Token
	ParseToken(tokenStr string) Result[*JWTData]
}

type _jwt struct {
	config *JWTConfig
}

func (j *_jwt) NewToken(data map[string]any) Result[string] {
	jwtData := JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.config.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Unix(time.Now().Unix()+j.config.ExpiresTime, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0)),
		},
		Data: data,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenStr, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return Result[string]{Code: CodeServerError.WithMsg(err.Error())}
	}
	return Result[string]{Data: tokenStr}
}
func (j *_jwt) FlushToken(jwtData *JWTData) Result[string] {
	jwtData.ExpiresAt = jwt.NewNumericDate(time.Unix(time.Now().Unix()+j.config.ExpiresTime, 0))
	jwtData.IssuedAt = jwt.NewNumericDate(time.Unix(time.Now().Unix(), 0))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenStr, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return Result[string]{Code: CodeServerError.WithMsg(err.Error())}
	}
	return Result[string]{Data: tokenStr}
}

func (j *_jwt) ParseToken(tokenStr string) Result[*JWTData] {
	var jwtData = new(JWTData)
	token, err := jwt.ParseWithClaims(tokenStr, jwtData, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	})
	if err != nil {
		return Result[*JWTData]{Code: CodeAccessInvlid.WithMsg(err.Error())}
	}
	if token.Valid {
		return Result[*JWTData]{Data: jwtData}
	}
	return Result[*JWTData]{Code: CodeAccessInvlid}
}

var jwtIns map[string]*_jwt

func init() {
	jwtIns = make(map[string]*_jwt)
}

// InitJWT JWT初始化
func InitJWT(config *JWTConfig, key ...string) {
	ins := &_jwt{config}
	if len(key) == 0 {
		jwtIns[""] = ins
	} else {
		jwtIns[key[0]] = ins
	}
}

// JWT 获取JWT
func JWT(key ...string) IJWT {
	if len(key) == 0 {
		return jwtIns[""]
	} else {
		return jwtIns[key[0]]
	}
}
