package stark

import (
	"net/http"

	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/rest"
	restSelector "github.com/Allenxuxu/stark/rest/client/selector"
	"github.com/Allenxuxu/stark/rpc"
	rpcSelector "github.com/Allenxuxu/stark/rpc/client/selector"
)

func NewRPCServer(rg registry.Registry, opt ...rpc.ServerOption) *rpc.Server {
	return rpc.NewServer(rg, opt...)
}

func NewRPCClient(name string, s rpcSelector.Selector, opt ...rpc.ClientOption) (*rpc.Client, error) {
	return rpc.NewClient(name, s, opt...)
}

func NewRestServer(rg registry.Registry, handler http.Handler, opts ...rest.ServerOption) *rest.Server {
	return rest.NewSever(rg, handler, opts...)
}

func NewRestClient(name string, s restSelector.Selector, opt ...rest.ClientOption) (*rest.Client, error) {
	return rest.NewClient(name, s, opt...)
}
