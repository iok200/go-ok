package server

import (
	"errors"
	"fmt"
	"github.com/iok200/go-ok/nacosclient"
	"github.com/iok200/go-ok/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	clusterName string
	groupName   string
	serviceName string
	ip          string
	port        int
	server      *grpc.Server
	mu          sync.Mutex
}

func New(clusterName, groupName, serviceName string) *Server {
	return &Server{clusterName: clusterName, groupName: groupName, serviceName: serviceName}
}

func (this *Server) Run() error {
	if this.server != nil {
		return nil
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server != nil {
		return nil
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
	this.ip = localIp.String()
	this.port = localPort
	grpcServer := grpc.NewServer()
	if ok, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
		Ip:          this.ip,
		Port:        uint64(this.port),
	}); !ok || err != nil {
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
	this.server = nil
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
		Ip:          this.ip,
		Port:        uint64(this.port),
		Cluster:     this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
	}); !ok || err != nil {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("注销服务失败:%+v\n", this)
	}
}

func (this *Server) GetServer(f func(s *grpc.Server)) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server == nil {
		return errors.New("server is not run")
	}
	f(this.server)
	return nil
}

func (this *Server) GetAddr() string {
	return this.ip + ":" + strconv.Itoa(this.port)
}
