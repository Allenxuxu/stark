package main

import (
	"context"
	"fmt"

	"github.com/Allenxuxu/stark/client/balancer"

	"github.com/Allenxuxu/stark/client"
	"github.com/Allenxuxu/stark/client/selector"
	"github.com/Allenxuxu/stark/client/selector/static"
	"github.com/Allenxuxu/stark/example/server/grpc/unary"
	"github.com/Allenxuxu/stark/pkg/registry"
	"google.golang.org/grpc"
)

func main() {
	service := []*registry.Service{
		{
			Name:      "",
			Version:   "",
			Metadata:  nil,
			Endpoints: nil,
			Nodes: []*registry.Node{
				{Address: "127.0.0.1:9092"},

				{Address: "127.0.0.1:9091"},
				{Address: "127.0.0.1:9092"},
			},
		},
	}

	s := static.NewSelector(
		service,
		selector.WithBalancerName(balancer.Random),
	)

	client, err := client.NewClient("stark.rpc.test", s,
		client.DialOption(
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024)),
		),
	)

	if err != nil {
		panic(err)
	}

	c := unary.NewGreeterClient(client.Conn())

	resp, err := c.Greet(context.Background(), &unary.Request{Name: "xuxu"})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Greet)
}
