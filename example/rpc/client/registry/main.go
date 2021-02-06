package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Allenxuxu/stark/log"

	"github.com/Allenxuxu/stark"
	"github.com/Allenxuxu/stark/example/rpc/routeguide"
	"github.com/Allenxuxu/stark/registry/mdns"
	"github.com/Allenxuxu/stark/rpc"
	"github.com/Allenxuxu/stark/rpc/client/balancer"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"github.com/Allenxuxu/stark/rpc/client/selector/registry"
	"google.golang.org/grpc"
)

func main() {
	//rg, err := consul.NewRegistry()
	rg, err := mdns.NewRegistry()
	//rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}

	s, err := registry.NewSelector(rg,
		selector.BalancerName(balancer.RoundRobin),
	)
	if err != nil {
		panic(err)
	}

	client, err := stark.NewRPCClient("stark.rpc.test", s,
		rpc.GrpcDialOption(
			grpc.WithInsecure(),
		),
	)
	if err != nil {
		panic(err)
	}
	c := routeguide.NewRouteGuideClient(client.Conn())

	for i := 0; i < 20; i++ {
		resp, err := c.GetFeature(context.Background(), &routeguide.Point{
			Latitude:  11,
			Longitude: 0,
		})
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println(i, resp.Name, resp.Location.Latitude, resp.Location.Latitude)
	}

	stream, err := c.RouteChat(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		if err := stream.Send(&routeguide.RouteNote{
			Location: &routeguide.Point{
				Latitude:  1,
				Longitude: 1,
			},
			Message: "xx",
		}); err != nil {
			panic(err)
		}

		in, err := stream.Recv()
		if err == io.EOF {
			panic(err)

		}
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 5)
		log.Infof("[RouteChat] %v", in)
	}
}
