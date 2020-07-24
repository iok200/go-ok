package server

import (
	"github.com/iok200/go-ok/config"
	"github.com/iok200/go-ok/log"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	config.SetDefaultConfigPath("../../ok.yaml")
	server := New("demoCluster", "demoGroup", "demoService")
	if err := server.Run(); err != nil {
		log.Infoln(err)
		return
	}
	time.Sleep(time.Hour * 24)
}
