package sdkwebmiddleware

import (
	"net/http"

	"github.com/StarfishProgram/starfish-sdk/sdkauth"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdkjwt"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

func Auth(jwt sdkjwt.Jwt, auth *sdkauth.Auth, domain string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Token")
		if token == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": sdkcodes.AccessLimited.Code(),
				"msg":  "token not found",
				"i18n": sdkcodes.AccessLimited.I18n(),
				"data": nil,
			})
			ctx.Abort()
			return
		}
		userClaims, err := jwt.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": sdkcodes.AccessLimited.Code(),
				"msg":  err.Error(),
				"i18n": sdkcodes.AccessLimited.I18n(),
				"data": nil,
			})
			ctx.Abort()
			return
		}
		roleId := userClaims.RoleId
		url := ctx.Request.URL.Path
		method := ctx.Request.Method

		ok, err := auth.Enforce(roleId, domain, url, method)
		if err != nil {
			sdklog.Error(err)
			ctx.JSON(http.StatusOK, gin.H{
				"code": sdkcodes.Internal.Code(),
				"msg":  sdkcodes.Internal.Msg(),
				"i18n": sdkcodes.Internal.I18n(),
				"data": nil,
			})
			ctx.Abort()
			return
		}
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"code": sdkcodes.AccessLimited.Code(),
				"msg":  sdkcodes.AccessLimited.Msg(),
				"i18n": sdkcodes.AccessLimited.I18n(),
				"data": nil,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
