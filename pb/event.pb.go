// Code generated by protoc-gen-go.
// source: event.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	event.proto

It has these top-level messages:
	Event
	ListEventsRequest
	ListEventsResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type AggregateType int32

const (
	AggregateType_CAMPAIGN AggregateType = 0
)

var AggregateType_name = map[int32]string{
	0: "CAMPAIGN",
}
var AggregateType_value = map[string]int32{
	"CAMPAIGN": 0,
}

func (x AggregateType) String() string {
	return proto.EnumName(AggregateType_name, int32(x))
}
func (AggregateType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Event struct {
	EventId       string                     `protobuf:"bytes,1,opt,name=event_id,json=eventId" json:"event_id,omitempty"`
	AggregateType AggregateType              `protobuf:"varint,2,opt,name=aggregate_type,json=aggregateType,enum=eventstore.v1.AggregateType" json:"aggregate_type,omitempty"`
	AggregateId   string                     `protobuf:"bytes,3,opt,name=aggregate_id,json=aggregateId" json:"aggregate_id,omitempty"`
	Payload       string                     `protobuf:"bytes,4,opt,name=payload" json:"payload,omitempty"`
	Timestamp     *google_protobuf.Timestamp `protobuf:"bytes,5,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Event) GetEventId() string {
	if m != nil {
		return m.EventId
	}
	return ""
}

func (m *Event) GetAggregateType() AggregateType {
	if m != nil {
		return m.AggregateType
	}
	return AggregateType_CAMPAIGN
}

func (m *Event) GetAggregateId() string {
	if m != nil {
		return m.AggregateId
	}
	return ""
}

func (m *Event) GetPayload() string {
	if m != nil {
		return m.Payload
	}
	return ""
}

func (m *Event) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

// *******************
// Requests
// *******************
type ListEventsRequest struct {
	AggregateType AggregateType `protobuf:"varint,1,opt,name=aggregate_type,json=aggregateType,enum=eventstore.v1.AggregateType" json:"aggregate_type,omitempty"`
}

func (m *ListEventsRequest) Reset()         { *m = ListEventsRequest{} }
func (m *ListEventsRequest) String() string { return proto.CompactTextString(m) }
func (*ListEventsRequest) ProtoMessage()    {}
func (*ListEventsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{1}
}

func (m *ListEventsRequest) GetAggregateType() AggregateType {
	if m != nil {
		return m.AggregateType
	}
	return AggregateType_CAMPAIGN
}

// *******************
// Responses
// *******************
type ListEventsResponse struct {
	Events []*Event `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
}

func (m *ListEventsResponse) Reset()         { *m = ListEventsResponse{} }
func (m *ListEventsResponse) String() string { return proto.CompactTextString(m) }
func (*ListEventsResponse) ProtoMessage()    {}
func (*ListEventsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2}
}

func (m *ListEventsResponse) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

func init() {
	proto.RegisterType((*Event)(nil), "eventstore.v1.Event")
	proto.RegisterType((*ListEventsRequest)(nil), "eventstore.v1.ListEventsRequest")
	proto.RegisterType((*ListEventsResponse)(nil), "eventstore.v1.ListEventsResponse")
	proto.RegisterEnum("eventstore.v1.AggregateType", AggregateType_name, AggregateType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EventStore service

type EventStoreClient interface {
	ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error)
}

type eventStoreClient struct {
	cc *grpc.ClientConn
}

func NewEventStoreClient(cc *grpc.ClientConn) EventStoreClient {
	return &eventStoreClient{cc}
}

func (c *eventStoreClient) ListEvents(ctx context.Context, in *ListEventsRequest, opts ...grpc.CallOption) (*ListEventsResponse, error) {
	out := new(ListEventsResponse)
	err := grpc.Invoke(ctx, "/eventstore.v1.EventStore/ListEvents", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for EventStore service

type EventStoreServer interface {
	ListEvents(context.Context, *ListEventsRequest) (*ListEventsResponse, error)
}

func RegisterEventStoreServer(s *grpc.Server, srv EventStoreServer) {
	s.RegisterService(&_EventStore_serviceDesc, srv)
}

func _EventStore_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventStoreServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eventstore.v1.EventStore/ListEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventStoreServer).ListEvents(ctx, req.(*ListEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EventStore_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eventstore.v1.EventStore",
	HandlerType: (*EventStoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListEvents",
			Handler:    _EventStore_ListEvents_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event.proto",
}

func init() { proto.RegisterFile("event.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 320 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x51, 0x5b, 0x4f, 0xf2, 0x30,
	0x18, 0xa6, 0x1f, 0x1f, 0xa7, 0x77, 0x40, 0x4c, 0xe3, 0x45, 0x5d, 0x34, 0x8e, 0x5d, 0x2d, 0xc6,
	0x14, 0x99, 0x37, 0xde, 0x12, 0xa2, 0x86, 0x44, 0x8d, 0x99, 0xdc, 0x93, 0x91, 0xbd, 0xce, 0x25,
	0x40, 0xeb, 0x5a, 0x48, 0xb8, 0xf6, 0x1f, 0xf8, 0x03, 0xfd, 0x2d, 0x86, 0x8e, 0x83, 0x23, 0x1e,
	0xc2, 0xe5, 0xd3, 0x3e, 0xc7, 0x16, 0x2c, 0x9c, 0xe3, 0x54, 0x73, 0x99, 0x0a, 0x2d, 0x68, 0xc3,
	0x00, 0xa5, 0x45, 0x8a, 0x7c, 0xde, 0xb1, 0x4f, 0x63, 0x21, 0xe2, 0x31, 0xb6, 0xcd, 0xe5, 0x68,
	0xf6, 0xdc, 0xd6, 0xc9, 0x04, 0x95, 0x0e, 0x27, 0x32, 0xe3, 0xbb, 0x1f, 0x04, 0x4a, 0xd7, 0x4b,
	0x09, 0x3d, 0x82, 0xaa, 0xd1, 0x0e, 0x93, 0x88, 0x11, 0x87, 0x78, 0xb5, 0xa0, 0x62, 0x70, 0x3f,
	0xa2, 0x3d, 0x68, 0x86, 0x71, 0x9c, 0x62, 0x1c, 0x6a, 0x1c, 0xea, 0x85, 0x44, 0xf6, 0xcf, 0x21,
	0x5e, 0xd3, 0x3f, 0xe6, 0xb9, 0x34, 0xde, 0x5d, 0x93, 0x06, 0x0b, 0x89, 0x41, 0x23, 0xfc, 0x0a,
	0x69, 0x0b, 0xea, 0x5b, 0x93, 0x24, 0x62, 0x45, 0x93, 0x61, 0x6d, 0xce, 0xfa, 0x11, 0x65, 0x50,
	0x91, 0xe1, 0x62, 0x2c, 0xc2, 0x88, 0xfd, 0xcf, 0x1a, 0xac, 0x20, 0xbd, 0x82, 0xda, 0xa6, 0x39,
	0x2b, 0x39, 0xc4, 0xb3, 0x7c, 0x9b, 0x67, 0xdb, 0xf8, 0x7a, 0x1b, 0x1f, 0xac, 0x19, 0xc1, 0x96,
	0xec, 0xbe, 0x40, 0xeb, 0x2e, 0x51, 0xda, 0x6c, 0x54, 0x37, 0x22, 0xcd, 0x77, 0xc4, 0xd7, 0x19,
	0x2a, 0xfd, 0xcd, 0x40, 0xb2, 0xf7, 0x40, 0x37, 0x00, 0xf7, 0xb7, 0x24, 0x25, 0xc5, 0x54, 0x21,
	0x3d, 0x87, 0x72, 0xe6, 0xc9, 0x88, 0x53, 0xf4, 0x2c, 0xff, 0x70, 0x27, 0xc2, 0xc8, 0x83, 0x15,
	0xe7, 0xec, 0x04, 0x1a, 0x39, 0x1b, 0x5a, 0x87, 0x6a, 0xaf, 0x7b, 0xff, 0xd8, 0xed, 0xdf, 0x3e,
	0x1c, 0x14, 0xfc, 0x77, 0x02, 0x60, 0x04, 0x4f, 0x4b, 0x39, 0x7d, 0x23, 0x60, 0xff, 0x5c, 0x81,
	0x5e, 0xec, 0x44, 0xfd, 0xf9, 0x2e, 0x76, 0x67, 0x0f, 0x45, 0xb6, 0xcf, 0x2d, 0x8c, 0xca, 0xe6,
	0x43, 0x2e, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x69, 0xd2, 0x62, 0x8e, 0x98, 0x02, 0x00, 0x00,
}
