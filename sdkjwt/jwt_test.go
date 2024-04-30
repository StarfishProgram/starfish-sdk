package sdkjwt

import (
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	Init(&Config{
		Issuer:      "starfish",
		SecretKey:   "abc123ABC",
		ExpiresTime: 2592000,
		ReissueTime: 604800,
	})
	// 颁发token
	tokenStr, err := Ins().NewToken(1, 2, "3")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("颁发token :", tokenStr)

	// 解析token
	userClaims, err := Ins().ParseToken(tokenStr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("解析token :", userClaims)

	// token续签
	time.Sleep(time.Second * 2)
	flushToken, err := Ins().FlushToken(userClaims)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token续签 :", flushToken)
}
