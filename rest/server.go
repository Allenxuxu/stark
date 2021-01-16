package rest

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/Allenxuxu/stark/log"
	"github.com/Allenxuxu/stark/registry"
)

var (
	DefaultName    = "stark.http.server"
	DefaultVersion = time.Now().Format("2006.01.02.15.04")
	DefaultId      = uuid.New().String()
	DefaultAddress = ":0"

	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Minute
)

type Server struct {
	opts     ServerOptions
	registry registry.Registry
	handler  http.Handler
	server   *http.Server
	service  *registry.Service
	exit     chan struct{}
}

func NewSever(rg registry.Registry, handler http.Handler, opts ...ServerOption) *Server {
	options := ServerOptions{
		Name:             DefaultName,
		Version:          DefaultVersion,
		Id:               DefaultId,
		Address:          DefaultAddress,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterInterval: DefaultRegisterInterval,
	}

	for _, o := range opts {
		o(&options)
	}

	s := &Server{
		opts:     options,
		registry: rg,
		handler:  handler,
		exit:     make(chan struct{}),
	}

	s.service = &registry.Service{
		Name:    s.opts.Name,
		Version: s.opts.Version,
		Nodes: []*registry.Node{{
			Id:       s.opts.Id,
			Address:  s.opts.Address,
			Metadata: s.opts.Metadata,
		}},
	}
	return s
}

func (s *Server) Start() error {
	s.server = &http.Server{Addr: s.opts.Address, Handler: s.handler}
	ln, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	s.service.Nodes[0].Address = ln.Addr().String()
	s.opts.Address = ln.Addr().String()
	if err := s.register(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case sig := <-ch:
			log.Infof("Received signal %s", sig)
			if err = s.Stop(); err != nil {
				log.Errorf("Server stop error :%v", err)
			}
		case <-s.exit:
		}

		if err := s.deregister(); err != nil {
			log.Errorf("deregister error %v", err)
		}
	}()

	log.Infof("Http server listen on %s", s.opts.Address)
	if len(s.opts.CertFile) > 0 && len(s.opts.KeyFile) > 0 {
		err = s.server.ServeTLS(ln, s.opts.CertFile, s.opts.KeyFile)
	} else {
		err = s.server.Serve(ln)
	}
	if err != nil && err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (s *Server) Options() ServerOptions {
	return s.opts
}

func (s *Server) Stop() error {
	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
		return s.server.Close()
	}
}

func (s *Server) String() string {
	return "http"
}

func (s *Server) register() error {
	ttlOpt := registry.RegisterTTL(s.opts.RegisterTTL)
	if err := s.registry.Register(s.service, ttlOpt); err != nil {
		return err
	}
	log.Infof("Registry [%s] register node: %s", s.registry.String(), s.service.Nodes[0].Id)

	if s.opts.RegisterInterval <= time.Duration(0) {
		return nil
	}

	go func() {
		t := time.NewTicker(s.opts.RegisterInterval)

		for {
			select {
			case <-t.C:
				if err := s.registry.Register(s.service, ttlOpt); err != nil {
					log.Errorf("Server register error: %v", err)
				}
			case <-s.exit:
				t.Stop()
				return
			}
		}
	}()

	return nil
}

func (s *Server) deregister() error {
	log.Infof("Registry [%s] deregister node: %s", s.registry.String(), s.service.Nodes[0].Id)
	return s.registry.Deregister(s.service)
}
