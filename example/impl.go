package example

import (
	context "context"
	"fmt"
)

type Impl struct {
	Addr string
}

func (this *Impl) SayHello(ctx context.Context, request *HelloRequest) (*HelloReply, error) {
	fmt.Printf("Server:%s", this.Addr)
	return &HelloReply{
		Message: "ddd",
	}, nil
}
