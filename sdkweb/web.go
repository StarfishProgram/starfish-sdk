package sdkweb

import (
	"net/http"
	"os"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	sdkwebmiddleware "github.com/StarfishProgram/starfish-sdk/sdkweb/middleware"
	"github.com/gin-gonic/gin"
)

// Config WEB配置
type Config struct {
	Listen string `toml:"listen"` // 监听地址
}

func ResponseData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, sdk.AnyMap{
		"code": sdkcodes.OK.Code(),
		"msg":  sdkcodes.OK.Msg(),
		"i18n": sdkcodes.OK.I18n(),
		"data": data,
	})
}

func ResponseError(ctx *gin.Context, code sdkcodes.Code) {
	ctx.JSON(http.StatusOK, sdk.AnyMap{
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
	server.Use(sdkwebmiddleware.Catch, sdkwebmiddleware.Cors)
	server.NoRoute(func(ctx *gin.Context) {
		ResponseError(ctx, sdkcodes.RequestNotFound)
	})
	if routers != nil {
		routers(server)
	}
	httpServer := http.Server{Addr: config.Listen, Handler: server}
	ch := make(chan os.Signal, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sdklog.Error("Web运行异常", err)
		}
		sdklog.Info("WEB服务已停止")
		close(ch)
	}()
	go func() {
		<-ch
		httpServer.Shutdown(sdk.Context())
	}()
	return ch
}
