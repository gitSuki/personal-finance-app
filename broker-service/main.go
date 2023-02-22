package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gitsuki/finance/broker/api"
	"github.com/gitsuki/finance/broker/pb"
	"github.com/gitsuki/finance/broker/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("[fatal] unable to load config", err)
	}

	go startgRPCServer(config) // runs on seperate goroutine to handle HTTP and gRPC requests concurrently
	startHTTPProxyServer(config)
}

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

func startHTTPProxyServer(config util.Config) {
	brokerServer, err := api.NewServer(config)
	if err != nil {
		log.Fatal("[fatal] unable to create broker server", err)
	}

	gRPCGatewayMux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // isn't called until the current function has returned

	err = pb.RegisterBrokerHandlerServer(ctx, gRPCGatewayMux, brokerServer)
	if err != nil {
		log.Fatal("[fatal] unable to register gateway handler server")
	}

	serveMux := http.NewServeMux()       // handles HTTP requests
	serveMux.Handle("/", gRPCGatewayMux) // forwards all requests to gRPC gateway mux

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("[fatal] unable to create listener", err)
	}

	log.Printf("start gRPC Gateway HTTP Proxy server at %s", listener.Addr().String())
	err = http.Serve(listener, serveMux)
	if err != nil {
		log.Fatal("[fatal] unable to launch gRPC Gateway HTTP Proxy server", err)
	}
}
