// Code generated by protoc-gen-go.
// source: review.proto
// DO NOT EDIT!

package vendasta_listingsproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// based on Review VObject
type Review struct {
	ReviewId          string                      `protobuf:"bytes,1,opt,name=review_id,json=reviewId" json:"review_id,omitempty"`
	ListingId         string                      `protobuf:"bytes,2,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	Url               string                      `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	StarRating        float32                     `protobuf:"fixed32,4,opt,name=star_rating,json=starRating" json:"star_rating,omitempty"`
	ReviewerName      string                      `protobuf:"bytes,5,opt,name=reviewer_name,json=reviewerName" json:"reviewer_name,omitempty"`
	ReviewerEmail     string                      `protobuf:"bytes,6,opt,name=reviewer_email,json=reviewerEmail" json:"reviewer_email,omitempty"`
	ReviewerUrl       string                      `protobuf:"bytes,7,opt,name=reviewer_url,json=reviewerUrl" json:"reviewer_url,omitempty"`
	Content           string                      `protobuf:"bytes,8,opt,name=content" json:"content,omitempty"`
	PublishedDate     *google_protobuf1.Timestamp `protobuf:"bytes,9,opt,name=published_date,json=publishedDate" json:"published_date,omitempty"`
	Title             string                      `protobuf:"bytes,10,opt,name=title" json:"title,omitempty"`
	ListingExternalId string                      `protobuf:"bytes,11,opt,name=listing_external_id,json=listingExternalId" json:"listing_external_id,omitempty"`
}

func (m *Review) Reset()                    { *m = Review{} }
func (m *Review) String() string            { return proto.CompactTextString(m) }
func (*Review) ProtoMessage()               {}
func (*Review) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *Review) GetReviewId() string {
	if m != nil {
		return m.ReviewId
	}
	return ""
}

func (m *Review) GetListingId() string {
	if m != nil {
		return m.ListingId
	}
	return ""
}

func (m *Review) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Review) GetStarRating() float32 {
	if m != nil {
		return m.StarRating
	}
	return 0
}

func (m *Review) GetReviewerName() string {
	if m != nil {
		return m.ReviewerName
	}
	return ""
}

func (m *Review) GetReviewerEmail() string {
	if m != nil {
		return m.ReviewerEmail
	}
	return ""
}

func (m *Review) GetReviewerUrl() string {
	if m != nil {
		return m.ReviewerUrl
	}
	return ""
}

func (m *Review) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func (m *Review) GetPublishedDate() *google_protobuf1.Timestamp {
	if m != nil {
		return m.PublishedDate
	}
	return nil
}

func (m *Review) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Review) GetListingExternalId() string {
	if m != nil {
		return m.ListingExternalId
	}
	return ""
}

type GetReviewRequest struct {
	ReviewId string `protobuf:"bytes,1,opt,name=review_id,json=reviewId" json:"review_id,omitempty"`
}

func (m *GetReviewRequest) Reset()                    { *m = GetReviewRequest{} }
func (m *GetReviewRequest) String() string            { return proto.CompactTextString(m) }
func (*GetReviewRequest) ProtoMessage()               {}
func (*GetReviewRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *GetReviewRequest) GetReviewId() string {
	if m != nil {
		return m.ReviewId
	}
	return ""
}

type DeleteReviewRequest struct {
	ReviewId string `protobuf:"bytes,1,opt,name=review_id,json=reviewId" json:"review_id,omitempty"`
}

func (m *DeleteReviewRequest) Reset()                    { *m = DeleteReviewRequest{} }
func (m *DeleteReviewRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteReviewRequest) ProtoMessage()               {}
func (*DeleteReviewRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *DeleteReviewRequest) GetReviewId() string {
	if m != nil {
		return m.ReviewId
	}
	return ""
}

type ListReviewsRequest struct {
	ListingId         string `protobuf:"bytes,1,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	ListingExternalId string `protobuf:"bytes,2,opt,name=listing_external_id,json=listingExternalId" json:"listing_external_id,omitempty"`
	// int64 offset = 3;
	PageSize int64  `protobuf:"varint,4,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
	Cursor   string `protobuf:"bytes,5,opt,name=cursor" json:"cursor,omitempty"`
}

func (m *ListReviewsRequest) Reset()                    { *m = ListReviewsRequest{} }
func (m *ListReviewsRequest) String() string            { return proto.CompactTextString(m) }
func (*ListReviewsRequest) ProtoMessage()               {}
func (*ListReviewsRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *ListReviewsRequest) GetListingId() string {
	if m != nil {
		return m.ListingId
	}
	return ""
}

func (m *ListReviewsRequest) GetListingExternalId() string {
	if m != nil {
		return m.ListingExternalId
	}
	return ""
}

func (m *ListReviewsRequest) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListReviewsRequest) GetCursor() string {
	if m != nil {
		return m.Cursor
	}
	return ""
}

type ListReviewsResponse struct {
	Reviews              []*Review `protobuf:"bytes,1,rep,name=reviews" json:"reviews,omitempty"`
	TotalNumberOfReviews int64     `protobuf:"varint,2,opt,name=total_number_of_reviews,json=totalNumberOfReviews" json:"total_number_of_reviews,omitempty"`
	// int64 offset = 3;
	PageSize int64  `protobuf:"varint,4,opt,name=page_size,json=pageSize" json:"page_size,omitempty"`
	Cursor   string `protobuf:"bytes,5,opt,name=cursor" json:"cursor,omitempty"`
}

func (m *ListReviewsResponse) Reset()                    { *m = ListReviewsResponse{} }
func (m *ListReviewsResponse) String() string            { return proto.CompactTextString(m) }
func (*ListReviewsResponse) ProtoMessage()               {}
func (*ListReviewsResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *ListReviewsResponse) GetReviews() []*Review {
	if m != nil {
		return m.Reviews
	}
	return nil
}

func (m *ListReviewsResponse) GetTotalNumberOfReviews() int64 {
	if m != nil {
		return m.TotalNumberOfReviews
	}
	return 0
}

func (m *ListReviewsResponse) GetPageSize() int64 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *ListReviewsResponse) GetCursor() string {
	if m != nil {
		return m.Cursor
	}
	return ""
}

func init() {
	proto.RegisterType((*Review)(nil), "vendasta.listingsproto.Review")
	proto.RegisterType((*GetReviewRequest)(nil), "vendasta.listingsproto.GetReviewRequest")
	proto.RegisterType((*DeleteReviewRequest)(nil), "vendasta.listingsproto.DeleteReviewRequest")
	proto.RegisterType((*ListReviewsRequest)(nil), "vendasta.listingsproto.ListReviewsRequest")
	proto.RegisterType((*ListReviewsResponse)(nil), "vendasta.listingsproto.ListReviewsResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ReviewService service

type ReviewServiceClient interface {
	Put(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Review, error)
	Get(ctx context.Context, in *GetReviewRequest, opts ...grpc.CallOption) (*Review, error)
	Delete(ctx context.Context, in *DeleteReviewRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	List(ctx context.Context, in *ListReviewsRequest, opts ...grpc.CallOption) (*ListReviewsResponse, error)
}

type reviewServiceClient struct {
	cc *grpc.ClientConn
}

func NewReviewServiceClient(cc *grpc.ClientConn) ReviewServiceClient {
	return &reviewServiceClient{cc}
}

func (c *reviewServiceClient) Put(ctx context.Context, in *Review, opts ...grpc.CallOption) (*Review, error) {
	out := new(Review)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ReviewService/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) Get(ctx context.Context, in *GetReviewRequest, opts ...grpc.CallOption) (*Review, error) {
	out := new(Review)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ReviewService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) Delete(ctx context.Context, in *DeleteReviewRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ReviewService/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) List(ctx context.Context, in *ListReviewsRequest, opts ...grpc.CallOption) (*ListReviewsResponse, error) {
	out := new(ListReviewsResponse)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ReviewService/List", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ReviewService service

type ReviewServiceServer interface {
	Put(context.Context, *Review) (*Review, error)
	Get(context.Context, *GetReviewRequest) (*Review, error)
	Delete(context.Context, *DeleteReviewRequest) (*google_protobuf.Empty, error)
	List(context.Context, *ListReviewsRequest) (*ListReviewsResponse, error)
}

func RegisterReviewServiceServer(s *grpc.Server, srv ReviewServiceServer) {
	s.RegisterService(&_ReviewService_serviceDesc, srv)
}

func _ReviewService_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Review)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ReviewService/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).Put(ctx, req.(*Review))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ReviewService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).Get(ctx, req.(*GetReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ReviewService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).Delete(ctx, req.(*DeleteReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListReviewsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ReviewService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).List(ctx, req.(*ListReviewsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReviewService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vendasta.listingsproto.ReviewService",
	HandlerType: (*ReviewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _ReviewService_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ReviewService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ReviewService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ReviewService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "review.proto",
}

func init() { proto.RegisterFile("review.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 547 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0xed, 0x6a, 0x13, 0x41,
	0x14, 0x65, 0xb3, 0x6d, 0x3e, 0x6e, 0x9a, 0x52, 0x6f, 0x4a, 0x5c, 0xb6, 0x68, 0x63, 0x44, 0x08,
	0x16, 0x36, 0x10, 0x11, 0xfc, 0x2b, 0x34, 0x94, 0x80, 0xb6, 0xb2, 0xd5, 0xdf, 0xcb, 0x24, 0x7b,
	0x13, 0x07, 0xf6, 0xcb, 0x99, 0xd9, 0xa8, 0x7d, 0x03, 0xdf, 0xc0, 0x77, 0xf1, 0x41, 0x7c, 0x1d,
	0xd9, 0x99, 0xdd, 0x40, 0xd3, 0xc6, 0x68, 0xff, 0xe5, 0x9e, 0x73, 0xee, 0xc7, 0xde, 0x73, 0x27,
	0x70, 0x20, 0x68, 0xc5, 0xe9, 0xab, 0x97, 0x89, 0x54, 0xa5, 0xd8, 0x5b, 0x51, 0x12, 0x32, 0xa9,
	0x98, 0x17, 0x71, 0xa9, 0x78, 0xb2, 0x94, 0x1a, 0x77, 0x4f, 0x96, 0x69, 0xba, 0x8c, 0x68, 0xa4,
	0xa3, 0x59, 0xbe, 0x18, 0x51, 0x9c, 0xa9, 0xef, 0x26, 0xc9, 0x3d, 0xdd, 0x24, 0x15, 0x8f, 0x49,
	0x2a, 0x16, 0x67, 0x46, 0x30, 0xf8, 0x61, 0x43, 0xdd, 0xd7, 0x6d, 0xf0, 0x04, 0x5a, 0xa6, 0x61,
	0xc0, 0x43, 0xc7, 0xea, 0x5b, 0xc3, 0x96, 0xdf, 0x34, 0xc0, 0x34, 0xc4, 0x27, 0x00, 0x65, 0xdb,
	0x82, 0xad, 0x69, 0xb6, 0x55, 0x22, 0xd3, 0x10, 0x8f, 0xc0, 0xce, 0x45, 0xe4, 0xd8, 0x1a, 0x2f,
	0x7e, 0xe2, 0x29, 0xb4, 0xa5, 0x62, 0x22, 0x10, 0xac, 0x90, 0x38, 0x7b, 0x7d, 0x6b, 0x58, 0xf3,
	0xa1, 0x80, 0x7c, 0x8d, 0xe0, 0x73, 0xe8, 0x98, 0xea, 0x24, 0x82, 0x84, 0xc5, 0xe4, 0xec, 0xeb,
	0xe4, 0x83, 0x0a, 0xbc, 0x64, 0x31, 0xe1, 0x0b, 0x38, 0x5c, 0x8b, 0x28, 0x66, 0x3c, 0x72, 0xea,
	0x5a, 0xb5, 0x4e, 0x9d, 0x14, 0x20, 0x3e, 0x83, 0x75, 0x5a, 0x50, 0xcc, 0xd1, 0xd0, 0xa2, 0x76,
	0x85, 0x7d, 0x12, 0x11, 0x3a, 0xd0, 0x98, 0xa7, 0x89, 0xa2, 0x44, 0x39, 0x4d, 0xcd, 0x56, 0x21,
	0xbe, 0x85, 0xc3, 0x2c, 0x9f, 0x45, 0x5c, 0x7e, 0xa6, 0x30, 0x08, 0x99, 0x22, 0xa7, 0xd5, 0xb7,
	0x86, 0xed, 0xb1, 0xeb, 0x99, 0xe5, 0x79, 0xd5, 0xf2, 0xbc, 0x8f, 0xd5, 0xf2, 0xfc, 0xce, 0x3a,
	0xe3, 0x9c, 0x29, 0xc2, 0x63, 0xd8, 0x57, 0x5c, 0x45, 0xe4, 0x80, 0x2e, 0x6d, 0x02, 0xf4, 0xa0,
	0x5b, 0xed, 0x8c, 0xbe, 0x29, 0x12, 0x09, 0x8b, 0x8a, 0xe5, 0xb5, 0xb5, 0xe6, 0x51, 0x49, 0x4d,
	0x4a, 0x66, 0x1a, 0x0e, 0x46, 0x70, 0x74, 0x41, 0xca, 0xb8, 0xe1, 0xd3, 0x97, 0x9c, 0xa4, 0xfa,
	0xab, 0x29, 0x83, 0x31, 0x74, 0xcf, 0x29, 0x22, 0x45, 0xff, 0x91, 0xf3, 0xd3, 0x02, 0x7c, 0xc7,
	0x65, 0xd9, 0x46, 0x56, 0x39, 0xb7, 0xfd, 0xb5, 0x36, 0xfd, 0xdd, 0xf2, 0x29, 0xb5, 0x2d, 0x9f,
	0x52, 0x8c, 0x90, 0xb1, 0x25, 0x05, 0x92, 0xdf, 0x90, 0xf6, 0xde, 0xf6, 0x9b, 0x05, 0x70, 0xcd,
	0x6f, 0x08, 0x7b, 0x50, 0x9f, 0xe7, 0x42, 0xa6, 0xa2, 0xb4, 0xbc, 0x8c, 0x06, 0xbf, 0x2c, 0xe8,
	0xde, 0x1a, 0x4d, 0x66, 0x69, 0x22, 0x09, 0xdf, 0x40, 0xc3, 0x8c, 0x2f, 0x1d, 0xab, 0x6f, 0x0f,
	0xdb, 0xe3, 0xa7, 0xde, 0xfd, 0x6f, 0xc1, 0x2b, 0xf7, 0x50, 0xc9, 0xf1, 0x35, 0x3c, 0x56, 0xa9,
	0x62, 0x51, 0x90, 0xe4, 0xf1, 0x8c, 0x44, 0x90, 0x2e, 0x82, 0xaa, 0x52, 0x4d, 0x0f, 0x75, 0xac,
	0xe9, 0x4b, 0xcd, 0x5e, 0x2d, 0xca, 0xc6, 0x0f, 0x9a, 0x7e, 0xfc, 0xbb, 0x06, 0x1d, 0x53, 0xe0,
	0x9a, 0xc4, 0x8a, 0xcf, 0x09, 0x27, 0x60, 0x7f, 0xc8, 0x15, 0xee, 0x98, 0xd6, 0xdd, 0xc1, 0xe3,
	0x15, 0xd8, 0x17, 0xa4, 0x70, 0xb8, 0x4d, 0xb6, 0x79, 0x33, 0x3b, 0x0b, 0xbe, 0x87, 0xba, 0x39,
	0x1b, 0x3c, 0xdb, 0xa6, 0xbc, 0xe7, 0xac, 0xdc, 0xde, 0x9d, 0xf7, 0x30, 0x29, 0xfe, 0x69, 0x30,
	0x80, 0xbd, 0xc2, 0x35, 0x7c, 0xb9, 0xad, 0xd8, 0xdd, 0x73, 0x73, 0xcf, 0xfe, 0x49, 0x6b, 0xfc,
	0x9f, 0xd5, 0x35, 0xf5, 0xea, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5b, 0x19, 0x0c, 0xe0, 0x10,
	0x05, 0x00, 0x00,
}