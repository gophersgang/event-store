// Code generated by protoc-gen-go.
// source: citations.proto
// DO NOT EDIT!

/*
Package citations is a generated protocol buffer package.

It is generated from these files:
	citations.proto

It has these top-level messages:
	Citation
	RangeCitationsRequest
	RangeCitationsResponse
*/
package citations

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

type Citation struct {
	AccountId string                     `protobuf:"bytes,1,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
	Url       string                     `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
	Title     string                     `protobuf:"bytes,3,opt,name=title" json:"title,omitempty"`
	Created   *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=created" json:"created,omitempty"`
}

func (m *Citation) Reset()                    { *m = Citation{} }
func (m *Citation) String() string            { return proto.CompactTextString(m) }
func (*Citation) ProtoMessage()               {}
func (*Citation) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Citation) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *Citation) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Citation) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Citation) GetCreated() *google_protobuf.Timestamp {
	if m != nil {
		return m.Created
	}
	return nil
}

type RangeCitationsRequest struct {
	AccountId string                     `protobuf:"bytes,1,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
	StartTime *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime" json:"start_time,omitempty"`
	EndTime   *google_protobuf.Timestamp `protobuf:"bytes,3,opt,name=end_time,json=endTime" json:"end_time,omitempty"`
}

func (m *RangeCitationsRequest) Reset()                    { *m = RangeCitationsRequest{} }
func (m *RangeCitationsRequest) String() string            { return proto.CompactTextString(m) }
func (*RangeCitationsRequest) ProtoMessage()               {}
func (*RangeCitationsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RangeCitationsRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *RangeCitationsRequest) GetStartTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *RangeCitationsRequest) GetEndTime() *google_protobuf.Timestamp {
	if m != nil {
		return m.EndTime
	}
	return nil
}

type RangeCitationsResponse struct {
	Citations []*Citation `protobuf:"bytes,1,rep,name=citations" json:"citations,omitempty"`
}

func (m *RangeCitationsResponse) Reset()                    { *m = RangeCitationsResponse{} }
func (m *RangeCitationsResponse) String() string            { return proto.CompactTextString(m) }
func (*RangeCitationsResponse) ProtoMessage()               {}
func (*RangeCitationsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RangeCitationsResponse) GetCitations() []*Citation {
	if m != nil {
		return m.Citations
	}
	return nil
}

func init() {
	proto.RegisterType((*Citation)(nil), "citations.Citation")
	proto.RegisterType((*RangeCitationsRequest)(nil), "citations.RangeCitationsRequest")
	proto.RegisterType((*RangeCitationsResponse)(nil), "citations.RangeCitationsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CitationService service

type CitationServiceClient interface {
	GetByRange(ctx context.Context, in *RangeCitationsRequest, opts ...grpc.CallOption) (*RangeCitationsResponse, error)
}

type citationServiceClient struct {
	cc *grpc.ClientConn
}

func NewCitationServiceClient(cc *grpc.ClientConn) CitationServiceClient {
	return &citationServiceClient{cc}
}

func (c *citationServiceClient) GetByRange(ctx context.Context, in *RangeCitationsRequest, opts ...grpc.CallOption) (*RangeCitationsResponse, error) {
	out := new(RangeCitationsResponse)
	err := grpc.Invoke(ctx, "/citations.CitationService/GetByRange", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CitationService service

type CitationServiceServer interface {
	GetByRange(context.Context, *RangeCitationsRequest) (*RangeCitationsResponse, error)
}

func RegisterCitationServiceServer(s *grpc.Server, srv CitationServiceServer) {
	s.RegisterService(&_CitationService_serviceDesc, srv)
}

func _CitationService_GetByRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RangeCitationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CitationServiceServer).GetByRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/citations.CitationService/GetByRange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CitationServiceServer).GetByRange(ctx, req.(*RangeCitationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CitationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "citations.CitationService",
	HandlerType: (*CitationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByRange",
			Handler:    _CitationService_GetByRange_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "citations.proto",
}

func init() { proto.RegisterFile("citations.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x91, 0xcd, 0x4e, 0xeb, 0x30,
	0x10, 0x85, 0xe5, 0x9b, 0x0b, 0x34, 0xd3, 0x45, 0x91, 0xf9, 0x51, 0x54, 0x09, 0x11, 0xba, 0xca,
	0x2a, 0x15, 0x01, 0x16, 0x6c, 0x61, 0x81, 0x10, 0x2b, 0x02, 0xfb, 0xca, 0x8d, 0x87, 0xc8, 0x52,
	0x6a, 0x87, 0x78, 0x82, 0xc4, 0x13, 0xf0, 0x2a, 0x3c, 0x26, 0x8a, 0x53, 0x13, 0x84, 0x50, 0xbb,
	0xf3, 0x39, 0x3e, 0xe3, 0xf9, 0x66, 0x0c, 0x93, 0x42, 0x91, 0x20, 0x65, 0xb4, 0x4d, 0xeb, 0xc6,
	0x90, 0xe1, 0xe1, 0xb7, 0x31, 0x3d, 0x2d, 0x8d, 0x29, 0x2b, 0x9c, 0xbb, 0x8b, 0x65, 0xfb, 0x32,
	0x27, 0xb5, 0x42, 0x4b, 0x62, 0x55, 0xf7, 0xd9, 0xd9, 0x07, 0x83, 0xd1, 0xed, 0x3a, 0xce, 0x4f,
	0x00, 0x44, 0x51, 0x98, 0x56, 0xd3, 0x42, 0xc9, 0x88, 0xc5, 0x2c, 0x09, 0xf3, 0x70, 0xed, 0xdc,
	0x4b, 0xbe, 0x0f, 0x41, 0xdb, 0x54, 0xd1, 0x3f, 0xe7, 0x77, 0x47, 0x7e, 0x08, 0x3b, 0xa4, 0xa8,
	0xc2, 0x28, 0x70, 0x5e, 0x2f, 0xf8, 0x25, 0xec, 0x15, 0x0d, 0x0a, 0x42, 0x19, 0xfd, 0x8f, 0x59,
	0x32, 0xce, 0xa6, 0x69, 0x8f, 0x91, 0x7a, 0x8c, 0xf4, 0xd9, 0x63, 0xe4, 0x3e, 0x3a, 0xfb, 0x64,
	0x70, 0x94, 0x0b, 0x5d, 0xa2, 0xc7, 0xb1, 0x39, 0xbe, 0xb6, 0x68, 0x69, 0x1b, 0xd6, 0x35, 0x80,
	0x25, 0xd1, 0xd0, 0xa2, 0x9b, 0xcd, 0xd1, 0x6d, 0xee, 0x18, 0xba, 0x74, 0xa7, 0xf9, 0x15, 0x8c,
	0x50, 0xcb, 0xbe, 0x30, 0xd8, 0x8e, 0x8a, 0x5a, 0x76, 0x6a, 0xf6, 0x00, 0xc7, 0xbf, 0x49, 0x6d,
	0x6d, 0xb4, 0x45, 0x7e, 0x0e, 0xc3, 0xf2, 0x23, 0x16, 0x07, 0xc9, 0x38, 0x3b, 0x48, 0x87, 0xff,
	0xf1, 0x05, 0xf9, 0x90, 0xca, 0x24, 0x4c, 0xbc, 0xfd, 0x84, 0xcd, 0x9b, 0x2a, 0x90, 0x3f, 0x02,
	0xdc, 0x21, 0xdd, 0xbc, 0xbb, 0x26, 0x3c, 0xfe, 0xf1, 0xc0, 0x9f, 0x0b, 0x9a, 0x9e, 0x6d, 0x48,
	0xf4, 0x60, 0xcb, 0x5d, 0x37, 0xcf, 0xc5, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xce, 0x6d, 0x07,
	0x41, 0x2d, 0x02, 0x00, 0x00,
}
