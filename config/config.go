package config

import (
	"github.com/iok200/go-ok/util"
	"github.com/magiconair/properties"
	"sync"
)

var Name = "ok.properties"
var _config *properties.Properties
var _config_mu sync.Mutex

func Load() *properties.Properties {
	if _config == nil {
		_config_mu.Lock()
		defer _config_mu.Unlock()
		if _config == nil {
			_config = util.GetConfig(Name)
		}
	}
	return _config
}

func Parse(cfg interface{}) {
	if err := Load().Decode(cfg); err != nil {
		panic(err)
	}
}
