package sdkweb

import (
	"testing"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	"github.com/gin-gonic/gin"
)

func TestWeb(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	ch := Init(&Config{Listen: "0.0.0.0:12345"}, func(eng *gin.Engine) {
		eng.GET("/", func(ctx *gin.Context) {
			sdk.Assert(false)
		})
	})
	signal := sdk.NewSignal()
	signal.Add(ch)
	signal.Waiting()
}
