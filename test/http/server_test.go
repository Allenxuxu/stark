package http

import (
	"testing"
	"time"

	"github.com/Allenxuxu/stark/rest"

	"github.com/Allenxuxu/stark/registry/memory"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	rg, err := memory.NewRegistry()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	name := "stark.http.test"
	s := rest.NewSever(rg, r, rest.Name(name))

	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	client := resty.New()
	services, err := rg.GetService(name)
	assert.Nil(t, err)
	assert.Equal(t, len(services), 1)
	assert.Equal(t, services[0].Name, name)

	addr := services[0].Nodes[0].Address
	t.Log(addr)
	resp, err := client.R().Get("http://" + addr + "/ping")
	assert.Nil(t, err)
	t.Log(resp)

	assert.Nil(t, s.Stop())
}
