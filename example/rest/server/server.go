package main

import (
	"github.com/Allenxuxu/stark/registry/mdns"
	"github.com/Allenxuxu/stark/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	//rg, err := consul.NewRegistry()
	rg, err := mdns.NewRegistry()
	//rg, err := etcd.NewRegistry()
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

	s := rest.NewSever(rg, r,
		rest.Name("stark.http.test"),
	)

	if err := s.Start(); err != nil {
		panic(err)
	}
}
