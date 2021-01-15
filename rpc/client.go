package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Allenxuxu/stark/rpc/client/selector"
	"google.golang.org/grpc"
)

var (
	DefaultTimeout = time.Second * 2
)

type Client struct {
	opts     *ClientOptions
	name     string
	selector selector.Selector

	conn *grpc.ClientConn
}

func NewClient(name string, s selector.Selector, opt ...ClientOption) (*Client, error) {
	opts := ClientOptions{
		Timeout: DefaultTimeout,
	}

	for _, o := range opt {
		o(&opts)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if len(s.Options().Balancer) != 0 {
		grpcOpts = append(grpcOpts, grpc.WithBalancerName(s.Options().Balancer))
	}

	opts.GrpcOpts = append(opts.GrpcOpts, grpcOpts...)
	client := &Client{
		opts:     &opts,
		name:     name,
		selector: s,
		conn:     nil,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *Client) connect() (err error) {
	timeCtx, cancel := context.WithTimeout(context.Background(), c.opts.Timeout)
	defer cancel()
	c.conn, err = grpc.DialContext(timeCtx, c.selector.Address(c.name), c.opts.GrpcOpts...)
	if err != nil {
		return fmt.Errorf("connect to %s error: %v", c.selector.Address(c.name), err)
	}

	return nil
}
