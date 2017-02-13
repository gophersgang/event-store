package api

import (
	"github.com/vendasta/gosdks/logging"
	grpc_context "golang.org/x/net/context"
	"github.com/vendasta/gosdks/pb/event-store/v1"
	"google.golang.org/grpc"
	"github.com/vendasta/event-store/pkg/event"
	"google.golang.org/grpc/codes"
	"github.com/vendasta/event-store/pkg/utils"
	"github.com/golang/protobuf/ptypes"
	"time"
)

type EventStoreServer struct{
	repository event.Repository
}

func (s *EventStoreServer) ListEvents(ctx grpc_context.Context, req *eventstore_v1.ListEventsRequest) (*eventstore_v1.ListEventsResponse, error) {
	logging.Debugf(ctx, "You should change this")
	res := &eventstore_v1.ListEventsResponse{
		Events: []*eventstore_v1.Event{&eventstore_v1.Event{EventId:"foo"}},
	}
	return res, nil
}

func (s *EventStoreServer) CreateEvent(ctx grpc_context.Context, req *eventstore_v1.CreateEventRequest) (*eventstore_v1.CreateEventResponse, error) {
	logging.Debugf(ctx, "You should change this")

	var err error
	var eventDate time.Time
	var aggregateType event.AggregateType
	var eventProto *eventstore_v1.Event

	aggregateType, err = event.AggregateTypeFromProto(req.AggregateType)
	if err != nil {
		logging.Errorf(ctx, "Failed to convert aggregateType from proto")
		return nil, toGrpcError(err)
	}

	if req.Timestamp != nil {
		eventDate, err = ptypes.Timestamp(req.Timestamp)
		if err != nil {
			logging.Errorf(ctx, "Failed to convert timestamp from proto")
			return nil, toGrpcError(utils.Error(utils.InvalidArgument, "Invalid timestamp"))
		}
	} else {
		eventDate = time.Time{}
	}

	e, err := s.repository.CreateEvent(ctx, aggregateType, req.AggregateId, req.Payload, eventDate)
	if err != nil {
		logging.Errorf(ctx, "Failed to create event")
		return nil, toGrpcError(err)
	}

	eventProto, err = e.Proto()
	if err != nil {
		logging.Errorf(ctx, "Failed to convert event to proto")
		return nil, toGrpcError(err)
	}
	res := &eventstore_v1.CreateEventResponse{
		Event: eventProto,
	}
	return res, nil
}

// Handle the error
func toGrpcError(err error) error {
	if e, ok := err.(utils.EventStoreError); ok {
		return e.GRPCError()
	}
	return grpc.Errorf(codes.Unknown, err.Error())
}

// NewEventServer returns a new event server.
func NewEventServer(s *grpc.Server, r event.Repository) {
	eventServer := EventStoreServer{r}
	eventstore_v1.RegisterEventStoreServer(s, &eventServer)
}
