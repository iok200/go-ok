package example

import (
	"context"
	"fmt"
	"github.com/iok200/go-ok/config"
	"github.com/iok200/go-ok/rpc/client"
	"github.com/iok200/go-ok/rpc/server"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	createServer()
	createServer()
	createServer()
	time.Sleep(time.Hour)
}

func TestClient(t *testing.T) {
	createClient()
	time.Sleep(time.Hour)
}

func createServer() {
	config.SetDefaultConfigPath("../ok.yaml")
	var ser *server.Server
	var err error
	ser = server.New("demoCluster", "demoGroup", "demoService")
	if err = ser.Run(); err != nil {
		fmt.Println(err)
		return
	}
	if err = ser.GetServer(func(s *grpc.Server) {
		impl := new(Impl)
		impl.Addr = ser.GetAddr()
		RegisterHelloServer(s, impl)
	}); err != nil {
		fmt.Println(err)
		return
	}
}

func createClient() {
	config.SetDefaultConfigPath("../ok.yaml")
	cli := client.New("demoCluster", "demoGroup", "demoService")
	if err := cli.Dial(); err != nil {
		fmt.Println(err)
		return
	}
	var helloClient HelloClient
	if err := cli.GetConn(func(conn *grpc.ClientConn) {
		helloClient = NewHelloClient(conn)
	}); err != nil {
		fmt.Println(err)
		return
	}
	for a := 0; a < 10; a++ {
		_, err := helloClient.SayHello(context.Background(), &HelloRequest{Name: "111"})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
