package sdk

import (
	"errors"
	"testing"
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdklog"
)

func TestXxx(t *testing.T) {
	sdklog.Init(&sdklog.Config{
		Level: "info",
	})
	Go(func() {
		panic(errors.New("123"))
	})

	time.Sleep(time.Hour)
}
