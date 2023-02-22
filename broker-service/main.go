package main

import (
	"log"
	"net"

	"github.com/gitsuki/finance/broker/api"
	"github.com/gitsuki/finance/broker/pb"
	"github.com/gitsuki/finance/broker/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("[fatal] cannot load config", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("[fatal] cannot create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBrokerServer(grpcServer, server)
	reflection.Register(grpcServer) // allows gRPC client to explore which RPCs are available on the server

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal("[fatal] cannot create listener", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("[fatal] unable to launch gRPC server", err)
	}
}
