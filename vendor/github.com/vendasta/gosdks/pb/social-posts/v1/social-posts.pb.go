// Code generated by protoc-gen-go.
// source: social-posts.proto
// DO NOT EDIT!

/*
Package socialposts_v1 is a generated protocol buffer package.

It is generated from these files:
	social-posts.proto

It has these top-level messages:
	SocialPost
	ListSocialPostResponse
	ListSocialPostRequest
*/
package socialposts_v1

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

type SocialPost_DeletionStatus int32

const (
	SocialPost_NONE        SocialPost_DeletionStatus = 0
	SocialPost_FAILED      SocialPost_DeletionStatus = 1
	SocialPost_IN_PROGRESS SocialPost_DeletionStatus = 2
)

var SocialPost_DeletionStatus_name = map[int32]string{
	0: "NONE",
	1: "FAILED",
	2: "IN_PROGRESS",
}
var SocialPost_DeletionStatus_value = map[string]int32{
	"NONE":        0,
	"FAILED":      1,
	"IN_PROGRESS": 2,
}

func (x SocialPost_DeletionStatus) String() string {
	return proto.EnumName(SocialPost_DeletionStatus_name, int32(x))
}
func (SocialPost_DeletionStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type SocialPost_Service int32

const (
	SocialPost_TWITTER     SocialPost_Service = 0
	SocialPost_FACEBOOK    SocialPost_Service = 1
	SocialPost_LINKED_IN   SocialPost_Service = 2
	SocialPost_GOOGLE_PLUS SocialPost_Service = 3
)

var SocialPost_Service_name = map[int32]string{
	0: "TWITTER",
	1: "FACEBOOK",
	2: "LINKED_IN",
	3: "GOOGLE_PLUS",
}
var SocialPost_Service_value = map[string]int32{
	"TWITTER":     0,
	"FACEBOOK":    1,
	"LINKED_IN":   2,
	"GOOGLE_PLUS": 3,
}

func (x SocialPost_Service) String() string {
	return proto.EnumName(SocialPost_Service_name, int32(x))
}
func (SocialPost_Service) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

type SocialPost struct {
	AccountId      string                     `protobuf:"bytes,1,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
	SocialPostId   string                     `protobuf:"bytes,2,opt,name=social_post_id,json=socialPostId" json:"social_post_id,omitempty"`
	PostText       string                     `protobuf:"bytes,3,opt,name=post_text,json=postText" json:"post_text,omitempty"`
	Posted         *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=posted" json:"posted,omitempty"`
	IsError        bool                       `protobuf:"varint,5,opt,name=is_error,json=isError" json:"is_error,omitempty"`
	DeletionStatus SocialPost_DeletionStatus  `protobuf:"varint,6,opt,name=deletion_status,json=deletionStatus,enum=socialposts.v1.SocialPost_DeletionStatus" json:"deletion_status,omitempty"`
	Service        SocialPost_Service         `protobuf:"varint,7,opt,name=service,enum=socialposts.v1.SocialPost_Service" json:"service,omitempty"`
}

func (m *SocialPost) Reset()                    { *m = SocialPost{} }
func (m *SocialPost) String() string            { return proto.CompactTextString(m) }
func (*SocialPost) ProtoMessage()               {}
func (*SocialPost) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SocialPost) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *SocialPost) GetSocialPostId() string {
	if m != nil {
		return m.SocialPostId
	}
	return ""
}

func (m *SocialPost) GetPostText() string {
	if m != nil {
		return m.PostText
	}
	return ""
}

func (m *SocialPost) GetPosted() *google_protobuf.Timestamp {
	if m != nil {
		return m.Posted
	}
	return nil
}

func (m *SocialPost) GetIsError() bool {
	if m != nil {
		return m.IsError
	}
	return false
}

func (m *SocialPost) GetDeletionStatus() SocialPost_DeletionStatus {
	if m != nil {
		return m.DeletionStatus
	}
	return SocialPost_NONE
}

func (m *SocialPost) GetService() SocialPost_Service {
	if m != nil {
		return m.Service
	}
	return SocialPost_TWITTER
}

type ListSocialPostResponse struct {
	SocialPosts []*SocialPost `protobuf:"bytes,1,rep,name=social_posts,json=socialPosts" json:"social_posts,omitempty"`
	NextCursor  string        `protobuf:"bytes,2,opt,name=next_cursor,json=nextCursor" json:"next_cursor,omitempty"`
	HasMore     bool          `protobuf:"varint,3,opt,name=has_more,json=hasMore" json:"has_more,omitempty"`
}

func (m *ListSocialPostResponse) Reset()                    { *m = ListSocialPostResponse{} }
func (m *ListSocialPostResponse) String() string            { return proto.CompactTextString(m) }
func (*ListSocialPostResponse) ProtoMessage()               {}
func (*ListSocialPostResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ListSocialPostResponse) GetSocialPosts() []*SocialPost {
	if m != nil {
		return m.SocialPosts
	}
	return nil
}

func (m *ListSocialPostResponse) GetNextCursor() string {
	if m != nil {
		return m.NextCursor
	}
	return ""
}

func (m *ListSocialPostResponse) GetHasMore() bool {
	if m != nil {
		return m.HasMore
	}
	return false
}

type ListSocialPostRequest struct {
	Start      *google_protobuf.Timestamp `protobuf:"bytes,1,opt,name=start" json:"start,omitempty"`
	End        *google_protobuf.Timestamp `protobuf:"bytes,2,opt,name=end" json:"end,omitempty"`
	AccountId  string                     `protobuf:"bytes,3,opt,name=account_id,json=accountId" json:"account_id,omitempty"`
	PartnerId  string                     `protobuf:"bytes,4,opt,name=partner_id,json=partnerId" json:"partner_id,omitempty"`
	NextCursor string                     `protobuf:"bytes,5,opt,name=next_cursor,json=nextCursor" json:"next_cursor,omitempty"`
}

func (m *ListSocialPostRequest) Reset()                    { *m = ListSocialPostRequest{} }
func (m *ListSocialPostRequest) String() string            { return proto.CompactTextString(m) }
func (*ListSocialPostRequest) ProtoMessage()               {}
func (*ListSocialPostRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ListSocialPostRequest) GetStart() *google_protobuf.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *ListSocialPostRequest) GetEnd() *google_protobuf.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

func (m *ListSocialPostRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *ListSocialPostRequest) GetPartnerId() string {
	if m != nil {
		return m.PartnerId
	}
	return ""
}

func (m *ListSocialPostRequest) GetNextCursor() string {
	if m != nil {
		return m.NextCursor
	}
	return ""
}

func init() {
	proto.RegisterType((*SocialPost)(nil), "socialposts.v1.SocialPost")
	proto.RegisterType((*ListSocialPostResponse)(nil), "socialposts.v1.ListSocialPostResponse")
	proto.RegisterType((*ListSocialPostRequest)(nil), "socialposts.v1.ListSocialPostRequest")
	proto.RegisterEnum("socialposts.v1.SocialPost_DeletionStatus", SocialPost_DeletionStatus_name, SocialPost_DeletionStatus_value)
	proto.RegisterEnum("socialposts.v1.SocialPost_Service", SocialPost_Service_name, SocialPost_Service_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SocialPosts service

type SocialPostsClient interface {
	List(ctx context.Context, in *ListSocialPostRequest, opts ...grpc.CallOption) (*ListSocialPostResponse, error)
}

type socialPostsClient struct {
	cc *grpc.ClientConn
}

func NewSocialPostsClient(cc *grpc.ClientConn) SocialPostsClient {
	return &socialPostsClient{cc}
}

func (c *socialPostsClient) List(ctx context.Context, in *ListSocialPostRequest, opts ...grpc.CallOption) (*ListSocialPostResponse, error) {
	out := new(ListSocialPostResponse)
	err := grpc.Invoke(ctx, "/socialposts.v1.SocialPosts/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SocialPosts service

type SocialPostsServer interface {
	List(context.Context, *ListSocialPostRequest) (*ListSocialPostResponse, error)
}

func RegisterSocialPostsServer(s *grpc.Server, srv SocialPostsServer) {
	s.RegisterService(&_SocialPosts_serviceDesc, srv)
}

func _SocialPosts_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSocialPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SocialPostsServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/socialposts.v1.SocialPosts/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SocialPostsServer).List(ctx, req.(*ListSocialPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SocialPosts_serviceDesc = grpc.ServiceDesc{
	ServiceName: "socialposts.v1.SocialPosts",
	HandlerType: (*SocialPostsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _SocialPosts_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "social-posts.proto",
}

func init() { proto.RegisterFile("social-posts.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x52, 0x41, 0x6f, 0xd3, 0x4c,
	0x10, 0xad, 0xe3, 0x24, 0x76, 0xc6, 0xf9, 0x52, 0x6b, 0xa5, 0x0f, 0x99, 0xa0, 0xaa, 0x91, 0x05,
	0x28, 0x48, 0xe0, 0x42, 0x38, 0x70, 0x81, 0x43, 0x69, 0xdc, 0xc8, 0x6a, 0x88, 0xa3, 0x75, 0x50,
	0x8f, 0x96, 0x1b, 0x6f, 0x5b, 0x4b, 0x89, 0x37, 0xec, 0x6e, 0xaa, 0xfc, 0x0f, 0xfe, 0x17, 0x47,
	0x7e, 0x0f, 0xda, 0xb5, 0x43, 0x1a, 0x0b, 0xc2, 0xcd, 0x7e, 0xf3, 0xe6, 0xed, 0xbc, 0x79, 0x03,
	0x88, 0xd3, 0x79, 0x96, 0x2c, 0xde, 0xac, 0x28, 0x17, 0xdc, 0x5b, 0x31, 0x2a, 0x28, 0xea, 0x14,
	0x58, 0x01, 0x3d, 0xbc, 0xeb, 0x9e, 0xde, 0x51, 0x7a, 0xb7, 0x20, 0x67, 0xaa, 0x7a, 0xb3, 0xbe,
	0x3d, 0x13, 0xd9, 0x92, 0x70, 0x91, 0x2c, 0x57, 0x45, 0x83, 0xfb, 0x43, 0x07, 0x88, 0x54, 0xcf,
	0x94, 0x72, 0x81, 0x4e, 0x00, 0x92, 0xf9, 0x9c, 0xae, 0x73, 0x11, 0x67, 0xa9, 0xa3, 0xf5, 0xb4,
	0x7e, 0x0b, 0xb7, 0x4a, 0x24, 0x48, 0xd1, 0x73, 0x28, 0x1f, 0x88, 0xe5, 0x0b, 0x92, 0x52, 0x53,
	0x94, 0x36, 0xff, 0x2d, 0x11, 0xa4, 0xe8, 0x19, 0xb4, 0x54, 0x59, 0x90, 0x8d, 0x70, 0x74, 0x45,
	0x30, 0x25, 0x30, 0x23, 0x1b, 0x81, 0x06, 0xd0, 0x94, 0xdf, 0x24, 0x75, 0xea, 0x3d, 0xad, 0x6f,
	0x0d, 0xba, 0x5e, 0x31, 0xa2, 0xb7, 0x1d, 0xd1, 0x9b, 0x6d, 0x47, 0xc4, 0x25, 0x13, 0x3d, 0x05,
	0x33, 0xe3, 0x31, 0x61, 0x8c, 0x32, 0xa7, 0xd1, 0xd3, 0xfa, 0x26, 0x36, 0x32, 0xee, 0xcb, 0x5f,
	0x84, 0xe1, 0x38, 0x25, 0x0b, 0x22, 0x32, 0x9a, 0xc7, 0x5c, 0x24, 0x62, 0xcd, 0x9d, 0x66, 0x4f,
	0xeb, 0x77, 0x06, 0xaf, 0xbc, 0xfd, 0x55, 0x78, 0x3b, 0x97, 0xde, 0xb0, 0xec, 0x88, 0x54, 0x03,
	0xee, 0xa4, 0x7b, 0xff, 0xe8, 0x23, 0x18, 0x9c, 0xb0, 0x87, 0x6c, 0x4e, 0x1c, 0x43, 0x69, 0xb9,
	0x07, 0xb4, 0xa2, 0x82, 0x89, 0xb7, 0x2d, 0xee, 0x07, 0xe8, 0xec, 0xeb, 0x23, 0x13, 0xea, 0x93,
	0x70, 0xe2, 0xdb, 0x47, 0x08, 0xa0, 0x79, 0x79, 0x1e, 0x8c, 0xfd, 0xa1, 0xad, 0xa1, 0x63, 0xb0,
	0x82, 0x49, 0x3c, 0xc5, 0xe1, 0x08, 0xfb, 0x51, 0x64, 0xd7, 0xdc, 0x21, 0x18, 0xa5, 0x18, 0xb2,
	0xc0, 0x98, 0x5d, 0x07, 0xb3, 0x99, 0x8f, 0xed, 0x23, 0xd4, 0x06, 0xf3, 0xf2, 0xfc, 0xc2, 0xff,
	0x1c, 0x86, 0x57, 0xb6, 0x86, 0xfe, 0x83, 0xd6, 0x38, 0x98, 0x5c, 0xf9, 0xc3, 0x38, 0x98, 0xd8,
	0x35, 0xa9, 0x32, 0x0a, 0xc3, 0xd1, 0xd8, 0x8f, 0xa7, 0xe3, 0xaf, 0x91, 0xad, 0xbb, 0xdf, 0x35,
	0x78, 0x32, 0xce, 0xb8, 0xd8, 0x8d, 0x88, 0x09, 0x5f, 0xd1, 0x9c, 0x13, 0xf4, 0x09, 0xda, 0x8f,
	0xd2, 0xe3, 0x8e, 0xd6, 0xd3, 0x55, 0x00, 0x7f, 0x35, 0x87, 0xad, 0x5d, 0xae, 0x1c, 0x9d, 0x82,
	0x95, 0x93, 0x8d, 0x88, 0xe7, 0x6b, 0xc6, 0x29, 0x2b, 0x93, 0x07, 0x09, 0x5d, 0x28, 0x44, 0xc6,
	0x74, 0x9f, 0xf0, 0x78, 0x49, 0x19, 0x51, 0xb1, 0x9b, 0xd8, 0xb8, 0x4f, 0xf8, 0x17, 0xca, 0x88,
	0xfb, 0x53, 0x83, 0xff, 0xab, 0x53, 0x7d, 0x5b, 0x13, 0x2e, 0xd0, 0x5b, 0x68, 0x70, 0x91, 0x30,
	0xa1, 0x8e, 0xed, 0xf0, 0x39, 0x14, 0x44, 0xf4, 0x1a, 0x74, 0x92, 0x17, 0x97, 0x77, 0x98, 0x2f,
	0x69, 0x95, 0x8b, 0xd6, 0xab, 0x17, 0x7d, 0x02, 0xb0, 0x4a, 0x98, 0xc8, 0x09, 0x93, 0xe5, 0x7a,
	0x51, 0x2e, 0x91, 0x20, 0xad, 0x7a, 0x6e, 0x54, 0x3d, 0x0f, 0x6e, 0xc1, 0x8a, 0x1e, 0xed, 0xe8,
	0x1a, 0xea, 0xd2, 0x26, 0x7a, 0x51, 0x5d, 0xea, 0x1f, 0xcd, 0x77, 0x5f, 0xfe, 0x8b, 0x56, 0x24,
	0xe7, 0x1e, 0xdd, 0x34, 0x95, 0xbf, 0xf7, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x67, 0x0a, 0x73,
	0x5e, 0xf5, 0x03, 0x00, 0x00,
}
