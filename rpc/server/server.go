package server

import (
	"errors"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	addr   string
	server *grpc.Server
	mu     sync.Mutex
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (this *Server) Run() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server != nil {
		return errors.New("server is runing")
	}
	lis, err := net.Listen("tcp", this.addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	this.server = grpcServer
	return nil
}

func (this *Server) Stop() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server == nil {
		return
	}
	this.server.Stop()
}
