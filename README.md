# stark

[![Github Actions](https://github.com/Allenxuxu/gev/workflows/CI/badge.svg)](https://github.com/Allenxuxu/stark/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Allenxuxu/stark)](https://goreportcard.com/report/github.com/Allenxuxu/stark)
[![LICENSE](https://img.shields.io/badge/LICENSE-MIT-blue)](https://github.com/Allenxuxu/stark/blob/master/LICENSE)

> stark 命名源自漫威 Tony Stark

一个用于构建分布式系统的工具集或者轻框架，支持 grpc 和 http ，支持多种注册中心 consul ，etcd , mdns 等。

## 功能

- 服务发现
- 负载均衡
- grpc/http Server 和 Client
- ...

## stark ctl

```shell
go get  github.com/Allenxuxu/stark/cmd/stark
```

查询注册服务

```shell
stark service -r consul -rg 127.0.0.1:8500 {service name}
```

## Example

[grpc server](example/rpc/server/main.go)

```go
rg, err := mdns.NewRegistry()
if err != nil {
    panic(err)
}

s := stark.NewRPCServer(rg,
    rpc.Name("stark.rpc.test"),
    rpc.Version("v0.0.1"),
)

rs := &routeGuideServer{}
pb.RegisterRouteGuideServer(s.GrpcServer(), rs)

if err := s.Start(); err != nil {
    panic(err)
}
```

```shell
cd example/rpc/server
go run main.go
```

使用 stark 工具 查看已经注册的服务

```shell
stark service stark.rpc.test
```

[grpc client](example/rpc/client/registry)

```go
rg, err := mdns.NewRegistry()
if err != nil {
    panic(err)
}

s, err := registry.NewSelector(rg,
    selector.BalancerName(balancer.RoundRobin),
)
if err != nil {
    panic(err)
}

client, err := stark.NewRPCClient("stark.rpc.test", s,
    rpc.GrpcDialOption(
        grpc.WithInsecure(),
    ),
)
if err != nil {
    panic(err)
}

c := routeguide.NewRouteGuideClient(client.Conn())

for i := 0; i < 10; i++ {
    resp, err := c.GetFeature(context.Background(), &routeguide.Point{
        Latitude:  0,
        Longitude: 0,
    })
    if err != nil {
        panic(err)
    }
}
```

[http server](example/rest/server/main.go)

```go
rg, err := mdns.NewRegistry()
if err != nil {
    panic(err)
}

r := gin.Default()
r.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
})

s := stark.NewRestServer(rg, r,
    rest.Name("stark.http.test"),
)

if err := s.Start(); err != nil {
    panic(err)
}
```

[http client](example/rest/client/main.go)

```go
rg, err := mdns.NewRegistry()
if err != nil {
    panic(err)
}

s, err := registry.NewSelector(rg)
if err != nil {
    panic(err)
}

c, err := stark.NewRestClient("stark.http.test", s)
if err != nil {
    panic(err)
}

for i := 0; i < 5; i++ {
    r, err := c.Request()
    if err != nil {
        panic(err)
    }

    resp, err := r.Get("/ping")
    if err != nil {
        panic(err)
    }

    log.Info(resp)
}
```

## 感谢

[go-micro](https://github.com/asim/go-micro) 
本项目就是受到 go-micro 的启发而来，本意就是构建一个比 go-micro 更加轻量，够用就好的微服务框架。
项目中 config，registry 都是从 go-micro v1.18 改造而来。
