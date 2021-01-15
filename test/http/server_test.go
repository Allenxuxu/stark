package http

import (
	"testing"
	"time"

	rg "github.com/Allenxuxu/stark/registry"
	"github.com/Allenxuxu/stark/registry/memory"
	"github.com/Allenxuxu/stark/rest"
	"github.com/Allenxuxu/stark/rest/client/selector/registry"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var name = "stark.http.test"

func TestServer(t *testing.T) {
	rg, err := memory.NewRegistry()
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	s := rest.NewSever(rg, r, rest.Name(name))

	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	client := newClient(rg)
	for i := 0; i < 10; i++ {
		request, err := client.Request()
		assert.Nil(t, err)
		resp, err := request.Get("/ping")
		assert.Nil(t, err)
		t.Log(resp)
	}

	assert.Nil(t, s.Stop())
}

func newClient(rg rg.Registry) *rest.Client {
	s, err := registry.NewSelector(rg)
	if err != nil {
		panic(err)
	}

	c, err := rest.NewClient(name, s)
	if err != nil {
		panic(err)
	}

	return c
}
