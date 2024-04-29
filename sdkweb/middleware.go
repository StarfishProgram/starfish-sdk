package sdkweb

import (
	"net/http"

	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

func MWCatch(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			sdklog.Ins().AddCallerSkip(1).Error(err)
			ctx.Abort()
		}
	}()
	ctx.Next()
}

func MWCors(ctx *gin.Context) {
	method := ctx.Request.Method
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Token, CAPTCHA_ID, CAPTCHA_CODE")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
