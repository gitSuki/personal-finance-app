package api

import (
	"context"

	"github.com/gitsuki/finance/broker/event"
	"github.com/gitsuki/finance/broker/pb"
)

func (server *Server) HelloWorld(ctx context.Context, req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	rsp1 := event.SendRequest(server.config, req.GetRequest())

	rsp2 := &pb.HelloWorldResponse{
		Response: rsp1,
	}

	return rsp2, nil
}
