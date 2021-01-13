package main

import (
	"github.com/Allenxuxu/stark/registry/consul"
	"github.com/Allenxuxu/stark/rest/server"
	"github.com/gin-gonic/gin"
)

func main() {
	rg, err := consul.NewRegistry()
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

	s := server.NewSever(rg, r)

	if err := s.Run(); err != nil {
		panic(err)
	}
}
