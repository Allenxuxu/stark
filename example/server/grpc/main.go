package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Allenxuxu/stark/example/server/grpc/unary"
	"github.com/Allenxuxu/stark/pkg/registry/etcd"
	"github.com/Allenxuxu/stark/server"
	"google.golang.org/grpc"
)

type GreetServer struct{}

func NewGreetServer() *GreetServer {
	return &GreetServer{}
}

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
		log.Println("steam send")
		time.Sleep(time.Second * 2)
		if err := steam.Send(&unary.Response{
			Greet: "hello from " + hostname,
		}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := time.Now()
		resp, err = handler(ctx, req)
		log.Printf("method: %s time: %v\n", info.FullMethod, time.Since(st))
		return resp, err
	}

	rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}
	s := server.NewServer(rg,
		server.Name("testserver"),
		server.Id("id123"),
		server.Address(":9091"),
		server.UnaryServerInterceptor(interceptor),
	)

	gs := NewGreetServer()
	unary.RegisterGreeterServer(s.GrpcServer(), gs)
	s.Register(gs)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-ch
		fmt.Println("stop")
		if err := s.Stop(); err != nil {
			panic(err)
		}
	}()

	if err := s.Start(); err != nil {
		panic(err)
	}
}
