package starfish_sdk

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	uuid "github.com/satori/go.uuid"
)

type IUUID interface {
	// ID ID
	ID() ID
	// UUID UUID
	UUID() string
	// TimeID 时间ID
	TimeID() string
}

type _uuid struct {
	ins *snowflake.Node
}

func (u *_uuid) ID() ID {
	return ID(u.ins.Generate().Int64())
}

func (u *_uuid) UUID() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}

func (u *_uuid) TimeID() string {
	now := time.Now().Format("20060102150405")
	n := rand.Int63n(1000000000000000000)
	return fmt.Sprintf("%s%018d", now, n)
}

var uuidIns map[string]IUUID

func init() {
	uuidIns = make(map[string]IUUID)
}

// InitUUID 初始化UUID
func InitUUID(nodeSeq int64, key ...string) {
	Check(nodeSeq < 1024)
	node, err := snowflake.NewNode(nodeSeq)
	CheckError(err)
	ins := &_uuid{node}
	if len(key) == 0 {
		uuidIns[""] = ins
	} else {
		uuidIns[key[0]] = ins
	}
}

// UUID 获取UUID
func UUID(key ...string) IUUID {
	if len(key) == 0 {
		return uuidIns[""]
	} else {
		return uuidIns[key[0]]
	}
}
