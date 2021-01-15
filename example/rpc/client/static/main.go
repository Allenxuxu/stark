package main

import (
	"context"
	"fmt"

	"github.com/Allenxuxu/stark/rpc"

	"github.com/Allenxuxu/stark/rpc/client/balancer"

	"github.com/Allenxuxu/stark/example/rpc/routeguide"
	"github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"github.com/Allenxuxu/stark/rpc/client/selector/static"
	"google.golang.org/grpc"
)

func main() {
	service := []*registry.Service{
		{
			Name:      "",
			Version:   "",
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
		selector.BalancerName(balancer.Random),
	)

	client, err := rpc.NewClient("stark.rpc.test", s,
		rpc.GrpcDialOption(
			grpc.WithInsecure(),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(1024*1024),
				grpc.MaxCallSendMsgSize(1024*1024)),
		),
	)

	if err != nil {
		panic(err)
	}

	c := routeguide.NewRouteGuideClient(client.Conn())

	resp, err := c.GetFeature(context.Background(), &routeguide.Point{
		Latitude:  0,
		Longitude: 0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Name, resp.Location.Latitude, resp.Location.Latitude)

}
