package api

import (
	"github.com/gitsuki/finance/listener/pb"
	"github.com/gitsuki/finance/listener/util"
)

// Server serves gRPC requests for our broker service
type Server struct {
	pb.UnimplementedBrokerServer
	config util.Config
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config) (*Server, error) {
	newServer := &Server{
		config: config,
	}

	return newServer, nil
}
