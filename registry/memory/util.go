package memory

import (
	"time"

	"github.com/Allenxuxu/stark/registry"
)

func serviceToRecord(s *registry.Service, ttl time.Duration) *record {
	nodes := make(map[string]*node)
	for _, n := range s.Nodes {
		nodes[n.Id] = &node{
			Node:     n,
			TTL:      ttl,
			LastSeen: time.Now(),
		}
	}

	endpoints := make([]*registry.Endpoint, len(s.Endpoints))
	for i, e := range s.Endpoints {
		endpoints[i] = e
	}

	return &record{
		Name:      s.Name,
		Version:   s.Version,
		Nodes:     nodes,
		Endpoints: endpoints,
	}
}

func recordToService(r *record) *registry.Service {
	metadata := make(map[string]string)
	for k, v := range r.Metadata {
		metadata[k] = v
	}

	endpoints := make([]*registry.Endpoint, len(r.Endpoints))
	for i, e := range r.Endpoints {
		request := new(registry.Value)
		if e.Request != nil {
			*request = *e.Request
		}
		response := new(registry.Value)
		if e.Response != nil {
			*response = *e.Response
		}

		metadata := make(map[string]string)
		for k, v := range e.Metadata {
			metadata[k] = v
		}

		endpoints[i] = &registry.Endpoint{
			Name:     e.Name,
			Request:  request,
			Response: response,
			Metadata: metadata,
		}
	}

	nodes := make([]*registry.Node, len(r.Nodes))
	i := 0
	for _, n := range r.Nodes {
		metadata := make(map[string]string)
		for k, v := range n.Metadata {
			metadata[k] = v
		}

		nodes[i] = &registry.Node{
			Id:       n.Id,
			Address:  n.Address,
			Metadata: metadata,
		}
		i++
	}

	return &registry.Service{
		Name:      r.Name,
		Version:   r.Version,
		Endpoints: endpoints,
		Nodes:     nodes,
	}
}
