#PROTOC="${HOME}/protoc_bin/bin/protoc"
PROTOC="protoc"
all: proto build test

build: ./proto/sum_grpc.pb.go ./proto/sum.pb.go
	go build

test:
	GOMAXPROCS=8 go test

./proto/sum_grpc.pb.go: ../sum.proto
	${PROTOC}  -I.. --go-grpc_out=./proto --go_out=./proto ../sum.proto

./proto/sum.pb.go: ../sum.proto
	${PROTOC}  -I.. --go-grpc_out=./proto --go_out=./proto ../sum.proto

setup_protoc_exts:
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
