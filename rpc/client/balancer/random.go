package balancer

import (
	"context"
	"math/rand"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const Random = "random"

func init() {
	balancer.Register(newRandom())
	rand.Seed(time.Now().UnixNano())
}

func newRandom() balancer.Builder {
	return base.NewBalancerBuilder(Random, &randomPickerBuilder{})
}

type randomPickerBuilder struct{}

func (*randomPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var conns []balancer.SubConn
	for _, conn := range readySCs {
		conns = append(conns, conn)
	}

	return &randomPicker{
		subConns: conns,
	}
}

type randomPicker struct {
	subConns []balancer.SubConn
}

func (p *randomPicker) Pick(ctx context.Context, info balancer.PickInfo) (conn balancer.SubConn, done func(balancer.DoneInfo), err error) {
	conn = p.subConns[rand.Int()%len(p.subConns)]
	return
}
