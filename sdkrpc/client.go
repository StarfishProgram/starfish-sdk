package sdkrpc

import (
	reflect "reflect"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

var clientIns map[string]*_Client

func init() {
	clientIns = make(map[string]*_Client)
}

type _Client struct {
	client GRPCServiceClient
}

func InitClient(url string, key ...string) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		sdklog.Ins().AddCallerSkip(1).Panic(err)
	}
	client := NewGRPCServiceClient(conn)
	ins := _Client{client: client}
	if len(key) == 0 {
		clientIns[""] = &ins
	} else {
		clientIns[key[0]] = &ins
	}
}

func Client(key ...string) *_Client {
	if len(key) == 0 {
		return clientIns[""]
	} else {
		return clientIns[key[0]]
	}
}

type CallResult[D protoreflect.ProtoMessage] struct {
	Code *Code
	Data D
}

func Call[P, R protoreflect.ProtoMessage](client *_Client, param P) CallResult[R] {
	var r R
	anyParam, err := anypb.New(param)
	if err != nil {
		return CallResult[R]{
			Code: &Code{
				Code: sdkcodes.Internal.Code(),
				Msg:  err.Error(),
				I18N: sdkcodes.Internal.I18n(),
			},
			Data: r,
		}
	}
	result, err := client.client.Call(sdk.Context(), anyParam)
	if err != nil {
		sdklog.Ins().AddCallerSkip(1).Error(err)
		return CallResult[R]{
			Code: &Code{
				Code: sdkcodes.Internal.Code(),
				Msg:  err.Error(),
				I18N: sdkcodes.Internal.I18n(),
			},
			Data: r,
		}
	}
	if result.Code != nil {
		return CallResult[R]{
			Code: &Code{
				Code: result.Code.Code,
				Msg:  result.Code.Msg,
				I18N: result.Code.I18N,
			},
			Data: r,
		}
	}
	realData := reflect.New(reflect.TypeOf(r).Elem()).Interface().(R)
	if err := result.Data.UnmarshalTo(realData); err != nil {
		return CallResult[R]{
			Code: &Code{
				Code: sdkcodes.Internal.Code(),
				Msg:  err.Error(),
				I18N: sdkcodes.Internal.I18n(),
			},
			Data: r,
		}
	}
	return CallResult[R]{Data: realData}
}
