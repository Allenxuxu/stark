package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/Allenxuxu/ratelimit/tokenbucket"
	"github.com/Allenxuxu/stark"
	pb "github.com/Allenxuxu/stark/example/rpc/routeguide"
	"github.com/Allenxuxu/stark/registry/consul"
	"github.com/Allenxuxu/stark/rpc"
	"github.com/Allenxuxu/stark/rpc/server/middleware/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

type routeGuideServer struct{}

func NewServer() *routeGuideServer {
	return &routeGuideServer{}
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	log.Println("[GetFeature]", point.Latitude)
	return &pb.Feature{Location: point}, nil
}

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	log.Printf("[ListFeatures] %v", rect)

	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.Feature{
			Name: "feature",
			Location: &pb.Point{
				Latitude:  int32(i),
				Longitude: int32(i),
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}

		log.Printf("[RecordRoute] %v", point)
	}
}

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("[RouteChat] %v", in)

		if err := stream.Send(&pb.RouteNote{
			Location: in.Location,
			Message:  "reply " + in.Message,
		}); err != nil {
			return err
		}
	}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	st := time.Now()
	resp, err = handler(ctx, req)

	p, _ := peer.FromContext(ctx)
	log.Printf("method: %s time: %v, peer : %s\n", info.FullMethod, time.Since(st), p.Addr)
	return resp, err
}

func main() {
	rg, err := consul.NewRegistry()
	//rg, err := mdns.NewRegistry()
	//rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}

	s := stark.NewRPCServer(rg,
		rpc.Name("stark.rpc.test"),
		rpc.Version("v2.0.1"),
		rpc.Metadata(map[string]string{
			"server": "rpc",
			"test":   "1",
		}),
		rpc.MetricsAddress(":9091"),
		rpc.UnaryServerInterceptor(
			interceptor,
			ratelimit.UnaryServerInterceptor(tokenbucket.New(10, 5)),
		),
		rpc.StreamServerInterceptor(
			ratelimit.StreamServerInterceptor(tokenbucket.New(10, 5)),
		),
	)

	rs := NewServer()

	reflection.Register(s.GrpcServer())

	pb.RegisterRouteGuideServer(s.GrpcServer(), rs)

	if err := s.Start(); err != nil {
		panic(err)
	}
}
