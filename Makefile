default:
	@echo "make genrpc : 生成GRPC代码"

genrpc:
	-@rm -rf sdkrpc/*.pb.go
	-@protoc --go_out=sdkrpc \
	--go_opt=paths=source_relative \
	--go-grpc_out=sdkrpc \
	--go-grpc_opt=paths=source_relative \
	--proto_path=sdkrpc \
	sdkrpc/*.proto