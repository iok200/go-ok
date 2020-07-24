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

type Config struct {
	Addr string `properties:"nacos.addr,default=127.0.0.1:8848"`
}

var _client *Client
var _client_mu sync.Mutex

func Load() *Client {
	if _client == nil {
		_client_mu.Lock()
		defer _client_mu.Unlock()
		if _client == nil {
			var cfg *Config
			config.Parse(cfg)
			client := initNacos(cfg)
			_client = client
		}
	}
	return _client
}

type Client struct {
	config       *Config
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
	mu           sync.Mutex
}

func initNacos(cfg *Config) *Client {
	return &Client{config: cfg}
}

func (this *Client) GetNamingClient() naming_client.INamingClient {
	this.initNamingClient()
	return this.namingClient
}
func (this *Client) GetConfigClient() config_client.IConfigClient {
	this.initConfigClient()
	return this.configClient
}

func (this *Client) initNamingClient() {
	if this.namingClient == nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		if this.namingClient == nil {
			cfg := this.buildConfig()
			namingClient, err := clients.CreateNamingClient(cfg)
			if err != nil {
				panic(err)
			}
			this.namingClient = namingClient
		}
	}
}

func (this *Client) initConfigClient() {
	if this.configClient == nil {
		this.mu.Lock()
		defer this.mu.Unlock()
		if this.configClient == nil {
			cfg := this.buildConfig()
			configClient, err := clients.CreateConfigClient(cfg)
			if err != nil {
				panic(err)
			}
			this.configClient = configClient
		}
	}
}

func (this *Client) buildConfig() map[string]interface{} {
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
	if strings.Index(this.config.Addr, ",") != -1 {
		addrs := strings.Split(this.config.Addr, ",")
		for _, v := range addrs {
			addrSplit := strings.Split(v, ":")
			ip := addrSplit[0]
			port, err := strconv.ParseInt(addrSplit[1], 10, 64)
			if err != nil {
				panic(err)
			}
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				ContextPath: "/nacos",
				IpAddr:      ip,
				Port:        uint64(port),
			})
		}

	} else {
		addrSplit := strings.Split(this.config.Addr, ":")
		ip := addrSplit[0]
		port, err := strconv.ParseInt(addrSplit[1], 10, 64)
		if err != nil {
			panic(err)
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
	}
}
