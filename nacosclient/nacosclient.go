package nacosclient

import (
	"github.com/iok200/go-ok/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"strconv"
	"strings"
	"sync"
)

var _client *Client
var _client_mu sync.Mutex

func Load() (*Client, error) {
	if _client != nil {
		return _client, nil
	}
	_client_mu.Lock()
	defer _client_mu.Unlock()
	if _client != nil {
		return _client, nil
	}
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}
	client := New(conf.Nacos.Addr)
	_client = client
	return _client, nil
}

type Client struct {
	addr         string
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
	mu           sync.Mutex
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (this *Client) GetNamingClient() (naming_client.INamingClient, error) {
	if err := this.initNamingClient(); err != nil {
		return nil, err
	}
	return this.namingClient, nil
}
func (this *Client) GetConfigClient() (config_client.IConfigClient, error) {
	if err := this.initConfigClient(); err != nil {
		return nil, err
	}
	return this.configClient, nil
}

func (this *Client) initNamingClient() error {
	if this.namingClient != nil {
		return nil
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.namingClient == nil {
		config, err := this.buildConfig()
		if err != nil {
			return err
		}
		namingClient, err := clients.CreateNamingClient(config)
		if err != nil {
			return err
		}
		this.namingClient = namingClient
	}
	return nil
}

func (this *Client) initConfigClient() error {
	if this.configClient != nil {
		return nil
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.configClient == nil {
		config, err := this.buildConfig()
		if err != nil {
			return err
		}
		configClient, err := clients.CreateConfigClient(config)
		if err != nil {
			return err
		}
		this.configClient = configClient
	}
	return nil
}

func (this *Client) buildConfig() (map[string]interface{}, error) {
	clientConfig := constant.ClientConfig{
		TimeoutMs:            10 * 1000,
		ListenInterval:       30 * 1000,
		BeatInterval:         5 * 1000,
		UpdateThreadNum:      5,
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: true,
		LogDir:               "__NacosCache__/logs",
		CacheDir:             "__NacosCache__/cache",
	}
	var serverConfigs []constant.ServerConfig
	if strings.Index(this.addr, ",") != -1 {
		addrs := strings.Split(this.addr, ",")
		for _, v := range addrs {
			addrSplit := strings.Split(v, ":")
			ip := addrSplit[0]
			port, err := strconv.ParseInt(addrSplit[1], 10, 64)
			if err != nil {
				return nil, err
			}
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				ContextPath: "/nacos",
				IpAddr:      ip,
				Port:        uint64(port),
			})
		}

	} else {
		addrSplit := strings.Split(this.addr, ":")
		ip := addrSplit[0]
		port, err := strconv.ParseInt(addrSplit[1], 10, 64)
		if err != nil {
			return nil, err
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      ip,
			Port:        uint64(port),
		})
	}
	return map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	}, nil
}
