package api

type Server struct {
}

func NewServer() (*Server, error) {
	newServer := &Server{}
	return newServer, nil
}
