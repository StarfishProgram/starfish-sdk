package sdkwebmiddleware

import (
	"net/http"

	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

func Catch(ctx *gin.Context) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if code, ok := err.(sdkcodes.Code); ok {
			sdklog.AddCallerSkip(3).Warn(code)
			ctx.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Msg(),
				"i18n": code.I18n(),
				"data": nil,
			})
			ctx.Abort()
			return
		}
		sdklog.AddCallerSkip(2).Error(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": sdkcodes.Internal.Code(),
			"msg":  sdkcodes.Internal.Msg(),
			"i18n": sdkcodes.Internal.I18n(),
			"data": nil,
		})
		ctx.Abort()
	}()
	ctx.Next()
}
