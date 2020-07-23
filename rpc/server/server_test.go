package server

import (
	"fmt"
	"github.com/iok200/go-ok/config"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	config.SetDefaultConfigPath("../../ok.yaml")
	server := New("demoCluster", "demoGroup", "demoService")
	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Hour)
}
