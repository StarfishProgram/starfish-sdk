package sdkuuid

import "github.com/StarfishProgram/starfish-sdk/sdktypes"

type UUID interface {
	Id() sdktypes.ID
	Uuid() string
	TimeID() string
}
