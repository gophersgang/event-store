// Code generated by protoc-gen-go.
// source: review.proto
// DO NOT EDIT!

package datalakeproto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Review struct {
	ReviewId      string                      `protobuf:"bytes,1,opt,name=review_id,json=reviewId" json:"review_id,omitempty"`
	ListingId     string                      `protobuf:"bytes,2,opt,name=listing_id,json=listingId" json:"listing_id,omitempty"`
	Url           string                      `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
	StarRating    float64                     `protobuf:"fixed64,4,opt,name=star_rating,json=starRating" json:"star_rating,omitempty"`
	ReviewerName  string                      `protobuf:"bytes,5,opt,name=reviewer_name,json=reviewerName" json:"reviewer_name,omitempty"`
	ReviewerEmail string                      `protobuf:"bytes,6,opt,name=reviewer_email,json=reviewerEmail" json:"reviewer_email,omitempty"`
	ReviewerUrl   string                      `protobuf:"bytes,7,opt,name=reviewer_url,json=reviewerUrl" json:"reviewer_url,omitempty"`
	Content       string                      `protobuf:"bytes,8,opt,name=content" json:"content,omitempty"`
	PublishedDate *google_protobuf1.Timestamp `protobuf:"bytes,9,opt,name=published_date,json=publishedDate" json:"published_date,omitempty"`
	Title         string                      `protobuf:"bytes,10,opt,name=title" json:"title,omitempty"`
	SourceId      int64                       `protobuf:"varint,11,opt,name=source_id,json=sourceId" json:"source_id,omitempty"`
	DeletedOn     *google_protobuf1.Timestamp `protobuf:"bytes,12,opt,name=deleted_on,json=deletedOn" json:"deleted_on,omitempty"`
}

func (m *Review) Reset()                    { *m = Review{} }
func (m *Review) String() string            { return proto.CompactTextString(m) }
func (*Review) ProtoMessage()               {}
func (*Review) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

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

func (m *Review) GetStarRating() float64 {
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

func (m *Review) GetSourceId() int64 {
	if m != nil {
		return m.SourceId
	}
	return 0
}

func (m *Review) GetDeletedOn() *google_protobuf1.Timestamp {
	if m != nil {
		return m.DeletedOn
	}
	return nil
}

func init() {
	proto.RegisterType((*Review)(nil), "datalakeproto.Review")
}

func init() { proto.RegisterFile("review.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x90, 0xcf, 0x4b, 0xc3, 0x30,
	0x14, 0xc7, 0xa9, 0x75, 0x3f, 0xfa, 0xba, 0x0d, 0x09, 0x1e, 0xc2, 0x44, 0x56, 0x15, 0xa1, 0xa7,
	0x0e, 0xf4, 0xe4, 0x51, 0xd0, 0xc3, 0x2e, 0x0a, 0x45, 0xcf, 0x25, 0x5b, 0x9e, 0x33, 0x98, 0x36,
	0x23, 0x7d, 0xd5, 0xbf, 0xcb, 0xff, 0x50, 0x92, 0xac, 0xbd, 0x7a, 0xeb, 0xfb, 0x7c, 0x7f, 0x34,
	0x7c, 0x61, 0x66, 0xf1, 0x5b, 0xe1, 0x4f, 0x71, 0xb0, 0x86, 0x0c, 0x9b, 0x4b, 0x41, 0x42, 0x8b,
	0x2f, 0xf4, 0xe7, 0x72, 0xb5, 0x37, 0x66, 0xaf, 0x71, 0xed, 0xaf, 0x6d, 0xf7, 0xb1, 0x26, 0x55,
	0x63, 0x4b, 0xa2, 0x3e, 0x04, 0xff, 0xf5, 0x6f, 0x0c, 0xe3, 0xd2, 0x17, 0xb0, 0x0b, 0x48, 0x42,
	0x55, 0xa5, 0x24, 0x8f, 0xb2, 0x28, 0x4f, 0xca, 0x69, 0x00, 0x1b, 0xc9, 0x2e, 0x01, 0xb4, 0x6a,
	0x49, 0x35, 0x7b, 0xa7, 0x9e, 0x78, 0x35, 0x39, 0x92, 0x8d, 0x64, 0x67, 0x10, 0x77, 0x56, 0xf3,
	0xd8, 0x73, 0xf7, 0xc9, 0x56, 0x90, 0xb6, 0x24, 0x6c, 0x65, 0x85, 0xb3, 0xf0, 0xd3, 0x2c, 0xca,
	0xa3, 0x12, 0x1c, 0x2a, 0x3d, 0x61, 0x37, 0x30, 0x0f, 0xed, 0x68, 0xab, 0x46, 0xd4, 0xc8, 0x47,
	0x3e, 0x3c, 0xeb, 0xe1, 0x8b, 0xa8, 0x91, 0xdd, 0xc2, 0x62, 0x30, 0x61, 0x2d, 0x94, 0xe6, 0x63,
	0xef, 0x1a, 0xa2, 0xcf, 0x0e, 0xb2, 0x2b, 0x18, 0x62, 0x95, 0x7b, 0xc7, 0xc4, 0x9b, 0xd2, 0x9e,
	0xbd, 0x5b, 0xcd, 0x38, 0x4c, 0x76, 0xa6, 0x21, 0x6c, 0x88, 0x4f, 0xbd, 0xda, 0x9f, 0xec, 0x11,
	0x16, 0x87, 0x6e, 0xab, 0x55, 0xfb, 0x89, 0xb2, 0x92, 0x82, 0x90, 0x27, 0x59, 0x94, 0xa7, 0x77,
	0xcb, 0x22, 0x8c, 0x57, 0xf4, 0xe3, 0x15, 0x6f, 0xfd, 0x78, 0xe5, 0x7c, 0x48, 0x3c, 0x09, 0x42,
	0x76, 0x0e, 0x23, 0x52, 0xa4, 0x91, 0x83, 0xaf, 0x0e, 0x87, 0x1b, 0xb4, 0x35, 0x9d, 0xdd, 0xa1,
	0x9b, 0x2c, 0xcd, 0xa2, 0x3c, 0x2e, 0xa7, 0x01, 0x6c, 0x24, 0x7b, 0x00, 0x90, 0xa8, 0x91, 0x50,
	0x56, 0xa6, 0xe1, 0xb3, 0x7f, 0xff, 0x98, 0x1c, 0xdd, 0xaf, 0xcd, 0x76, 0xec, 0xe5, 0xfb, 0xbf,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x58, 0xdd, 0xdc, 0xfa, 0x01, 0x00, 0x00,
}
