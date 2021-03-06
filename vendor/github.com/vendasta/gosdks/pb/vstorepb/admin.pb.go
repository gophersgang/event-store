// Code generated by protoc-gen-go.
// source: admin.proto
// DO NOT EDIT!

/*
Package vstorepb is a generated protocol buffer package.

It is generated from these files:
	admin.proto
	api.proto
	common.proto
	schema.proto
	secondary_index.proto

It has these top-level messages:
	CreateNamespaceRequest
	UpdateNamespaceRequest
	DeleteNamespaceRequest
	CreateKindRequest
	UpdateKindRequest
	GetKindRequest
	GetKindResponse
	DeleteKindRequest
	BackupConfig
	Property
	Entity
	Struct
	ListValue
	Value
	KeySet
	CreateRequest
	GetRequest
	GetResponse
	UpdateRequest
	EntityResult
	LookupFilter
	LookupRequest
	LookupResponse
	GeoPoint
	Schema
	NamespaceConfig
	SecondaryIndexPropertyConfig
	SecondaryIndex
	ElasticsearchRawConfig
	ElasticsearchConfig
	CloudSQLConfig
	ElasticsearchAnalysis
	ElasticsearchAnalyzer
	ElasticsearchFilter
	ElasticsearchCharFilter
	ElasticsearchTokenizer
*/
package vstorepb

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

type BackupConfig_BackupFrequency int32

const (
	BackupConfig_WEEKLY  BackupConfig_BackupFrequency = 0
	BackupConfig_DAILY   BackupConfig_BackupFrequency = 1
	BackupConfig_MONTHLY BackupConfig_BackupFrequency = 2
)

var BackupConfig_BackupFrequency_name = map[int32]string{
	0: "WEEKLY",
	1: "DAILY",
	2: "MONTHLY",
}
var BackupConfig_BackupFrequency_value = map[string]int32{
	"WEEKLY":  0,
	"DAILY":   1,
	"MONTHLY": 2,
}

func (x BackupConfig_BackupFrequency) String() string {
	return proto.EnumName(BackupConfig_BackupFrequency_name, int32(x))
}
func (BackupConfig_BackupFrequency) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{8, 0}
}

type Property_Type int32

const (
	Property_STRING    Property_Type = 0
	Property_INT64     Property_Type = 1
	Property_DOUBLE    Property_Type = 2
	Property_BOOL      Property_Type = 3
	Property_TIMESTAMP Property_Type = 4
	Property_GEOPOINT  Property_Type = 5
	Property_STRUCT    Property_Type = 6
)

var Property_Type_name = map[int32]string{
	0: "STRING",
	1: "INT64",
	2: "DOUBLE",
	3: "BOOL",
	4: "TIMESTAMP",
	5: "GEOPOINT",
	6: "STRUCT",
}
var Property_Type_value = map[string]int32{
	"STRING":    0,
	"INT64":     1,
	"DOUBLE":    2,
	"BOOL":      3,
	"TIMESTAMP": 4,
	"GEOPOINT":  5,
	"STRUCT":    6,
}

func (x Property_Type) String() string {
	return proto.EnumName(Property_Type_name, int32(x))
}
func (Property_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{9, 0} }

type CreateNamespaceRequest struct {
	// Unique namespace id unique to your project/microservice. Must be in lower snake case format.
	// Example(s): repcore, partner-central, central-identity-service, marketing-automation.
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// List of service account ids that are authorized to access the data in this namespace.
	AuthorizedServiceAccounts []string `protobuf:"bytes,2,rep,name=authorized_service_accounts,json=authorizedServiceAccounts" json:"authorized_service_accounts,omitempty"`
}

func (m *CreateNamespaceRequest) Reset()                    { *m = CreateNamespaceRequest{} }
func (m *CreateNamespaceRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateNamespaceRequest) ProtoMessage()               {}
func (*CreateNamespaceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CreateNamespaceRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *CreateNamespaceRequest) GetAuthorizedServiceAccounts() []string {
	if m != nil {
		return m.AuthorizedServiceAccounts
	}
	return nil
}

type UpdateNamespaceRequest struct {
	// Id of an existing namespace.
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// List of service account ids that are authorized to access the data in this namespace.
	// Replaces the list of authorized service accounts that are currently on the namespace.
	AuthorizedServiceAccounts []string `protobuf:"bytes,2,rep,name=authorized_service_accounts,json=authorizedServiceAccounts" json:"authorized_service_accounts,omitempty"`
}

func (m *UpdateNamespaceRequest) Reset()                    { *m = UpdateNamespaceRequest{} }
func (m *UpdateNamespaceRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateNamespaceRequest) ProtoMessage()               {}
func (*UpdateNamespaceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UpdateNamespaceRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *UpdateNamespaceRequest) GetAuthorizedServiceAccounts() []string {
	if m != nil {
		return m.AuthorizedServiceAccounts
	}
	return nil
}

type DeleteNamespaceRequest struct {
	// Id of an existing namespace.
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
}

func (m *DeleteNamespaceRequest) Reset()                    { *m = DeleteNamespaceRequest{} }
func (m *DeleteNamespaceRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteNamespaceRequest) ProtoMessage()               {}
func (*DeleteNamespaceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DeleteNamespaceRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

type CreateKindRequest struct {
	// Id of an existing namespace
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// Unique name of the kind that are creating. Must be in capital camel case format.
	// Example(s): AccountGroup, Partner, Review, Listing.
	Kind string `protobuf:"bytes,2,opt,name=kind" json:"kind,omitempty"`
	// List of fields that compose of the primary key. The order is important as it is used for building keysets,
	// as well as lookups can be done by the leading pieces of a keyset.
	PrimaryKey []string `protobuf:"bytes,3,rep,name=primary_key,json=primaryKey" json:"primary_key,omitempty"`
	// Schema for the kind. Indexing of any entities into this namespace/kind requires that a type has been set
	// for every field being indexed.  No inference is done and explicit types are required. Fields also are not
	// able to have their types changed or deleted, and only additive changes are allowed once a kind has been created.
	Properties []*Property `protobuf:"bytes,4,rep,name=properties" json:"properties,omitempty"`
	// Configured set of secondary indexes that you would like vStore to replicate to.
	SecondaryIndexes []*SecondaryIndex `protobuf:"bytes,5,rep,name=secondary_indexes,json=secondaryIndexes" json:"secondary_indexes,omitempty"`
	// Backup configuration
	BackupConfig *BackupConfig `protobuf:"bytes,6,opt,name=backup_config,json=backupConfig" json:"backup_config,omitempty"`
}

func (m *CreateKindRequest) Reset()                    { *m = CreateKindRequest{} }
func (m *CreateKindRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateKindRequest) ProtoMessage()               {}
func (*CreateKindRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CreateKindRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *CreateKindRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *CreateKindRequest) GetPrimaryKey() []string {
	if m != nil {
		return m.PrimaryKey
	}
	return nil
}

func (m *CreateKindRequest) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *CreateKindRequest) GetSecondaryIndexes() []*SecondaryIndex {
	if m != nil {
		return m.SecondaryIndexes
	}
	return nil
}

func (m *CreateKindRequest) GetBackupConfig() *BackupConfig {
	if m != nil {
		return m.BackupConfig
	}
	return nil
}

type UpdateKindRequest struct {
	// Id of an existing namespace
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// Id of an existing kind
	Kind string `protobuf:"bytes,2,opt,name=kind" json:"kind,omitempty"`
	// Schema for the kind with any new fields included in the request. Changes to any existing fields will cause
	// the request to fail.
	Properties []*Property `protobuf:"bytes,3,rep,name=properties" json:"properties,omitempty"`
	// Configured set of secondary indexes that you would like vStore to replicate to.
	SecondaryIndexes []*SecondaryIndex `protobuf:"bytes,5,rep,name=secondary_indexes,json=secondaryIndexes" json:"secondary_indexes,omitempty"`
}

func (m *UpdateKindRequest) Reset()                    { *m = UpdateKindRequest{} }
func (m *UpdateKindRequest) String() string            { return proto.CompactTextString(m) }
func (*UpdateKindRequest) ProtoMessage()               {}
func (*UpdateKindRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *UpdateKindRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *UpdateKindRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

func (m *UpdateKindRequest) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *UpdateKindRequest) GetSecondaryIndexes() []*SecondaryIndex {
	if m != nil {
		return m.SecondaryIndexes
	}
	return nil
}

type GetKindRequest struct {
	// Id of an existing namespace
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// Id of an existing kind
	Kind string `protobuf:"bytes,2,opt,name=kind" json:"kind,omitempty"`
}

func (m *GetKindRequest) Reset()                    { *m = GetKindRequest{} }
func (m *GetKindRequest) String() string            { return proto.CompactTextString(m) }
func (*GetKindRequest) ProtoMessage()               {}
func (*GetKindRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetKindRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *GetKindRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

type GetKindResponse struct {
	// List of fields that compose of the primary key. The order is important as it is used for building keysets,
	// as well as lookups can be done by the leading pieces of a keyset.
	PrimaryKey []string `protobuf:"bytes,1,rep,name=primary_key,json=primaryKey" json:"primary_key,omitempty"`
	// Schema for the kind. Indexing of any entities into this namespace/kind requires that a type has been set
	// for every field being indexed.  No inference is done and explicit types are required. Fields also are not
	// able to have their types changed or deleted, and only additive changes are allowed once a kind has been created.
	Properties []*Property `protobuf:"bytes,2,rep,name=properties" json:"properties,omitempty"`
	// Configured set of secondary indexes VStore is replicating to.
	SecondaryIndexes []*SecondaryIndex `protobuf:"bytes,3,rep,name=secondary_indexes,json=secondaryIndexes" json:"secondary_indexes,omitempty"`
	// Backup configuration
	BackupConfig *BackupConfig `protobuf:"bytes,4,opt,name=backup_config,json=backupConfig" json:"backup_config,omitempty"`
}

func (m *GetKindResponse) Reset()                    { *m = GetKindResponse{} }
func (m *GetKindResponse) String() string            { return proto.CompactTextString(m) }
func (*GetKindResponse) ProtoMessage()               {}
func (*GetKindResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *GetKindResponse) GetPrimaryKey() []string {
	if m != nil {
		return m.PrimaryKey
	}
	return nil
}

func (m *GetKindResponse) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *GetKindResponse) GetSecondaryIndexes() []*SecondaryIndex {
	if m != nil {
		return m.SecondaryIndexes
	}
	return nil
}

func (m *GetKindResponse) GetBackupConfig() *BackupConfig {
	if m != nil {
		return m.BackupConfig
	}
	return nil
}

type DeleteKindRequest struct {
	// Id of an existing namespace
	Namespace string `protobuf:"bytes,1,opt,name=namespace" json:"namespace,omitempty"`
	// Id of an existing kind
	Kind string `protobuf:"bytes,2,opt,name=kind" json:"kind,omitempty"`
}

func (m *DeleteKindRequest) Reset()                    { *m = DeleteKindRequest{} }
func (m *DeleteKindRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteKindRequest) ProtoMessage()               {}
func (*DeleteKindRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *DeleteKindRequest) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *DeleteKindRequest) GetKind() string {
	if m != nil {
		return m.Kind
	}
	return ""
}

type BackupConfig struct {
	Frequency BackupConfig_BackupFrequency `protobuf:"varint,1,opt,name=frequency,enum=vstorepb.BackupConfig_BackupFrequency" json:"frequency,omitempty"`
}

func (m *BackupConfig) Reset()                    { *m = BackupConfig{} }
func (m *BackupConfig) String() string            { return proto.CompactTextString(m) }
func (*BackupConfig) ProtoMessage()               {}
func (*BackupConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *BackupConfig) GetFrequency() BackupConfig_BackupFrequency {
	if m != nil {
		return m.Frequency
	}
	return BackupConfig_WEEKLY
}

type Property struct {
	// Unique identifier for this property. Must be in snake case format.
	// Example(s): account_group_id, listing_id, company_name
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Type for this property
	Type Property_Type `protobuf:"varint,2,opt,name=type,enum=vstorepb.Property_Type" json:"type,omitempty"`
	// Indicates if the field is repeated.
	Repeated bool `protobuf:"varint,3,opt,name=repeated" json:"repeated,omitempty"`
	// Indicates if the field is required. Only validates that the field has been supplied in create/update requests,
	// and not the actual value.
	// Example(s):
	// 1) A required string field would allow an empty string if the field was passed, but would fail if the field
	// was not present in the request.
	// 2) A required int property would allow 0 as a value, but would fail if the field was not supplied in the request.
	Required bool `protobuf:"varint,4,opt,name=required" json:"required,omitempty"`
	// Can only be specified if the Type supplied is a STRUCT.
	// Is the schema of the structured property.
	Properties            []*Property                              `protobuf:"bytes,5,rep,name=properties" json:"properties,omitempty"`
	SecondaryIndexConfigs map[string]*SecondaryIndexPropertyConfig `protobuf:"bytes,6,rep,name=secondary_index_configs,json=secondaryIndexConfigs" json:"secondary_index_configs,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Property) Reset()                    { *m = Property{} }
func (m *Property) String() string            { return proto.CompactTextString(m) }
func (*Property) ProtoMessage()               {}
func (*Property) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *Property) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Property) GetType() Property_Type {
	if m != nil {
		return m.Type
	}
	return Property_STRING
}

func (m *Property) GetRepeated() bool {
	if m != nil {
		return m.Repeated
	}
	return false
}

func (m *Property) GetRequired() bool {
	if m != nil {
		return m.Required
	}
	return false
}

func (m *Property) GetProperties() []*Property {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *Property) GetSecondaryIndexConfigs() map[string]*SecondaryIndexPropertyConfig {
	if m != nil {
		return m.SecondaryIndexConfigs
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateNamespaceRequest)(nil), "vstorepb.CreateNamespaceRequest")
	proto.RegisterType((*UpdateNamespaceRequest)(nil), "vstorepb.UpdateNamespaceRequest")
	proto.RegisterType((*DeleteNamespaceRequest)(nil), "vstorepb.DeleteNamespaceRequest")
	proto.RegisterType((*CreateKindRequest)(nil), "vstorepb.CreateKindRequest")
	proto.RegisterType((*UpdateKindRequest)(nil), "vstorepb.UpdateKindRequest")
	proto.RegisterType((*GetKindRequest)(nil), "vstorepb.GetKindRequest")
	proto.RegisterType((*GetKindResponse)(nil), "vstorepb.GetKindResponse")
	proto.RegisterType((*DeleteKindRequest)(nil), "vstorepb.DeleteKindRequest")
	proto.RegisterType((*BackupConfig)(nil), "vstorepb.BackupConfig")
	proto.RegisterType((*Property)(nil), "vstorepb.Property")
	proto.RegisterEnum("vstorepb.BackupConfig_BackupFrequency", BackupConfig_BackupFrequency_name, BackupConfig_BackupFrequency_value)
	proto.RegisterEnum("vstorepb.Property_Type", Property_Type_name, Property_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for VStoreAdmin service

type VStoreAdminClient interface {
	CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	UpdateNamespace(ctx context.Context, in *UpdateNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	CreateKind(ctx context.Context, in *CreateKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	UpdateKind(ctx context.Context, in *UpdateKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	GetKind(ctx context.Context, in *GetKindRequest, opts ...grpc.CallOption) (*GetKindResponse, error)
	DeleteKind(ctx context.Context, in *DeleteKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type vStoreAdminClient struct {
	cc *grpc.ClientConn
}

func NewVStoreAdminClient(cc *grpc.ClientConn) VStoreAdminClient {
	return &vStoreAdminClient{cc}
}

func (c *vStoreAdminClient) CreateNamespace(ctx context.Context, in *CreateNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/CreateNamespace", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) UpdateNamespace(ctx context.Context, in *UpdateNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/UpdateNamespace", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) DeleteNamespace(ctx context.Context, in *DeleteNamespaceRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/DeleteNamespace", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) CreateKind(ctx context.Context, in *CreateKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/CreateKind", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) UpdateKind(ctx context.Context, in *UpdateKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/UpdateKind", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) GetKind(ctx context.Context, in *GetKindRequest, opts ...grpc.CallOption) (*GetKindResponse, error) {
	out := new(GetKindResponse)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/GetKind", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vStoreAdminClient) DeleteKind(ctx context.Context, in *DeleteKindRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/vstorepb.VStoreAdmin/DeleteKind", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for VStoreAdmin service

type VStoreAdminServer interface {
	CreateNamespace(context.Context, *CreateNamespaceRequest) (*google_protobuf.Empty, error)
	UpdateNamespace(context.Context, *UpdateNamespaceRequest) (*google_protobuf.Empty, error)
	DeleteNamespace(context.Context, *DeleteNamespaceRequest) (*google_protobuf.Empty, error)
	CreateKind(context.Context, *CreateKindRequest) (*google_protobuf.Empty, error)
	UpdateKind(context.Context, *UpdateKindRequest) (*google_protobuf.Empty, error)
	GetKind(context.Context, *GetKindRequest) (*GetKindResponse, error)
	DeleteKind(context.Context, *DeleteKindRequest) (*google_protobuf.Empty, error)
}

func RegisterVStoreAdminServer(s *grpc.Server, srv VStoreAdminServer) {
	s.RegisterService(&_VStoreAdmin_serviceDesc, srv)
}

func _VStoreAdmin_CreateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).CreateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/CreateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).CreateNamespace(ctx, req.(*CreateNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_UpdateNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).UpdateNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/UpdateNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).UpdateNamespace(ctx, req.(*UpdateNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_DeleteNamespace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNamespaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).DeleteNamespace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/DeleteNamespace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).DeleteNamespace(ctx, req.(*DeleteNamespaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_CreateKind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).CreateKind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/CreateKind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).CreateKind(ctx, req.(*CreateKindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_UpdateKind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateKindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).UpdateKind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/UpdateKind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).UpdateKind(ctx, req.(*UpdateKindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_GetKind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).GetKind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/GetKind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).GetKind(ctx, req.(*GetKindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VStoreAdmin_DeleteKind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKindRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VStoreAdminServer).DeleteKind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vstorepb.VStoreAdmin/DeleteKind",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VStoreAdminServer).DeleteKind(ctx, req.(*DeleteKindRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _VStoreAdmin_serviceDesc = grpc.ServiceDesc{
	ServiceName: "vstorepb.VStoreAdmin",
	HandlerType: (*VStoreAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNamespace",
			Handler:    _VStoreAdmin_CreateNamespace_Handler,
		},
		{
			MethodName: "UpdateNamespace",
			Handler:    _VStoreAdmin_UpdateNamespace_Handler,
		},
		{
			MethodName: "DeleteNamespace",
			Handler:    _VStoreAdmin_DeleteNamespace_Handler,
		},
		{
			MethodName: "CreateKind",
			Handler:    _VStoreAdmin_CreateKind_Handler,
		},
		{
			MethodName: "UpdateKind",
			Handler:    _VStoreAdmin_UpdateKind_Handler,
		},
		{
			MethodName: "GetKind",
			Handler:    _VStoreAdmin_GetKind_Handler,
		},
		{
			MethodName: "DeleteKind",
			Handler:    _VStoreAdmin_DeleteKind_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}

func init() { proto.RegisterFile("admin.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 790 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xbc, 0x55, 0xed, 0x6a, 0xdb, 0x48,
	0x14, 0x8d, 0x2c, 0xd9, 0xb1, 0xaf, 0xf3, 0x21, 0x0f, 0xc4, 0x51, 0x9c, 0x85, 0x35, 0xfa, 0xb1,
	0x18, 0x96, 0x75, 0xc0, 0xdb, 0x86, 0xd2, 0x96, 0x52, 0x7f, 0xa8, 0xa9, 0x89, 0x3f, 0x82, 0xac,
	0xb4, 0x04, 0x0a, 0x46, 0x96, 0x26, 0xa9, 0x48, 0x2c, 0x29, 0x92, 0x6c, 0xaa, 0xbe, 0x42, 0xde,
	0xa1, 0xcf, 0xd2, 0xe7, 0xe9, 0x8f, 0x3e, 0x43, 0x19, 0x8d, 0x1c, 0xc9, 0x8a, 0xe3, 0x26, 0x35,
	0xe4, 0x9f, 0xe6, 0xde, 0x3b, 0x67, 0xce, 0xdc, 0x33, 0xba, 0x07, 0xf2, 0xaa, 0x3e, 0x36, 0xcc,
	0xaa, 0xed, 0x58, 0x9e, 0x85, 0xb2, 0x53, 0xd7, 0xb3, 0x1c, 0x6c, 0x8f, 0x4a, 0xfb, 0x17, 0x96,
	0x75, 0x71, 0x85, 0x0f, 0x82, 0xf8, 0x68, 0x72, 0x7e, 0x80, 0xc7, 0xb6, 0xe7, 0xd3, 0xb2, 0xd2,
	0x8e, 0x8b, 0x35, 0xcb, 0xd4, 0x55, 0xc7, 0x1f, 0x1a, 0xa6, 0x8e, 0xbf, 0xd0, 0xb0, 0x38, 0x85,
	0x62, 0xd3, 0xc1, 0xaa, 0x87, 0x7b, 0xea, 0x18, 0xbb, 0xb6, 0xaa, 0x61, 0x19, 0x5f, 0x4f, 0xb0,
	0xeb, 0xa1, 0xbf, 0x20, 0x67, 0xce, 0x62, 0x02, 0x53, 0x66, 0x2a, 0x39, 0x39, 0x0a, 0xa0, 0x37,
	0xb0, 0xaf, 0x4e, 0xbc, 0xcf, 0x96, 0x63, 0x7c, 0xc5, 0xfa, 0xd0, 0xc5, 0xce, 0xd4, 0xd0, 0xf0,
	0x50, 0xd5, 0x34, 0x6b, 0x62, 0x7a, 0xae, 0x90, 0x2a, 0xb3, 0x95, 0x9c, 0xbc, 0x17, 0x95, 0x0c,
	0x68, 0x45, 0x3d, 0x2c, 0x20, 0xe7, 0x9e, 0xda, 0xfa, 0xd3, 0x9f, 0x7b, 0x08, 0xc5, 0x16, 0xbe,
	0xc2, 0x8f, 0x3d, 0x57, 0xfc, 0x96, 0x82, 0x02, 0x6d, 0xd4, 0xb1, 0x61, 0xea, 0x0f, 0xe3, 0x8a,
	0x80, 0xbb, 0x34, 0x4c, 0x5d, 0x48, 0x05, 0x89, 0xe0, 0x1b, 0xfd, 0x0d, 0x79, 0xdb, 0x31, 0xc6,
	0x44, 0x86, 0x4b, 0xec, 0x0b, 0x6c, 0xc0, 0x17, 0xc2, 0xd0, 0x31, 0xf6, 0x51, 0x0d, 0xc0, 0x76,
	0x2c, 0x1b, 0x3b, 0x9e, 0x81, 0x5d, 0x81, 0x2b, 0xb3, 0x95, 0x7c, 0x0d, 0x55, 0x67, 0x1a, 0x57,
	0x4f, 0x68, 0xce, 0x97, 0x63, 0x55, 0x48, 0x82, 0x42, 0x42, 0x5d, 0xec, 0x0a, 0xe9, 0x60, 0xab,
	0x10, 0x6d, 0x1d, 0xcc, 0x4a, 0xda, 0xa4, 0x42, 0xe6, 0xdd, 0xb9, 0x35, 0x76, 0xd1, 0x2b, 0xd8,
	0x1c, 0xa9, 0xda, 0xe5, 0xc4, 0x1e, 0x6a, 0x96, 0x79, 0x6e, 0x5c, 0x08, 0x99, 0x32, 0x53, 0xc9,
	0xd7, 0x8a, 0x11, 0x44, 0x23, 0x48, 0x37, 0x83, 0xac, 0xbc, 0x31, 0x8a, 0xad, 0xc4, 0xef, 0x0c,
	0x14, 0xa8, 0xa2, 0xab, 0x35, 0x68, 0xfe, 0xfe, 0xec, 0x13, 0xde, 0x5f, 0x6c, 0xc0, 0xd6, 0x11,
	0xf6, 0x56, 0xa2, 0x2f, 0xfe, 0x60, 0x60, 0xfb, 0x16, 0xc4, 0xb5, 0x2d, 0xd3, 0xc5, 0x49, 0xcd,
	0x99, 0xdf, 0x68, 0x9e, 0xfa, 0xf3, 0x3b, 0xb3, 0xab, 0x6b, 0xce, 0x3d, 0x42, 0x73, 0x09, 0x0a,
	0xf4, 0x67, 0x5a, 0xad, 0x67, 0x37, 0x0c, 0x6c, 0xc4, 0x4f, 0x41, 0x2d, 0xc8, 0x9d, 0x3b, 0x04,
	0xce, 0xd4, 0xfc, 0x00, 0x62, 0xab, 0xf6, 0xcf, 0x62, 0x42, 0xe1, 0xe2, 0xdd, 0xac, 0x5a, 0x8e,
	0x36, 0x8a, 0xcf, 0x61, 0x3b, 0x91, 0x45, 0x00, 0x99, 0x8f, 0x92, 0x74, 0xdc, 0x39, 0xe3, 0xd7,
	0x50, 0x0e, 0xd2, 0xad, 0x7a, 0xbb, 0x73, 0xc6, 0x33, 0x28, 0x0f, 0xeb, 0xdd, 0x7e, 0x4f, 0x79,
	0xdf, 0x39, 0xe3, 0x53, 0xe2, 0x4f, 0x16, 0xb2, 0xb3, 0x8e, 0x13, 0xba, 0x84, 0x7b, 0x78, 0x8f,
	0xe0, 0x1b, 0xfd, 0x0b, 0x9c, 0xe7, 0xdb, 0x38, 0xb8, 0xc2, 0x56, 0x6d, 0xf7, 0xae, 0x4e, 0x55,
	0xc5, 0xb7, 0xb1, 0x1c, 0x14, 0xa1, 0x12, 0x64, 0x1d, 0x6c, 0x93, 0xb9, 0xa1, 0x0b, 0x6c, 0x99,
	0xa9, 0x64, 0xe5, 0xdb, 0x35, 0xcd, 0x5d, 0x4f, 0x0c, 0x07, 0xeb, 0x41, 0xdb, 0x83, 0x1c, 0x5d,
	0x27, 0x9e, 0x44, 0xfa, 0x41, 0x4f, 0x02, 0xc3, 0x6e, 0xe2, 0x49, 0x84, 0xa2, 0xba, 0x42, 0x26,
	0x00, 0xf8, 0x6f, 0x01, 0xd7, 0xf9, 0x17, 0x42, 0xbb, 0xea, 0x4a, 0xa6, 0xe7, 0xf8, 0xf2, 0x8e,
	0xbb, 0x28, 0x57, 0xb2, 0xa1, 0x74, 0xff, 0x26, 0xc4, 0x03, 0x4b, 0x1f, 0x39, 0x69, 0x18, 0xf9,
	0x44, 0xaf, 0x21, 0x3d, 0x55, 0xaf, 0x26, 0xb4, 0x61, 0xf9, 0xb8, 0x92, 0xf3, 0x30, 0x33, 0x4a,
	0xe1, 0x53, 0xa3, 0x9b, 0x5e, 0xa6, 0x5e, 0x30, 0xe2, 0x27, 0xe0, 0x48, 0x4b, 0x89, 0x7c, 0x03,
	0x45, 0x6e, 0xf7, 0x8e, 0xa8, 0x7c, 0xed, 0x9e, 0x72, 0xf8, 0x8c, 0x67, 0x48, 0xb8, 0xd5, 0x3f,
	0x6d, 0x74, 0x24, 0x3e, 0x85, 0xb2, 0xc0, 0x35, 0xfa, 0xfd, 0x0e, 0xcf, 0xa2, 0x4d, 0xc8, 0x29,
	0xed, 0xae, 0x34, 0x50, 0xea, 0xdd, 0x13, 0x9e, 0x43, 0x1b, 0x90, 0x3d, 0x92, 0xfa, 0x27, 0xfd,
	0x76, 0x4f, 0xe1, 0xd3, 0x21, 0xd2, 0x69, 0x53, 0xe1, 0x33, 0xb5, 0x1b, 0x0e, 0xf2, 0x1f, 0x06,
	0x84, 0x52, 0x9d, 0xd8, 0x2a, 0xea, 0xc2, 0x76, 0xc2, 0x12, 0x51, 0x39, 0xe2, 0xbc, 0xd8, 0x2d,
	0x4b, 0xc5, 0x2a, 0x35, 0xdf, 0xea, 0xcc, 0x7c, 0xab, 0x12, 0x31, 0x5f, 0x71, 0x8d, 0xc0, 0x25,
	0x9c, 0x2e, 0x0e, 0xb7, 0xd8, 0x04, 0x97, 0xc3, 0x25, 0x0c, 0x2c, 0x0e, 0xb7, 0xd8, 0xdb, 0x96,
	0xc0, 0x35, 0x01, 0x22, 0x5b, 0x43, 0xfb, 0xc9, 0x7b, 0xc6, 0x7e, 0xec, 0xe5, 0x20, 0xd1, 0xe8,
	0x8f, 0x83, 0xdc, 0x31, 0x84, 0x25, 0x20, 0x6f, 0x61, 0x3d, 0x1c, 0x9c, 0x28, 0x36, 0xc0, 0xe6,
	0x07, 0x72, 0x69, 0x6f, 0x41, 0x86, 0x4e, 0x59, 0x4a, 0x23, 0x1a, 0x47, 0x71, 0x1a, 0x77, 0x86,
	0xd4, 0xfd, 0x34, 0x46, 0x99, 0x20, 0xf2, 0xff, 0xaf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7b, 0x7d,
	0xeb, 0x9e, 0x64, 0x09, 0x00, 0x00,
}
