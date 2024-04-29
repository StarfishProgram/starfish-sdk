# 海星SDK

---

### 环境安装

```bash
# GRPC安装
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# 环境配置
export PATH="$PATH:$(go env GOPATH)/bin"

# 生成GRPC代码
protoc --go_out=sdkrpc \
    --go_opt=paths=source_relative \
    --go-grpc_out=sdkrpc \
    --go-grpc_opt=paths=source_relative \
    --proto_path=sdkrpc \
    sdkrpc/*.proto
```
