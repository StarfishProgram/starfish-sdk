package sdkrpc

import (
	context "context"
	"net"
	"os"
	reflect "reflect"

	"github.com/StarfishProgram/starfish-sdk/sdk"
	"github.com/StarfishProgram/starfish-sdk/sdkcodes"
	"github.com/StarfishProgram/starfish-sdk/sdklog"
	grpc "google.golang.org/grpc"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type _Server struct {
	UnimplementedGRPCServiceServer
	calls map[string]func(*anypb.Any) *anypb.Any
}

// ServerRegisterCall 注册服务
func ServerRegisterCall[P, R protoreflect.ProtoMessage](server *_Server, call func(param P) R) {
	var p P
	paramAny, err := anypb.New(p)
	if err != nil {
		sdklog.Ins().AddCallerSkip(1).Panic(err)
	}
	pt := reflect.TypeOf(p).Elem()
	server.calls[paramAny.TypeUrl] = func(param *anypb.Any) *anypb.Any {
		realParam := reflect.New(pt).Interface().(P)
		err := param.UnmarshalTo(realParam)
		sdk.CheckError(err)
		callResult := call(realParam)
		resultData, err := anypb.New(callResult)
		sdk.CheckError(err)
		return resultData
	}
	sdklog.Ins().AddCallerSkip(1).Info("RPC服务注册 :", paramAny.TypeUrl)
}

func (s *_Server) Call(ctx context.Context, param *anypb.Any) (result *Result, err error) {
	result = &Result{Code: nil, Data: nil}
	call, ok := s.calls[param.TypeUrl]
	if !ok {
		result.Code = &Code{
			Code: sdkcodes.RequestNotFound.Code(),
			Msg:  sdkcodes.RequestNotFound.Msg(),
			I18N: sdkcodes.RequestNotFound.I18n(),
		}
		return
	}
	defer func() {
		if err := recover(); err != nil {
			result.Data = nil
			if code, ok := err.(sdkcodes.Code); ok {
				result.Code = &Code{
					Code: code.Code(),
					Msg:  code.Msg(),
					I18N: code.I18n(),
				}
				sdklog.Ins().AddCallerSkip(3).Warn(code)
				return
			}
			sdklog.Ins().AddCallerSkip(2).Error(err)
			result.Code = &Code{
				Code: sdkcodes.Internal.Code(),
				Msg:  sdkcodes.Internal.Msg(),
				I18N: sdkcodes.Internal.I18n(),
			}
		}
	}()
	result.Data = call(param)
	return
}

func InitServer(listener string) (*_Server, chan os.Signal) {
	lis, err := net.Listen("tcp", listener)
	if err != nil {
		sdklog.Ins().Panicf("GRPC服务创建失败 : %s", err.Error())
	}
	server := _Server{
		calls: map[string]func(*anypb.Any) *anypb.Any{},
	}
	rpcServer := grpc.NewServer()
	RegisterGRPCServiceServer(rpcServer, &server)
	ch := make(chan os.Signal, 1)
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			sdklog.Ins().Error("GRPC服务运行异常", err)
		}
		sdklog.Ins().Info("GRPC服务已停止")
		close(ch)
	}()
	go func() {
		<-ch
		rpcServer.Stop()
	}()
	return &server, ch
}
