// Code generated by protoc-gen-go.
// source: listing.proto
// DO NOT EDIT!

/*
Package vendasta_listingsproto is a generated protocol buffer package.

It is generated from these files:
	listing.proto
	review.proto

It has these top-level messages:
	Geo
	Listing
	GetListingRequest
	DeleteListingRequest
	Review
	GetReviewRequest
	DeleteReviewRequest
	ListReviewsRequest
	ListReviewsResponse
*/
package vendasta_listingsproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

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

type Geo struct {
	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude" json:"longitude,omitempty"`
}

func (m *Geo) Reset()                    { *m = Geo{} }
func (m *Geo) String() string            { return proto.CompactTextString(m) }
func (*Geo) ProtoMessage()               {}
func (*Geo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Geo) GetLatitude() float64 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *Geo) GetLongitude() float64 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

// based on RawListing VObject
// Modified for datariver to remove the number_of_reviews and average_review_rating from the listing information
type Listing struct {
	ListingId  string `protobuf:"bytes,1,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	ExternalId string `protobuf:"bytes,2,opt,name=external_id,json=externalId" json:"external_id,omitempty"`
	Url        string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	// Basic NAP data
	CompanyName            string   `protobuf:"bytes,4,opt,name=company_name,json=companyName" json:"company_name,omitempty"`
	Address                string   `protobuf:"bytes,5,opt,name=address" json:"address,omitempty"`
	City                   string   `protobuf:"bytes,6,opt,name=city" json:"city,omitempty"`
	State                  string   `protobuf:"bytes,7,opt,name=state" json:"state,omitempty"`
	Country                string   `protobuf:"bytes,8,opt,name=country" json:"country,omitempty"`
	ZipCode                string   `protobuf:"bytes,9,opt,name=zip_code,json=zipCode" json:"zip_code,omitempty"`
	Location               *Geo     `protobuf:"bytes,10,opt,name=location" json:"location,omitempty"`
	Phone                  string   `protobuf:"bytes,11,opt,name=phone" json:"phone,omitempty"`
	AdditionalPhoneNumbers []string `protobuf:"bytes,12,rep,name=additional_phone_numbers,json=additionalPhoneNumbers" json:"additional_phone_numbers,omitempty"`
	Website                string   `protobuf:"bytes,13,opt,name=website" json:"website,omitempty"`
	// Extended NAP data
	// int32 number_of_reviews = 14;
	// float average_review_rating = 15;
	BusinessCategories []string `protobuf:"bytes,16,rep,name=business_categories,json=businessCategories" json:"business_categories,omitempty"`
}

func (m *Listing) Reset()                    { *m = Listing{} }
func (m *Listing) String() string            { return proto.CompactTextString(m) }
func (*Listing) ProtoMessage()               {}
func (*Listing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Listing) GetListingId() string {
	if m != nil {
		return m.ListingId
	}
	return ""
}

func (m *Listing) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func (m *Listing) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *Listing) GetCompanyName() string {
	if m != nil {
		return m.CompanyName
	}
	return ""
}

func (m *Listing) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Listing) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *Listing) GetState() string {
	if m != nil {
		return m.State
	}
	return ""
}

func (m *Listing) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

func (m *Listing) GetZipCode() string {
	if m != nil {
		return m.ZipCode
	}
	return ""
}

func (m *Listing) GetLocation() *Geo {
	if m != nil {
		return m.Location
	}
	return nil
}

func (m *Listing) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *Listing) GetAdditionalPhoneNumbers() []string {
	if m != nil {
		return m.AdditionalPhoneNumbers
	}
	return nil
}

func (m *Listing) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *Listing) GetBusinessCategories() []string {
	if m != nil {
		return m.BusinessCategories
	}
	return nil
}

type GetListingRequest struct {
	ListingId  string `protobuf:"bytes,1,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	ExternalId string `protobuf:"bytes,2,opt,name=external_id,json=externalId" json:"external_id,omitempty"`
}

func (m *GetListingRequest) Reset()                    { *m = GetListingRequest{} }
func (m *GetListingRequest) String() string            { return proto.CompactTextString(m) }
func (*GetListingRequest) ProtoMessage()               {}
func (*GetListingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *GetListingRequest) GetListingId() string {
	if m != nil {
		return m.ListingId
	}
	return ""
}

func (m *GetListingRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

type DeleteListingRequest struct {
	ListingId  string `protobuf:"bytes,1,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	ExternalId string `protobuf:"bytes,2,opt,name=external_id,json=externalId" json:"external_id,omitempty"`
}

func (m *DeleteListingRequest) Reset()                    { *m = DeleteListingRequest{} }
func (m *DeleteListingRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteListingRequest) ProtoMessage()               {}
func (*DeleteListingRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DeleteListingRequest) GetListingId() string {
	if m != nil {
		return m.ListingId
	}
	return ""
}

func (m *DeleteListingRequest) GetExternalId() string {
	if m != nil {
		return m.ExternalId
	}
	return ""
}

func init() {
	proto.RegisterType((*Geo)(nil), "vendasta.listingsproto.Geo")
	proto.RegisterType((*Listing)(nil), "vendasta.listingsproto.Listing")
	proto.RegisterType((*GetListingRequest)(nil), "vendasta.listingsproto.GetListingRequest")
	proto.RegisterType((*DeleteListingRequest)(nil), "vendasta.listingsproto.DeleteListingRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ListingService service

type ListingServiceClient interface {
	Put(ctx context.Context, in *Listing, opts ...grpc.CallOption) (*Listing, error)
	Get(ctx context.Context, in *GetListingRequest, opts ...grpc.CallOption) (*Listing, error)
	Delete(ctx context.Context, in *DeleteListingRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type listingServiceClient struct {
	cc *grpc.ClientConn
}

func NewListingServiceClient(cc *grpc.ClientConn) ListingServiceClient {
	return &listingServiceClient{cc}
}

func (c *listingServiceClient) Put(ctx context.Context, in *Listing, opts ...grpc.CallOption) (*Listing, error) {
	out := new(Listing)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ListingService/Put", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listingServiceClient) Get(ctx context.Context, in *GetListingRequest, opts ...grpc.CallOption) (*Listing, error) {
	out := new(Listing)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ListingService/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listingServiceClient) Delete(ctx context.Context, in *DeleteListingRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vendasta.listingsproto.ListingService/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ListingService service

type ListingServiceServer interface {
	Put(context.Context, *Listing) (*Listing, error)
	Get(context.Context, *GetListingRequest) (*Listing, error)
	Delete(context.Context, *DeleteListingRequest) (*google_protobuf.Empty, error)
}

func RegisterListingServiceServer(s *grpc.Server, srv ListingServiceServer) {
	s.RegisterService(&_ListingService_serviceDesc, srv)
}

func _ListingService_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Listing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListingServiceServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ListingService/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListingServiceServer).Put(ctx, req.(*Listing))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListingService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListingServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ListingService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListingServiceServer).Get(ctx, req.(*GetListingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListingService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteListingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListingServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vendasta.listingsproto.ListingService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListingServiceServer).Delete(ctx, req.(*DeleteListingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ListingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vendasta.listingsproto.ListingService",
	HandlerType: (*ListingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _ListingService_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ListingService_Get_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ListingService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "listing.proto",
}

func init() { proto.RegisterFile("listing.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 473 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x93, 0x41, 0x6f, 0xd3, 0x40,
	0x10, 0x85, 0x95, 0xba, 0x4d, 0xe2, 0x49, 0x8b, 0xca, 0x50, 0x45, 0x4b, 0x0a, 0x6a, 0xc8, 0x29,
	0x48, 0xc8, 0x91, 0xca, 0x01, 0x6e, 0x1c, 0x0a, 0x8a, 0x2a, 0xa1, 0xa8, 0xb8, 0x12, 0x57, 0x6b,
	0xed, 0x1d, 0xcc, 0x4a, 0xf6, 0xae, 0xf1, 0xae, 0x0b, 0xe9, 0x5f, 0xe4, 0x2f, 0x71, 0x40, 0xde,
	0x75, 0x52, 0x09, 0x35, 0xea, 0x81, 0xde, 0xfc, 0xde, 0x37, 0x7e, 0x3b, 0x9a, 0xd1, 0xc0, 0x51,
	0x21, 0x8d, 0x95, 0x2a, 0x8f, 0xaa, 0x5a, 0x5b, 0x8d, 0xe3, 0x1b, 0x52, 0x82, 0x1b, 0xcb, 0xa3,
	0xce, 0x37, 0xce, 0x9f, 0x9c, 0xe6, 0x5a, 0xe7, 0x05, 0x2d, 0x9c, 0x4a, 0x9b, 0x6f, 0x0b, 0x2a,
	0x2b, 0xbb, 0xf6, 0x3f, 0xcd, 0x3e, 0x40, 0xb0, 0x24, 0x8d, 0x13, 0x18, 0x16, 0xdc, 0x4a, 0xdb,
	0x08, 0x62, 0xbd, 0x69, 0x6f, 0xde, 0x8b, 0xb7, 0x1a, 0x5f, 0x40, 0x58, 0x68, 0x95, 0x7b, 0xb8,
	0xe7, 0xe0, 0x9d, 0x31, 0xfb, 0x1d, 0xc0, 0xe0, 0xb3, 0x7f, 0x0f, 0x5f, 0x02, 0x74, 0x4f, 0x27,
	0x52, 0xb8, 0x9c, 0x30, 0x0e, 0x3b, 0xe7, 0x52, 0xe0, 0x19, 0x8c, 0xe8, 0x97, 0xa5, 0x5a, 0xf1,
	0xa2, 0xe5, 0x7b, 0x8e, 0xc3, 0xc6, 0xba, 0x14, 0x78, 0x0c, 0x41, 0x53, 0x17, 0x2c, 0x70, 0xa0,
	0xfd, 0xc4, 0x57, 0x70, 0x98, 0xe9, 0xb2, 0xe2, 0x6a, 0x9d, 0x28, 0x5e, 0x12, 0xdb, 0x77, 0x68,
	0xd4, 0x79, 0x2b, 0x5e, 0x12, 0x32, 0x18, 0x70, 0x21, 0x6a, 0x32, 0x86, 0x1d, 0x38, 0xba, 0x91,
	0x88, 0xb0, 0x9f, 0x49, 0xbb, 0x66, 0x7d, 0x67, 0xbb, 0x6f, 0x3c, 0x81, 0x03, 0x63, 0xb9, 0x25,
	0x36, 0x70, 0xa6, 0x17, 0x6d, 0x46, 0xa6, 0x1b, 0x65, 0xeb, 0x35, 0x1b, 0xfa, 0x8c, 0x4e, 0xe2,
	0x73, 0x18, 0xde, 0xca, 0x2a, 0xc9, 0xb4, 0x20, 0x16, 0x7a, 0x74, 0x2b, 0xab, 0x0b, 0x2d, 0x08,
	0xdf, 0xc1, 0xb0, 0xd0, 0x19, 0xb7, 0x52, 0x2b, 0x06, 0xd3, 0xde, 0x7c, 0x74, 0x7e, 0x1a, 0xdd,
	0xbf, 0x82, 0x68, 0x49, 0x3a, 0xde, 0x16, 0xb7, 0x3d, 0x54, 0xdf, 0xb5, 0x22, 0x36, 0xf2, 0x3d,
	0x38, 0x81, 0xef, 0x81, 0x71, 0x21, 0x64, 0x5b, 0xc1, 0x8b, 0xc4, 0x79, 0x89, 0x6a, 0xca, 0x94,
	0x6a, 0xc3, 0x0e, 0xa7, 0xc1, 0x3c, 0x8c, 0xc7, 0x77, 0xfc, 0xaa, 0xc5, 0x2b, 0x4f, 0xdb, 0xee,
	0x7f, 0x52, 0x6a, 0xa4, 0x25, 0x76, 0xe4, 0x5b, 0xec, 0x24, 0x2e, 0xe0, 0x59, 0xda, 0x18, 0xa9,
	0xc8, 0x98, 0x24, 0xe3, 0x96, 0x72, 0x5d, 0x4b, 0x32, 0xec, 0xd8, 0xc5, 0xe1, 0x06, 0x5d, 0x6c,
	0xc9, 0xec, 0x1a, 0x9e, 0x2e, 0xc9, 0x76, 0xfb, 0x8c, 0xe9, 0x47, 0x43, 0xc6, 0xfe, 0xef, 0x5a,
	0x67, 0x5f, 0xe1, 0xe4, 0x23, 0x15, 0x64, 0xe9, 0x71, 0x73, 0xcf, 0xff, 0xf4, 0xe0, 0x49, 0x17,
	0x79, 0x4d, 0xf5, 0x8d, 0xcc, 0x08, 0x97, 0x10, 0x5c, 0x35, 0x16, 0xcf, 0x76, 0x2d, 0xa2, 0x2b,
	0x9f, 0x3c, 0x54, 0x80, 0x5f, 0xda, 0xbb, 0xb0, 0xf8, 0x7a, 0xf7, 0x46, 0xff, 0x99, 0xd2, 0xc3,
	0x91, 0x2b, 0xe8, 0xfb, 0x31, 0xe0, 0x9b, 0x5d, 0xa5, 0xf7, 0x8d, 0x69, 0x32, 0x8e, 0xfc, 0x01,
	0x47, 0x9b, 0x03, 0x8e, 0x3e, 0xb5, 0x07, 0x9c, 0xf6, 0x9d, 0x7e, 0xfb, 0x37, 0x00, 0x00, 0xff,
	0xff, 0xcd, 0xaa, 0x58, 0x27, 0x07, 0x04, 0x00, 0x00,
}
