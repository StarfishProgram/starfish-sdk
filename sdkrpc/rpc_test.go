package sdkrpc

import (
	"testing"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
)

func TestServer(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	server, ch := InitServer("0.0.0.0:12345")
	ServerRegisterCall(server, demoFn)
	signal := sdk.NewSignal()
	signal.Add(ch)
	signal.Waiting()
}

func demoFn(*Code) *Code {
	sdk.Check(false)
	return &Code{
		Code: 123,
		Msg:  "456",
		I18N: "789",
	}
}

func TestClient(t *testing.T) {
	sdklog.Init(&sdklog.Config{Level: "info"})
	InitClient("0.0.0.0:12345")
	result := Call[*Code, *Code](Client(), &Code{
		Code: 1,
		Msg:  "2",
		I18N: "3",
	})
	sdklog.Info(result)
}
