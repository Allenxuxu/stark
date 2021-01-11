package service

import (
	"fmt"

	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/registry/consul"
	"github.com/Allenxuxu/stark/registry/etcd"
	"github.com/Allenxuxu/stark/registry/mdns"
)

func newRegistry(name, addr string) (registry.Registry, error) {
	switch name {
	case "mdns":
		return mdns.NewRegistry()
	case "etcd":
		if len(addr) > 0 {
			return etcd.NewRegistry(registry.Addrs(addr))
		}
		return etcd.NewRegistry()
	case "consul":
		if len(addr) > 0 {
			return consul.NewRegistry(registry.Addrs(addr))
		}
		return consul.NewRegistry()
	default:
		return nil, fmt.Errorf("%s not supported", name)
	}
}
