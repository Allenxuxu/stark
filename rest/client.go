package rest

import (
	"fmt"
	"sync"
	"time"

	"github.com/Allenxuxu/stark/log"

	"github.com/Allenxuxu/stark/rest/client/selector"
	"github.com/go-resty/resty/v2"
)

var (
	DefaultTimeout = time.Second * 2
)

type Client struct {
	opts     *ClientOptions
	name     string
	selector selector.Selector

	mu      sync.RWMutex
	clients map[string]*resty.Client
}

func NewClient(name string, s selector.Selector, opt ...ClientOption) (*Client, error) {
	opts := ClientOptions{
		Timeout: DefaultTimeout,
	}

	for _, o := range opt {
		o(&opts)
	}

	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}

	client := &Client{
		opts:     &opts,
		name:     name,
		selector: s,
		clients:  make(map[string]*resty.Client),
	}

	return client, nil
}

func (c *Client) Request() (*resty.Request, error) {
	node, err := c.selector.Next(c.name)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	client, ok := c.clients[node.Address]
	if !ok {
		client = resty.New()
		client.SetHostURL(fmt.Sprintf("http://%s", node.Address))
		client.SetTimeout(c.opts.Timeout)
		c.clients[node.Address] = client
	}
	c.mu.RUnlock()

	log.Info(node.Address)
	return client.R(), nil
}
