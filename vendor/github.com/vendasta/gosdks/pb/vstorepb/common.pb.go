// Code generated by protoc-gen-go.
// source: common.proto
// DO NOT EDIT!

package vstorepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type GeoPoint struct {
	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude" json:"longitude,omitempty"`
}

func (m *GeoPoint) Reset()                    { *m = GeoPoint{} }
func (m *GeoPoint) String() string            { return proto.CompactTextString(m) }
func (*GeoPoint) ProtoMessage()               {}
func (*GeoPoint) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *GeoPoint) GetLatitude() float64 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *GeoPoint) GetLongitude() float64 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func init() {
	proto.RegisterType((*GeoPoint)(nil), "vstorepb.GeoPoint")
}

func init() { proto.RegisterFile("common.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 101 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x28, 0x2b, 0x2e, 0xc9, 0x2f, 0x4a,
	0x2d, 0x48, 0x52, 0x72, 0xe1, 0xe2, 0x70, 0x4f, 0xcd, 0x0f, 0xc8, 0xcf, 0xcc, 0x2b, 0x11, 0x92,
	0xe2, 0xe2, 0xc8, 0x49, 0x2c, 0xc9, 0x2c, 0x29, 0x4d, 0x49, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x60,
	0x0c, 0x82, 0xf3, 0x85, 0x64, 0xb8, 0x38, 0x73, 0xf2, 0xf3, 0xd2, 0x21, 0x92, 0x4c, 0x60, 0x49,
	0x84, 0x40, 0x12, 0x1b, 0xd8, 0x58, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xae, 0x82,
	0xd9, 0x66, 0x00, 0x00, 0x00,
}
