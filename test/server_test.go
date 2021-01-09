package test

import (
	"context"
	"testing"
	"time"

	"github.com/Allenxuxu/stark/registry/mdns"

	rg "github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/registry/memory"
	"github.com/Allenxuxu/stark/rpc/client"
	"github.com/Allenxuxu/stark/rpc/client/balancer"
	"github.com/Allenxuxu/stark/rpc/client/selector"
	"github.com/Allenxuxu/stark/rpc/client/selector/registry"
	"github.com/Allenxuxu/stark/rpc/server"
	pb "github.com/Allenxuxu/stark/test/routeguide"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var serverName = "stark.rpc.test"

func TestServer(t *testing.T) {
	memoryRegistry, err := memory.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}

	testServer(t, memoryRegistry)

	mdnsRegistry, err := mdns.NewRegistry()
	if err != nil {
		t.Fatal(err)
	}

	testServer(t, mdnsRegistry)
}

func testServer(t *testing.T, registry rg.Registry) {
	s := newServer(registry)

	go func() {
		time.Sleep(time.Second * 3)
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		c := pb.NewRouteGuideClient(newClient(registry).Conn())
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

func newServer(rg rg.Registry) *server.Server {
	s := server.NewServer(rg,
		server.Name(serverName),
	)

	rs := &routeGuideServer{}
	pb.RegisterRouteGuideServer(s.GrpcServer(), rs)

	return s
}

func newClient(rg rg.Registry) *client.Client {
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
