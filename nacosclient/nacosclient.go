package nacosclient

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"strconv"
	"strings"
	"sync"
)

var _addr = "127.0.0.1:8848"
var _namingClient naming_client.INamingClient
var _namingClientMu sync.Mutex
var _configClient config_client.IConfigClient
var _configClientMu sync.Mutex

func SetAddr(addr string) {
	_addr = addr
}

func LoadNamingClient() (naming_client.INamingClient, error) {
	if _namingClient == nil {
		_namingClientMu.Lock()
		if _namingClient == nil {
			conf, err := buildConfig()
			if err != nil {
				return nil, err
			}
			namingClient, err := clients.CreateNamingClient(conf)
			if err != nil {
				return nil, err
			}
			_namingClient = namingClient
		}
		_namingClientMu.Unlock()
	}
	return _namingClient, nil
}

func LoadConfigClient() (config_client.IConfigClient, error) {
	if _configClient == nil {
		_configClientMu.Lock()
		if _configClient == nil {
			conf, err := buildConfig()
			if err != nil {
				return nil, err
			}
			configClient, err := clients.CreateConfigClient(conf)
			if err != nil {
				return nil, err
			}
			_configClient = configClient
		}
		_configClientMu.Unlock()
	}
	return _configClient, nil
}

func buildConfig() (map[string]interface{}, error) {
	clientConfig := constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 30 * 1000,
		BeatInterval:   5 * 1000,
		LogDir:         "nacos/logs",
		CacheDir:       "nacos/cache",
	}
	var serverConfigs []constant.ServerConfig
	if strings.Index(_addr, ",") != -1 {
		addrs := strings.Split(_addr, ",")
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
		addrSplit := strings.Split(_addr, ":")
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
