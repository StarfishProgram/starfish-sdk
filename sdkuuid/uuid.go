package sdkuuid

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdktypes"
	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
)

type _UUID struct {
	ins *snowflake.Node
}

func (u *_UUID) Id() sdktypes.ID {
	return sdktypes.ID(u.ins.Generate().Int64())
}

func (u *_UUID) Uuid() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func (u *_UUID) TimeID() string {
	now := time.Now().Format("20060102150405")
	n := rand.Int63n(1000000000000000000)
	return fmt.Sprintf("%s%018d", now, n)
}

var ins map[string]UUID

func init() {
	ins = make(map[string]UUID)
}

// Init 初始化UUID
func Init(nodeSeq int64, key ...string) {
	sdk.Check(nodeSeq < 1024)
	node, err := snowflake.NewNode(nodeSeq)
	sdk.CheckError(err)
	_ins := &_UUID{node}
	if len(key) == 0 {
		ins[""] = _ins
	} else {
		ins[key[0]] = _ins
	}
}

// Ins 获取UUID
func Ins(key ...string) UUID {
	if len(key) == 0 {
		return ins[""]
	} else {
		return ins[key[0]]
	}
}
