package server

import (
	"errors"
	"fmt"
	"github.com/iok200/go-ok/nacosclient"
	"github.com/iok200/go-ok/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	config vo.RegisterInstanceParam
	server *grpc.Server
	mu     sync.Mutex
}

func New(clusterName, groupName, serviceName string) *Server {
	return &Server{config: vo.RegisterInstanceParam{
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: clusterName,
		GroupName:   groupName,
		ServiceName: serviceName,
	}}
}

func (this *Server) Run() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server != nil {
		return errors.New("server is runing")
	}
	nacosClient, err := nacosclient.Load()
	if err != nil {
		return err
	}
	namingClient, err := nacosClient.GetNamingClient()
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}
	localPort, err := util.GetAddrPort(lis.Addr().String())
	if err != nil {
		return err
	}
	localIp, err := util.GetIP()
	if err != nil {
		return err
	}
	this.config.Ip = localIp.String()
	this.config.Port = uint64(localPort)
	grpcServer := grpc.NewServer()
	if ok, err := namingClient.RegisterInstance(this.config); !ok || err != nil {
		if err != nil {
			return err
		}
		return errors.New("service register failed")
	}
	go grpcServer.Serve(lis)
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
	nacosClient, err := nacosclient.Load()
	if err != nil {
		fmt.Println(err)
		return
	}
	namingClient, err := nacosClient.GetNamingClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ephemeral:   true,
		Ip:          this.config.Ip,
		Port:        this.config.Port,
		Cluster:     this.config.ClusterName,
		GroupName:   this.config.GroupName,
		ServiceName: this.config.ServiceName,
	}); !ok || err != nil {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("注销服务失败:%+v", this.config)
	}
}

func (this *Server) GetService(f func(s *grpc.Server)) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server == nil {
		return errors.New("server is not run")
	}
	f(this.server)
	return nil
}
