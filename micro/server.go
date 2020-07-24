package micro

import (
	"errors"
	"github.com/iok200/go-ok/log"
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

func NewServer(clusterName, groupName, serviceName string) *Server {
	return &Server{clusterName: clusterName, groupName: groupName, serviceName: serviceName}
}

func (this *Server) Run() {
	if this.server == nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		if this.server == nil {
			lis, err := net.Listen("tcp", ":0")
			if err != nil {
				panic(err)
			}
			localPort, err := util.GetAddrPort(lis.Addr().String())
			if err != nil {
				panic(err)
			}
			localIp, err := util.GetIP()
			if err != nil {
				panic(err)
			}
			this.ip = localIp.String()
			this.port = localPort
			grpcServer := grpc.NewServer()
			if ok, err := nacosclient.Load().GetNamingClient().RegisterInstance(vo.RegisterInstanceParam{
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
					panic(err)
				}
			}
			go grpcServer.Serve(lis)
			this.server = grpcServer
		}
	}
}

func (this *Server) Stop() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server == nil {
		return
	}
	this.server.Stop()
	this.server = nil
	if ok, err := nacosclient.Load().GetNamingClient().DeregisterInstance(vo.DeregisterInstanceParam{
		Ephemeral:   true,
		Ip:          this.ip,
		Port:        uint64(this.port),
		Cluster:     this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
	}); !ok || err != nil {
		if err != nil {
			log.Infoln(err)
			return
		}
		log.Infof("注销服务失败:%+v\n", this)
	}
}

func (this *Server) GetServer(f func(s *grpc.Server)) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.server == nil {
		panic(errors.New("server is not run"))
	}
	f(this.server)
}

func (this *Server) GetAddr() string {
	return this.ip + ":" + strconv.Itoa(this.port)
}
