package api

import (
	"github.com/vendasta/gosdks/logging"
	grpc_context "golang.org/x/net/context"
	"github.com/vendasta/gosdks/pb/event-store/v1"
)

type EventStoreServer struct{}

func (s *EventStoreServer) ListEvents(ctx grpc_context.Context, req *eventstore_v1.ListEventsRequest) (*eventstore_v1.ListEventsResponse, error) {
	logging.Debugf(ctx, "You should change this")
	res := &eventstore_v1.ListEventsResponse{
		Events: []*eventstore_v1.Event{&eventstore_v1.Event{EventId:"foo"}},
	}
	return res, nil
}

func (s *EventStoreServer) CreateEvent(ctx grpc_context.Context, req *eventstore_v1.CreateEventRequest) (*eventstore_v1.CreateEventResponse, error) {
	logging.Debugf(ctx, "You should change this")
	res := &eventstore_v1.CreateEventResponse{
		Event: &eventstore_v1.Event{EventId:"foo"},
	}
	return res, nil
}

