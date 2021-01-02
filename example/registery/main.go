package main

import (
	"github.com/Allenxuxu/stark/log"
	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/Allenxuxu/stark/pkg/registry/etcd"
)

var (
	testData = map[string][]*registry.Service{
		"foo": {
			{
				Name:    "foo",
				Version: "1.0.0",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.0-123",
						Address: "localhost:9999",
					},
					{
						Id:      "foo-1.0.0-321",
						Address: "localhost:9999",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.1",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.1-321",
						Address: "localhost:6666",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.3",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.3-345",
						Address: "localhost:8888",
					},
				},
			},
		},
		"bar": {
			{
				Name:    "bar",
				Version: "default",
				Nodes: []*registry.Node{
					{
						Id:      "bar-1.0.0-123",
						Address: "localhost:9999",
					},
					{
						Id:      "bar-1.0.0-321",
						Address: "localhost:9999",
					},
				},
			},
			{
				Name:    "bar",
				Version: "latest",
				Nodes: []*registry.Node{
					{
						Id:      "bar-1.0.1-321",
						Address: "localhost:6666",
					},
				},
			},
		},
	}
)

func main() {
	rg, err := etcd.NewRegistry(
		registry.Addrs("localhost:2379"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info(rg.String())

	for k, v := range testData {
		log.Info("service name:", k)

		for _, s := range v {
			if err := rg.Register(s); err != nil {
				log.Fatalf("Register service err : %v", err)
			}
			log.Infof("Register service  %s , version %s", k, s.Version)
		}
	}

	for k, _ := range testData {
		services, err := rg.GetService(k)
		if err != nil {
			log.Fatalf("GetService err :%v", err)
		}

		for _, s := range services {
			log.Infof("GetService service : %s %s", s.Name, s.Version)
		}
	}

	for _, s := range testData["foo"] {
		if err := rg.Deregister(s); err != nil {
			log.Fatal(err)
		}
	}

	// just use name and node.id
	if err := rg.Deregister(&registry.Service{
		Name:      "bar",
		Version:   "",
		Metadata:  nil,
		Endpoints: nil,
		Nodes: []*registry.Node{
			{
				Id:      "bar-1.0.1-321",
				Address: "localhost:6666",
			},
		},
	}); err != nil {
		log.Fatal(err)
	}

	services, err := rg.ListServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range services {
		log.Infof("ListServices services %s %s", v.Name, v.Version)
		for _, node := range v.Nodes {
			log.Infof("ListServices node %s %s", node.Id, node.Address)
		}
	}

	w, err := rg.Watch()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		r, err := w.Next()
		if err != nil {
			log.Fatal(err)
		}

		log.Infof("Watch service %s, action %s", r.Service.Name, r.Action)

		if i == 3 {
			w.Stop()
		}
	}

}
