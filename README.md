
# AUTH
Authorize between different agent nodes

# Prerequisites
- Ubuntu 22.04.4 LTS
- Golang 1.22.3

# Installation
```
cd <path-to-node-gateway-go>
go env -w GOPROXY='https://goproxy.cn,direct'
go install
```

# Build
```
# GOOS: darwin、freebsd、linux、windows
# GOARCH: 386、amd64、arm

# linux x86-64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dorylus-auth-gateway .
```

# Run
```
go run .
```

# Run Test
```
go clean -testcache

go test -v -run ^TestSendAuthRequest$ hajime/node-auth
go test -v -run ^TestSendLLMRequest$ hajime/node-auth
go test -v -run ^TestSendKbSearchRequest$ hajime/node-auth
go test -v -run ^TestGatewaySendLLMRequest$ hajime/node-auth
go test -v -run ^TestGatewaySendKbSearchRequest$ hajime/node-auth
go test -v -run ^TestCheckNodeValidation$ hajime/node-auth
go test -v -run ^TestLoadNodeInfo$ hajime/node-auth
```