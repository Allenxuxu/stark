package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Allenxuxu/stark/example/server/grpc/unary"
	"google.golang.org/grpc/metadata"
)

type GreetServer struct{}

func (gs *GreetServer) Greet(ctx context.Context, req *unary.Request) (*unary.Response, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &unary.Response{
		Greet: "hello " + req.GetName() + " from " + hostname,
	}, nil
}
func (gs *GreetServer) Hello(in *unary.Request, steam unary.Greeter_HelloServer) error {
	name := in.GetName()
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	if err := steam.SendHeader(metadata.New(map[string]string{
		"name": name,
	})); err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 2)
		if err := steam.Send(&unary.Response{
			Greet: "hello from " + hostname,
		}); err != nil {
			return err
		}
	}
	return nil
}

func Test_extractEndpoints(t *testing.T) {
	service := &GreetServer{}
	tests := []string{"GreetServer.Greet", "GreetServer.Hello"}

	endpoints := extractEndpoints(service)
	assert.Equal(t, len(endpoints), 2)

	for _, e := range endpoints {
		assert.Contains(t, tests, e.Name)
	}

	endpoints = extractEndpoints(*service)
	assert.Equal(t, len(endpoints), 0)
}
