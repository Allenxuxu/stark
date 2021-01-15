package main

import (
	"github.com/Allenxuxu/stark/log"
	"github.com/Allenxuxu/stark/registry/consul"
	"github.com/Allenxuxu/stark/rest"
	"github.com/Allenxuxu/stark/rest/client/selector/registry"
)

func main() {
	rg, err := consul.NewRegistry()
	//rg, err := mdns.NewRegistry()
	//rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}

	s, err := registry.NewSelector(rg)
	if err != nil {
		panic(err)
	}

	c, err := rest.NewClient("stark.http.test", s)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		r, err := c.Request()
		if err != nil {
			panic(err)
		}

		resp, err := r.Get("/ping")
		if err != nil {
			panic(err)
		}

		log.Info(resp)
	}

}
