# launches the server
server:
	go run main.go

# generate golang code for gRPC services and methods defined in the proto file
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

# client for sending test requests to the gRPC server
evans:
	evans --host localhost --port 8000 -r repl

.PHONY: proto evans