package api

import "github.com/gitsuki/finance/broker/util"

type Server struct {
	config util.Config
}

func NewServer(config util.Config) (*Server, error) {
	newServer := &Server{
		config: config,
	}

	return newServer, nil
}
