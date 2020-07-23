package registry

import (
	"google.golang.org/grpc/resolver"
)

type _Nacos struct {
}

func (this *_Nacos) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	panic("implement me")
}

func (this *_Nacos) Scheme() string {
	return string(RegistryTypeNacos)
}
