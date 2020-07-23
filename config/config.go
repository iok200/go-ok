package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

var _defaultFilePath = "ok.yaml"

func SetFilePath(path string) {
	_defaultFilePath = path
}

type Config struct {
	Nacos struct {
		Addr string `yaml:"addr"`
	} `yaml:"nacos"`
}

var _conf *Config
var _confMu sync.Mutex

func Load() (*Config, error) {
	if _conf == nil {
		_confMu.Lock()
		if _conf == nil {
			err := initConf()
			if err != nil {
				return nil, err
			}
		}
		_confMu.Unlock()
	}
	return _conf, nil
}

func initConf() error {
	fileData, err := ioutil.ReadFile(_defaultFilePath)
	if err != nil {
		return err
	}
	conf := new(Config)
	err = yaml.Unmarshal(fileData, conf)
	if err != nil {
		return err
	}
	_conf = conf
	return nil
}
