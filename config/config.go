package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var _default_config_path = "ok.yaml"
var _config *Config
var _config_mu sync.Mutex

func Load() (*Config, error) {
	if _config != nil {
		return _config, nil
	}
	_config_mu.Lock()
	defer _config_mu.Unlock()
	if _config != nil {
		return _config, nil
	}
	config, err := New(_default_config_path)
	if err != nil {
		return nil, err
	}
	_config = config
	return _config, nil
}

func SetDefaultConfigPath(p string) {
	_default_config_path = p
}

type Config struct {
	Nacos struct {
		Addr string `yaml:"addr"`
	} `yaml:"nacos"`
}

func New(fp string) (*Config, error) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
