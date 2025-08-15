$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o rest ./cmd/rest
go build -o grpc ./cmd/grpc

