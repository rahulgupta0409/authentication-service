cd proto
protoc --go_out=. --go-grpc_out=. --proto_path=. user.proto