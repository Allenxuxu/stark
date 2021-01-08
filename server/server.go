package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Allenxuxu/stark/log"
	"github.com/Allenxuxu/toolkit/sync"
	"google.golang.org/grpc/reflection"

	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	DefaultAddress          = ":0"
	DefaultName             = "stark.rpc.server"
	DefaultVersion          = time.Now().Format("2006.01.02.15.04")
	DefaultId               = uuid.New().String()
	DefaultRegisterCheck    = func(context.Context) error { return nil }
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)

type Server struct {
	opts      *Options
	registry  registry.Registry
	grpcSever *grpc.Server
	service   *registry.Service
	exit      chan struct{}
	sw        sync.WaitGroupWrapper

	options            []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

func NewServer(rg registry.Registry, opt ...Option) *Server {
	opts := Options{
		Metadata:         nil,
		Name:             DefaultName,
		Address:          DefaultAddress,
		Id:               DefaultId,
		Version:          DefaultVersion,
		RegisterCheck:    DefaultRegisterCheck,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterInterval: DefaultRegisterInterval,
	}

	for _, o := range opt {
		o(&opts)
	}

	g := &Server{
		opts:     &opts,
		registry: rg,
		exit:     make(chan struct{}),
	}
	g.grpcSever = grpc.NewServer(opts.GrpcOpts...)
	g.service = &registry.Service{
		Name:      g.opts.Name,
		Version:   g.opts.Version,
		Metadata:  g.opts.Metadata,
		Endpoints: nil,
		Nodes: []*registry.Node{{
			Id:       g.opts.Id,
			Address:  g.opts.Address,
			Metadata: g.opts.Metadata},
		},
	}
	return g
}

func (g *Server) Register(service ...interface{}) {
	var endpoints []*registry.Endpoint
	for _, s := range service {
		endpoints = append(endpoints, extractEndpoints(s)...)
	}

	g.service.Endpoints = endpoints
}

func (g *Server) GrpcServer() *grpc.Server {
	return g.grpcSever
}

func (g *Server) Start() error {
	listener, err := net.Listen("tcp", g.opts.Address)
	if err != nil {
		return err
	}

	g.opts.Address = listener.Addr().String()
	g.service.Nodes[0].Address = listener.Addr().String()

	if err := g.register(); err != nil {
		return err
	}

	reflection.Register(g.grpcSever)
	log.Infof("%s server listen on %s", g.opts.Name, g.opts.Address)
	return g.grpcSever.Serve(listener)
}

func (g *Server) Stop() error {
	select {
	case <-g.exit:
		return nil
	default:
		close(g.exit)
		g.sw.Wait()
	}

	g.grpcSever.GracefulStop()
	return nil
}

func (g *Server) Options() Options {
	return *g.opts
}

func (g *Server) String() string {
	return "grpc"
}

func (g *Server) nodeId() string {
	return fmt.Sprintf("%s-%s", g.opts.Name, g.opts.Id)
}

func (g *Server) register() error {
	ttlOpt := registry.RegisterTTL(g.opts.RegisterTTL)
	if err := g.registry.Register(g.service, ttlOpt); err != nil {
		return err
	}

	g.sw.AddAndRun(func() {
		t := new(time.Ticker)
		if g.opts.RegisterInterval > time.Duration(0) {
			t = time.NewTicker(g.opts.RegisterInterval)
		}

	Loop:
		for {
			select {
			case <-t.C:
				if err := g.registry.Register(g.service, ttlOpt); err != nil {
					log.Log("Server register error: ", err)
				}
			case <-g.exit:
				break Loop
			}
		}

		if err := g.registry.Deregister(g.service); err != nil {
			log.Log("Server deregister error: ", err)
		}
	})

	return nil
}
