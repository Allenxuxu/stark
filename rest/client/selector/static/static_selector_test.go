package static

import (
	"github.com/Allenxuxu/stark/registry"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestStaticSelector_Next(t *testing.T) {
	service := []*registry.Service{
		{
			Name:      "",
			Version:   "",
			Endpoints: nil,
			Nodes: []*registry.Node{
				{Address: "127.0.0.1:9092"},
				{Address: "127.0.0.1:9091"},
				{Address: "127.0.0.1:9092"},
			},
		},
	}

	selector := New(service)
	for _, node := range service[0].Nodes {
		n, err := selector.Next("")
		assert.Nil(t, err)
		assert.Equal(t, n.Address, node.Address)
	}
}
