package example

import (
	"context"
	"github.com/iok200/go-ok/log"
	"github.com/iok200/go-ok/micro"
	"google.golang.org/grpc"
	"sync"
	"testing"
	"time"
)

func init() {
	//config.Name = "../ok.properties"
}

func TestServer(t *testing.T) {
	for a := 0; a < 100; a++ {
		createServer()
	}
	log.Infoln("服务创建完成")
	time.Sleep(time.Hour)
}

func TestClient(t *testing.T) {
	cli := micro.NewClient("demoCluster", "demoGroup", "demoService")
	cli.Dial()
	var helloClient HelloClient
	cli.GetConn(func(conn *grpc.ClientConn) {
		helloClient = NewHelloClient(conn)
	})
	count := 20000
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

	for a := 0; a < 100; a++ {
		wg.Add(count)
		beginTime = time.Now()
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
		endTime = time.Now()
		log.Infoln(endTime.Sub(beginTime).Milliseconds())
	}
}

func createServer() *micro.Server {
	var ser *micro.Server
	ser = micro.NewServer("demoCluster", "demoGroup", "demoService")
	ser.Run()
	ser.GetServer(func(s *grpc.Server) {
		impl := new(Impl)
		impl.Addr = ser.GetAddr()
		RegisterHelloServer(s, impl)
	})
	return ser
}
