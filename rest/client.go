package rest

import (
	"fmt"
	"sync"
	"time"

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

	return c.nextClient(node.Address).R(), nil
}

func (c *Client) nextClient(address string) *resty.Client {
	c.mu.RLock()
	client, ok := c.clients[address]
	c.mu.RUnlock()

	if !ok {
		client = c.newRestyClient(address)

		c.mu.Lock()
		c.clients[address] = client
		c.mu.Unlock()
	}

	return client
}

func (c *Client) newRestyClient(address string) *resty.Client {
	return resty.New().
		SetHostURL(fmt.Sprintf("http://%s", address)).
		SetTimeout(c.opts.Timeout)
}
