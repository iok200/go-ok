package example

import (
	"context"
	"github.com/iok200/go-ok/config"
	"github.com/iok200/go-ok/log"
	"github.com/iok200/go-ok/rpc/client"
	"github.com/iok200/go-ok/rpc/server"
	"google.golang.org/grpc"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	for a := 0; a < 200; a++ {
		createServer()
	}
	log.Infoln("服务创建完成")
	time.Sleep(time.Hour)
}

func TestClient(t *testing.T) {
	config.SetDefaultConfigPath("../ok.yaml")
	cli := client.New("demoCluster", "demoGroup", "demoService")
	if err := cli.Dial(); err != nil {
		log.Infoln(err)
		return
	}
	var helloClient HelloClient
	if err := cli.GetConn(func(conn *grpc.ClientConn) {
		helloClient = NewHelloClient(conn)
	}); err != nil {
		log.Infoln(err)
		return
	}
	count := 10000
	var wg sync.WaitGroup
	wg.Add(count)
	beginTime := time.Now()
	for a := 0; a < count; a++ {
		go func() {
			_, err := helloClient.SayHello(context.Background(), &HelloRequest{Name: "111"})
			wg.Done()
			if err != nil {
				log.Infoln(err)
				return
			}
		}()
	}
	wg.Wait()
	endTime := time.Now()
	log.Infoln(endTime.Sub(beginTime).Milliseconds())
}

func createServer() *server.Server {
	config.SetDefaultConfigPath("../ok.yaml")
	var ser *server.Server
	var err error
	ser = server.New("demoCluster", "demoGroup", "demoService")
	if err = ser.Run(); err != nil {
		log.Infoln(err)
		return nil
	}
	if err = ser.GetServer(func(s *grpc.Server) {
		impl := new(Impl)
		impl.Addr = ser.GetAddr()
		RegisterHelloServer(s, impl)
	}); err != nil {
		log.Infoln(err)
		return nil
	}
	return ser
}
