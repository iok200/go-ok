package example

import context "context"

type Impl struct {
}

func (this *Impl) SayHello(ctx context.Context, request *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "ddd",
	}, nil
}
