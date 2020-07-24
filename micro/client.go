package micro

import (
	"errors"
	"fmt"
	"github.com/iok200/go-ok/log"
	"github.com/iok200/go-ok/nacosclient"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"strconv"
	"strings"
	"sync"
)

func init() {
	resolver.Register(&nacosBuilder{})
}

type nacosBuilder struct {
}

func (this *nacosBuilder) Scheme() string {
	return "nacos"
}

func (this *nacosBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	serviceInfo := strings.Split(target.Endpoint, "|")
	clusterName := strings.Split(serviceInfo[0], "***")
	groupName := serviceInfo[1]
	serviceName := serviceInfo[2]
	r := &nacosResovler{
		clusterName:  clusterName,
		groupName:    groupName,
		serviceName:  serviceName,
		namingClient: nacosclient.Load().GetNamingClient(),
		cc:           cc,
		opts:         opts,
	}
	r.fetch()
	r.subscribe()
	return r, nil
}

type nacosResovler struct {
	clusterName  []string
	groupName    string
	serviceName  string
	cc           resolver.ClientConn
	opts         resolver.BuildOptions
	namingClient naming_client.INamingClient
}

func (this *nacosResovler) Close() {
	this.unsubscribe()
}

func (this *nacosResovler) fetch() {
	instances, err := this.namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		Clusters:    this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
	})
	var serviceConfig *serviceconfig.ParseResult
	addrs := make([]resolver.Address, len(instances))
	if err != nil || len(instances) == 0 {
		serviceConfig = &serviceconfig.ParseResult{Err: errors.New("service is not available")}
	} else {
		for i, s := range instances {
			addrs[i] = resolver.Address{
				Addr:       s.Ip + ":" + strconv.Itoa(int(s.Port)),
				ServerName: s.ServiceName,
			}
		}
	}
	this.cc.UpdateState(resolver.State{
		Addresses:     addrs,
		ServiceConfig: serviceConfig,
	})
}

func (this *nacosResovler) subscribe() {
	if err := this.namingClient.Subscribe(&vo.SubscribeParam{
		Clusters:    this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			addrs := make([]resolver.Address, len(services))
			if len(services) != 0 {
				log.Debugln("----------客户端服务发现----------")
			}
			for i, s := range services {
				addrs[i] = resolver.Address{
					Addr:       s.Ip + ":" + strconv.Itoa(int(s.Port)),
					ServerName: s.ServiceName,
				}
				log.Debugln(s.Ip + ":" + strconv.Itoa(int(s.Port)))
			}
			if len(services) != 0 {
				log.Debugln("----------客户端服务发现----------")
			}
			var serviceConfig *serviceconfig.ParseResult
			if len(services) == 0 {
				serviceConfig = &serviceconfig.ParseResult{Err: errors.New("service is not available")}
			}
			this.cc.UpdateState(resolver.State{
				Addresses:     addrs,
				ServiceConfig: serviceConfig,
			})
		},
	}); err != nil {
		panic(err)
	}
}

func (this *nacosResovler) unsubscribe() {
	_ = this.namingClient.Unsubscribe(&vo.SubscribeParam{
		Clusters:    this.clusterName,
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
	})
}

func (this *nacosResovler) ResolveNow(options resolver.ResolveNowOptions) {
}

type Client struct {
	clusterName []string
	groupName   string
	serviceName string
	conn        *grpc.ClientConn
	mu          sync.Mutex
}

func NewClient(clusterName []string, groupName, serviceName string) *Client {
	if clusterName == nil {
		clusterName = []string{}
	}
	return &Client{clusterName: clusterName, groupName: groupName, serviceName: serviceName}
}

func (this *Client) Dial() {
	if this.conn != nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn != nil {
		return
	}
	conn, err := grpc.Dial("nacos:///"+strings.Join(this.clusterName, "***")+"|"+this.groupName+"|"+this.serviceName, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		panic(err)
	}
	this.conn = conn
}

func (this *Client) Close() {
	if this.conn == nil {
		return
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn == nil {
		return
	}
	_ = this.conn.Close()
	this.conn = nil
}

func (this *Client) GetConn(f func(conn *grpc.ClientConn)) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn == nil {
		panic(errors.New("client is not connection"))
	}
	f(this.conn)
}
