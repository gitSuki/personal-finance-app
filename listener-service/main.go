package main

import (
	"log"
	"net"

	"github.com/gitsuki/finance/listener/api"
	"github.com/gitsuki/finance/listener/event"
	"github.com/gitsuki/finance/listener/pb"
	"github.com/gitsuki/finance/listener/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("[fatal] unable to load config", err)
	}

	// run on seperate goroutines to handle HTTP and gRPC requests and listen to rabbitMQ responses concurrently
	go event.RecieveRequests(config)
	startgRPCServer(config)
}

// startgRPCServer starts the gRPC server responsible for handling protocol buffer format requests
func startgRPCServer(config util.Config) {
	brokerServer, err := api.NewServer(config)
	if err != nil {
		log.Fatal("[fatal] unable to create broker server", err)
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterBrokerServer(gRPCServer, brokerServer)
	reflection.Register(gRPCServer) // allows gRPC client to explore which RPCs are available on the server

	listener, err := net.Listen("tcp", config.ProtobufServerAddress)
	if err != nil {
		log.Fatal("[fatal] unable to create listener", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = gRPCServer.Serve(listener)
	if err != nil {
		log.Fatal("[fatal] unable to launch gRPC server", err)
	}
}
