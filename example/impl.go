package example

import (
	context "context"
)

type Impl struct {
	Addr string
}

func (this *Impl) SayHello(ctx context.Context, request *HelloRequest) (*HelloReply, error) {
	//log.Infof("Server:%s", this.Addr)
	return &HelloReply{
		Message: "ddd",
	}, nil
}
