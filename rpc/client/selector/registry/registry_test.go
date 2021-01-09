package registry

import (
	"testing"

	"github.com/Allenxuxu/stark/rpc/client/selector"

	"github.com/Allenxuxu/stark/pkg/registry"
	"github.com/Allenxuxu/stark/pkg/registry/memory"
	"github.com/stretchr/testify/assert"
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
	r, err := memory.NewRegistry(memory.Services(testData))
	assert.Nil(t, err)
	cache, err := NewSelector(r)
	assert.Nil(t, err)

	service, err := cache.GetService("foo")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(service))

	for _, s := range service {
		assert.Equal(t, s.Name, "foo")
		assert.Contains(t, []string{"1.0.3", "1.0.0", "1.0.1"}, s.Version)
	}
}

func TestRegistrySelectorFilter(t *testing.T) {

	version := "1.0.0"
	r, err := memory.NewRegistry(memory.Services(testData))
	assert.Nil(t, err)
	cache, err := NewSelector(r, selector.WithFilter(selector.FilterVersion(version)))
	assert.Nil(t, err)

	service, err := cache.GetService("foo")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(service))
	assert.Equal(t, service[0].Name, "foo")
	assert.Equal(t, service[0].Version, version)
}
