package sdkwebmiddleware

import (
	"fmt"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkauth"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdkjwt"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

func Auth(
	jwt sdkjwt.Jwt,
	auth *sdkauth.Auth,
	domain string,
) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Token")
		sdk.Assert(token != "", sdkcodes.AccessLimited.WithMsg("Token不存在"))

		userClaims, err := jwt.ParseToken(token)
		if err != nil {
			sdk.Assert(false, sdkcodes.AccessLimited.WithMsg(err.Error()))
		}
		sub := fmt.Sprintf("role:%s", userClaims.RoleId.String())
		obj := fmt.Sprintf("%s:%s", domain, ctx.Request.URL.Path)

		ok, err := auth.Enforce(sub, obj)
		if err != nil {
			sdklog.Error(err)
			sdk.Assert(false, sdkcodes.Internal)
		}
		sdk.Assert(ok, sdkcodes.AccessLimited.WithMsg("`%s` 访问受限", ctx.Request.URL.Path))

		if jwt.NeedFlush(userClaims) {
			newToken, err := jwt.FlushToken(userClaims)
			if err != nil {
				sdklog.Error(err)
			} else {
				ctx.Header("Token", newToken)
			}
		}

		SetUserClaims(ctx, userClaims)

		ctx.Next()
	}
}

func GetUserClaims(ctx *gin.Context) (*sdkjwt.UserClaims, bool) {
	v, ok := ctx.Get("__UserClaims__")
	if !ok {
		return nil, false
	}
	userClaims, ok := v.(*sdkjwt.UserClaims)
	if !ok {
		return nil, false
	}
	return userClaims, true
}

func SetUserClaims(ctx *gin.Context, userClaims *sdkjwt.UserClaims) {
	ctx.Set("__UserClaims__", userClaims)
}
