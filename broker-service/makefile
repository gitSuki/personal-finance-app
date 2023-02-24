# builds docker image
docker:
	docker build -t broker-service:1.1 .

# launches rabbitmq docker container
rabbitmq:
	docker run --hostname my-rabbit --name rabbitmq3 -p 5672:5672 -d rabbitmq:3.11-alpine

# launches the server
server:
	go run main.go

# generate golang code for gRPC services and methods defined in the proto file
proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

# client for sending test requests to the gRPC server
evans:
	evans --host localhost --port 9000 -r repl

.PHONY: rabbitmiq server proto evans