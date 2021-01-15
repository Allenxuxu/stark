package selector

import (
	"testing"

	"github.com/Allenxuxu/stark/registry"
)

func TestStrategies(t *testing.T) {
	testData := []*registry.Service{
		{
			Name:    "test1",
			Version: "latest",
			Nodes: []*registry.Node{
				{
					Id:      "test1-1",
					Address: "10.0.0.1:1001",
				},
				{
					Id:      "test1-2",
					Address: "10.0.0.2:1002",
				},
			},
		},
		{
			Name:    "test1",
			Version: "default",
			Nodes: []*registry.Node{
				{
					Id:      "test1-3",
					Address: "10.0.0.3:1003",
				},
				{
					Id:      "test1-4",
					Address: "10.0.0.4:1004",
				},
			},
		},
	}

	for name, strategy := range map[string]Strategy{"random": Random(), "roundrobin": RoundRobin()} {

		counts := make(map[string]int)

		for i := 0; i < 100; i++ {
			node, err := strategy(testData)
			if err != nil {
				t.Fatal(err)
			}
			counts[node.Id]++
		}

		t.Logf("%s: %+v\n", name, counts)
	}
}
