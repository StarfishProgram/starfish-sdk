package sdkrpc

import (
	"testing"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
)

func TestServer(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	server, ch := InitServer("0.0.0.0:12345")
	ServerRegisterCall(server, func(*Code) *Code {
		return &Code{
			Code: 123,
			Msg:  "456",
			I18N: "789",
		}
	})
	signal := sdk.NewSignal()
	signal.Add(ch)
	signal.Waiting()
}

func TestClient(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	InitClient("0.0.0.0:12345")
	result := Call[*Code, *Code](Client(), &Code{
		Code: 0,
		Msg:  "",
		I18N: "",
	})
	sdklog.Ins().Info(result)
}
