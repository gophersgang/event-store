package api

import (
	"github.com/vendasta/gosdks/logging"
	grpc_context "golang.org/x/net/context"
	"github.com/vendasta/event-store/pb"
)

type EventStoreServer struct{}

func (s *EventStoreServer) ListEventsForAggregateType(ctx grpc_context.Context, req *pb.ListEventsForAggregateTypeRequest) (*pb.ListEventsForAggregateTypeResponse, error) {
	//TODO: CHANGE THIS TO DO STUFF
	logging.Debugf(ctx, "You should change this")
	res := &pb.ListEventsForAggregateTypeResponse{
		Events: []pb.Event{pb.AggregateType(pb.AggregateType_CAMPAIGN)},
	}
	return res, nil
}

