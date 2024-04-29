package sdk

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// WebConfig WEB配置
type WebConfig struct {
	Listen string `toml:"listen" ` // 监听地址
}

func WebMWCatch(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			Log().AddCallerSkip(1).Error(err)
			ctx.Abort()
		}
	}()
	ctx.Next()
}

func WebMWCors(ctx *gin.Context) {
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

func WebMWRequestParam[T any](call func(*gin.Context, *T)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var p T
		err := ctx.ShouldBind(&p)
		CheckError(err, CodeParamInvalid.WithError(err))
		call(ctx, &p)
	}
}

func WebRpsJson(ctx *gin.Context, code ICode, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Code(),
		"msg":  code.Msg(),
		"i18n": code.I18n(),
		"data": data,
	})
}

// InitWeb Web初始化
func InitWeb(config *WebConfig, routers func(eng *gin.Engine)) chan os.Signal {
	gin.DisableConsoleColor()
	server := gin.New()
	server.Use(WebMWCatch, WebMWCors)
	server.NoRoute(func(ctx *gin.Context) {
		WebRpsJson(ctx, CodeServerError, nil)
	})
	if routers != nil {
		routers(server)
	}
	httpServer := http.Server{Addr: config.Listen, Handler: server}
	ch := make(chan os.Signal, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Log().Error("Web运行异常", err)
		}
		Log().Info("WEB服务已停止")
		close(ch)
	}()
	go func() {
		<-ch
		httpServer.Shutdown(Context())
	}()
	return ch
}
