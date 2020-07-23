package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Nacos struct {
		Addr string `yaml:"addr"`
	} `yaml:"nacos"`
}

func Get(fp string) (*Config, error) {
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
