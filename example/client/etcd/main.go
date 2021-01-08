package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Allenxuxu/stark/client/balancer"

	"google.golang.org/grpc"

	"github.com/Allenxuxu/stark/client"
	"github.com/Allenxuxu/stark/client/selector"
	"github.com/Allenxuxu/stark/client/selector/registry"
	"github.com/Allenxuxu/stark/example/server/grpc/unary"
	"github.com/Allenxuxu/stark/pkg/registry/etcd"
)

func main() {
	rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}

	s, err := registry.NewSelector(rg,
		selector.WithBalancerName(balancer.RoundRobin),
	)
	if err != nil {
		panic(err)
	}

	client, err := client.NewClient("stark.rpc.test", s,
		client.DialOption(
			grpc.WithInsecure(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := unary.NewGreeterClient(client.Conn())

	for i := 0; i < 10; i++ {

		resp, err := c.Greet(context.Background(), &unary.Request{Name: "xuxu"})
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.Greet)
		time.Sleep(time.Second)
	}
}
