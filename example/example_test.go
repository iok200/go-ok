package example

import (
	"fmt"
	"github.com/iok200/go-ok/config"
	"github.com/iok200/go-ok/rpc/server"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	config.SetDefaultConfigPath("../ok.yaml")
	var ser *server.Server
	var err error
	ser = server.New("demoCluster", "demoGroup", "demoService")
	if err = ser.Run(); err != nil {
		fmt.Println(err)
		return
	}
	if err = ser.GetService(func(s *grpc.Server) {
		RegisterHelloServer(s, new(Impl))
	}); err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Hour)
}
