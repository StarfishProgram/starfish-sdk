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
	tokenResult := Ins().NewToken(map[string]any{
		"userId":   666,
		"roleId":   777,
		"username": "阿强",
	})
	if tokenResult.Code != nil {
		t.Fatal(tokenResult.Code)
	}
	t.Log("颁发token :", tokenResult.Data)

	// 解析token
	claimsResult := Ins().ParseToken(tokenResult.Data)
	if claimsResult.Code != nil {
		t.Fatal(claimsResult.Code)
	}
	t.Log("解析token :", claimsResult.Data)

	// token续签
	time.Sleep(time.Second * 2)
	flushTokenResult := Ins().FlushToken(claimsResult.Data)
	if flushTokenResult.Code != nil {
		t.Fatal(flushTokenResult.Code)
	}
	t.Log("token续签 :", flushTokenResult.Data)
}
