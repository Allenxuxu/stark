package rpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Allenxuxu/stark/rpc/metrics"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/Allenxuxu/stark/log"
	"github.com/Allenxuxu/stark/registry"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var (
	DefaultAddress          = ":0"
	DefaultName             = "stark.rpc.server"
	DefaultVersion          = time.Now().Format("2006.01.02.15.04")
	DefaultId               = uuid.New().String()
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)

type Server struct {
	opts      *ServerOptions
	registry  registry.Registry
	grpcSever *grpc.Server
	service   *registry.Service
	exit      chan struct{}
}

func NewServer(rg registry.Registry, opt ...ServerOption) *Server {
	opts := ServerOptions{
		Metadata:         nil,
		Name:             DefaultName,
		Address:          DefaultAddress,
		Id:               DefaultId,
		Version:          DefaultVersion,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterInterval: DefaultRegisterInterval,
	}

	for _, o := range opt {
		o(&opts)
	}

	opts.GrpcOpts = append(opts.GrpcOpts, grpc.ChainStreamInterceptor(
		grpc_recovery.StreamServerInterceptor(),
	))
	opts.GrpcOpts = append(opts.GrpcOpts, grpc.ChainUnaryInterceptor(
		grpc_recovery.UnaryServerInterceptor(),
	))

	g := &Server{
		opts:     &opts,
		registry: rg,
		exit:     make(chan struct{}),
	}

	g.grpcSever = grpc.NewServer(opts.GrpcOpts...)
	g.service = &registry.Service{
		Name:    g.opts.Name,
		Version: g.opts.Version,
		Nodes: []*registry.Node{{
			Id:       g.opts.Id,
			Address:  g.opts.Address,
			Metadata: g.opts.Metadata},
		},
	}

	// metrics
	if len(opts.Metadata[metaDataMetricsAddressKey]) > 0 {
		grpcMetrics := grpc_prometheus.NewServerMetrics()
		metrics.PrometheusMustRegister(grpcMetrics)

		opts.GrpcOpts = append(opts.GrpcOpts, grpc.ChainStreamInterceptor(
			grpcMetrics.StreamServerInterceptor(),
		))
		opts.GrpcOpts = append(opts.GrpcOpts, grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
		))

		grpcMetrics.InitializeMetrics(g.grpcSever)
	}

	return g
}

func (g *Server) RegisterEndpoints(service ...interface{}) {
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
	if err = g.register(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case sig := <-ch:
			log.Infof("Received signal %s", sig)
			if err = g.Stop(); err != nil {
				log.Error("Server stop error :%v", err)
			}
		case <-g.exit:
		}

		if err := g.deregister(); err != nil {
			log.Errorf("deregister error %v", err)
		}
	}()

	// metrics
	if len(g.opts.Metadata[metaDataMetricsAddressKey]) > 0 {
		go func() {
			if err := metrics.Run(g.opts.MetricsPath, g.opts.Metadata[metaDataMetricsAddressKey]); err != nil {
				log.Fatalf("Run metrics server error: %v", err)
			}
		}()
	}

	log.Infof("RPC server listen on %s", g.opts.Address)
	return g.grpcSever.Serve(listener)
}

func (g *Server) Stop() error {
	select {
	case <-g.exit:
		return nil
	default:
		close(g.exit)
		g.grpcSever.GracefulStop()
		return nil
	}
}

func (g *Server) Options() ServerOptions {
	return *g.opts
}

func (g *Server) String() string {
	return "grpc"
}

func (g *Server) register() error {
	ttlOpt := registry.RegisterTTL(g.opts.RegisterTTL)
	if err := g.registry.Register(g.service, ttlOpt); err != nil {
		return err
	}
	log.Infof("Registry [%s] register node: %s", g.registry.String(), g.service.Nodes[0].Id)

	if g.opts.RegisterInterval <= time.Duration(0) {
		return nil
	}

	go func() {
		t := time.NewTicker(g.opts.RegisterInterval)

		for {
			select {
			case <-t.C:
				if err := g.registry.Register(g.service, ttlOpt); err != nil {
					log.Errorf("Server register error: %v", err)
				}
			case <-g.exit:
				t.Stop()
				return
			}
		}
	}()

	return nil
}

func (g *Server) deregister() error {
	log.Infof("Registry [%s] deregister node: %s", g.registry.String(), g.service.Nodes[0].Id)
	return g.registry.Deregister(g.service)
}
