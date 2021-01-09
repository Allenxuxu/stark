package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Allenxuxu/stark/registry/mdns"
	"github.com/Allenxuxu/stark/rpc/client"
	"github.com/Allenxuxu/stark/rpc/client/balancer"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"github.com/Allenxuxu/stark/rpc/client/selector/registry"
	"github.com/Allenxuxu/stark/rpc/server"
	pb "github.com/Allenxuxu/stark/test/routeguide"
	"google.golang.org/grpc"
)

var serverName = "stark.rpc.test"

func TestServer(t *testing.T) {
	s := newServer()

	go func() {
		time.Sleep(time.Second * 3)
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		c := pb.NewRouteGuideClient(newClient().Conn())
		for i := 0; i < 10; i++ {
			p := &pb.Point{
				Latitude:  int32(i),
				Longitude: int32(i),
			}
			resp, err := c.GetFeature(context.Background(), p)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, p.Latitude, resp.Location.Latitude)
			assert.Equal(t, p.Longitude, resp.Location.Longitude)
		}
	}()

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
}

func newServer() *server.Server {
	rg, err := mdns.NewRegistry()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(rg,
		server.Name(serverName),
	)

	rs := &routeGuideServer{}
	pb.RegisterRouteGuideServer(s.GrpcServer(), rs)

	return s
}

func newClient() *client.Client {
	rg, err := mdns.NewRegistry()
	if err != nil {
		panic(err)
	}

	s, err := registry.NewSelector(rg,
		selector.WithBalancerName(balancer.RoundRobin),
	)
	if err != nil {
		panic(err)
	}

	c, err := client.NewClient(serverName, s,
		client.DialOption(
			grpc.WithInsecure(),
		),
	)
	if err != nil {
		panic(err)
	}

	return c
}
