package api

import (
	"github.com/vendasta/gosdks/logging"
	grpc_context "golang.org/x/net/context"
	"github.com/vendasta/event-store/pb"
)

type EventStoreServer struct{}

func (s *EventStoreServer) ListEvents(ctx grpc_context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	//TODO: CHANGE THIS TO DO STUFF
	logging.Debugf(ctx, "You should change this")
	res := &pb.ListEventsResponse{}
	return res, nil
}

