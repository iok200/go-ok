package server

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	config vo.RegisterInstanceParam
	server *grpc.Server
	mu     sync.Mutex
}

func New(config vo.RegisterInstanceParam) *Server {
	return &Server{config: config}
}

func (this *Server) Run() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server != nil {
		return errors.New("server is runing")
	}
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(this.config.Port)))
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

//func (this *Server) Service() error {
//	this.mu.Lock()
//	defer this.mu.Unlock()
//	if this.server == nil {
//		return errors.New("server is not run")
//	}
//	this.server.RegisterService()
//	return nil
//}
