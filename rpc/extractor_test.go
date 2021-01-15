package rpc

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	pb "github.com/Allenxuxu/stark/example/rpc/routeguide"
	"github.com/stretchr/testify/assert"
)

type routeGuideServer struct{}

func newServer() *routeGuideServer {
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

func Test_extractEndpoints(t *testing.T) {
	service := newServer()
	tests := []string{
		"routeGuideServer.RouteChat",
		"routeGuideServer.RecordRoute",
		"routeGuideServer.GetFeature",
		"routeGuideServer.ListFeatures",
	}

	endpoints := extractEndpoints(service)
	assert.Equal(t, len(endpoints), 4)

	for _, e := range endpoints {
		assert.Contains(t, tests, e.Name)
	}

	endpoints = extractEndpoints(*service)
	assert.Equal(t, len(endpoints), 0)
}
