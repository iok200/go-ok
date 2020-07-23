package client

import (
	"errors"
	"fmt"
	"github.com/iok200/go-ok/nacosclient"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
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
	clusterName := serviceInfo[0]
	groupName := serviceInfo[1]
	serviceName := serviceInfo[2]
	nacosClient, err := nacosclient.Load()
	if err != nil {
		return nil, err
	}
	namingClient, err := nacosClient.GetNamingClient()
	if err != nil {
		return nil, err
	}
	r := &nacosResovler{
		clusterName:  clusterName,
		groupName:    groupName,
		serviceName:  serviceName,
		namingClient: namingClient,
		cc:           cc,
		opts:         opts,
	}
	r.fetch()
	if err := r.subscribe(); err != nil {
		return nil, err
	}
	return r, nil
}

type nacosResovler struct {
	clusterName  string
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
	instances, err := this.namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"a"},
		HealthyOnly: true,
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

func (this *nacosResovler) subscribe() error {
	return this.namingClient.Subscribe(&vo.SubscribeParam{
		Clusters:    []string{this.clusterName},
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("服务发现:%+v", services)
			addrs := make([]resolver.Address, len(services))
			for i, s := range services {
				addrs[i] = resolver.Address{
					Addr:       s.Ip + ":" + strconv.Itoa(int(s.Port)),
					ServerName: s.ServiceName,
				}
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
	})
}

func (this *nacosResovler) unsubscribe() {
	_ = this.namingClient.Unsubscribe(&vo.SubscribeParam{
		Clusters:    []string{this.clusterName},
		GroupName:   this.groupName,
		ServiceName: this.serviceName,
	})
}

func (this *nacosResovler) ResolveNow(options resolver.ResolveNowOptions) {
}

type Client struct {
	clusterName string
	groupName   string
	serviceName string
	conn        *grpc.ClientConn
	mu          sync.Mutex
}

func New(clusterName, groupName, serviceName string) *Client {
	return &Client{clusterName: clusterName, groupName: groupName, serviceName: serviceName}
}

func (this *Client) Dial() error {
	if this.conn != nil {
		return nil
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn != nil {
		return nil
	}
	conn, err := grpc.Dial("nacos:///"+this.clusterName+"|"+this.groupName+"|"+this.serviceName, grpc.WithInsecure())
	if err != nil {
		return err
	}
	this.conn = conn
	return nil
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
	this.conn.Close()
	this.conn = nil
}

func (this *Client) GetConn(f func(conn *grpc.ClientConn)) error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.conn == nil {
		return errors.New("client is not connection")
	}
	f(this.conn)
	return nil
}
