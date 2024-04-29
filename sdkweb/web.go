package sdkweb

import (
	"net/http"
	"os"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

// Config WEB配置
type Config struct {
	Listen string `toml:"listen"` // 监听地址
}

type MapData map[string]any

func SuccessResponse(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": sdkcodes.OK.Code(),
		"msg":  sdkcodes.OK.Msg(),
		"i18n": sdkcodes.OK.I18n(),
		"data": data,
	})
}

func ErrorResponse(ctx *gin.Context, code sdkcodes.Code) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code.Code(),
		"msg":  code.Msg(),
		"i18n": code.I18n(),
		"data": nil,
	})
}

// Init Web初始化
func Init(config *Config, routers func(eng *gin.Engine)) chan os.Signal {
	gin.DisableConsoleColor()
	server := gin.New()
	server.Use(MWCatch, MWCors)
	server.NoRoute(func(ctx *gin.Context) {
		ErrorResponse(ctx, sdkcodes.RequestNotFound)
	})
	if routers != nil {
		routers(server)
	}
	httpServer := http.Server{Addr: config.Listen, Handler: server}
	ch := make(chan os.Signal, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sdklog.Ins().Error("Web运行异常", err)
		}
		sdklog.Ins().Info("WEB服务已停止")
		close(ch)
	}()
	go func() {
		<-ch
		httpServer.Shutdown(sdk.Context())
	}()
	return ch
}
