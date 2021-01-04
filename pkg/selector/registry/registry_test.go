package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/Allenxuxu/stark/pkg/selector"

	"github.com/Allenxuxu/stark/pkg/registry/memory"
)

var (
	// mock data
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
	}
)

func TestRegistrySelector(t *testing.T) {
	counts := map[string]int{}

	r, err := memory.NewRegistry(memory.Services(testData))
	assert.Nil(t, err)
	cache, err := NewSelector(selector.Registry(r))
	assert.Nil(t, err)

	for i := 0; i < 100; i++ {
		node, err := cache.Next("foo")
		if err != nil {
			t.Errorf("Expected node err, got err: %v", err)
		}
		counts[node.Id]++
	}

	t.Logf("Selector Counts %v", counts)
}
