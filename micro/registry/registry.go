package registry

import "google.golang.org/grpc/resolver"

type RegistryType string

const (
	RegistryTypeNacos = RegistryType("nacos")
	RegistryTypeEtcd  = RegistryType("etcd")
)

type Registry interface {
	resolver.Builder
}

func New(t RegistryType) Registry {
	switch t {
	case RegistryTypeNacos:
		return nil
	case RegistryTypeEtcd:
		return nil
	default:
		return nil
	}
}
