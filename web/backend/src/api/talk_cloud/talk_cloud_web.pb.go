// Code generated by protoc-gen-go. DO NOT EDIT.
// source: talk_cloud_web.proto

package talk_cloud

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type UpdateGroupReq struct {
	DeviceIds            []int64   `protobuf:"varint,2,rep,packed,name=deviceIds,proto3" json:"deviceIds,omitempty"`
	DeviceInfos          []*Member `protobuf:"bytes,3,rep,name=deviceInfos,proto3" json:"deviceInfos,omitempty"`
	GroupInfo            *Group    `protobuf:"bytes,4,opt,name=groupInfo,proto3" json:"groupInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *UpdateGroupReq) Reset()         { *m = UpdateGroupReq{} }
func (m *UpdateGroupReq) String() string { return proto.CompactTextString(m) }
func (*UpdateGroupReq) ProtoMessage()    {}
func (*UpdateGroupReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{0}
}

func (m *UpdateGroupReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateGroupReq.Unmarshal(m, b)
}
func (m *UpdateGroupReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateGroupReq.Marshal(b, m, deterministic)
}
func (m *UpdateGroupReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateGroupReq.Merge(m, src)
}
func (m *UpdateGroupReq) XXX_Size() int {
	return xxx_messageInfo_UpdateGroupReq.Size(m)
}
func (m *UpdateGroupReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateGroupReq.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateGroupReq proto.InternalMessageInfo

func (m *UpdateGroupReq) GetDeviceIds() []int64 {
	if m != nil {
		return m.DeviceIds
	}
	return nil
}

func (m *UpdateGroupReq) GetDeviceInfos() []*Member {
	if m != nil {
		return m.DeviceInfos
	}
	return nil
}

func (m *UpdateGroupReq) GetGroupInfo() *Group {
	if m != nil {
		return m.GroupInfo
	}
	return nil
}

type UpdateGroupResp struct {
	ResultMsg            *Result  `protobuf:"bytes,1,opt,name=resultMsg,proto3" json:"resultMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateGroupResp) Reset()         { *m = UpdateGroupResp{} }
func (m *UpdateGroupResp) String() string { return proto.CompactTextString(m) }
func (*UpdateGroupResp) ProtoMessage()    {}
func (*UpdateGroupResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{1}
}

func (m *UpdateGroupResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateGroupResp.Unmarshal(m, b)
}
func (m *UpdateGroupResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateGroupResp.Marshal(b, m, deterministic)
}
func (m *UpdateGroupResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateGroupResp.Merge(m, src)
}
func (m *UpdateGroupResp) XXX_Size() int {
	return xxx_messageInfo_UpdateGroupResp.Size(m)
}
func (m *UpdateGroupResp) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateGroupResp.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateGroupResp proto.InternalMessageInfo

func (m *UpdateGroupResp) GetResultMsg() *Result {
	if m != nil {
		return m.ResultMsg
	}
	return nil
}

type DeleteGroupResp struct {
	ResultMsg            *Result  `protobuf:"bytes,1,opt,name=resultMsg,proto3" json:"resultMsg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteGroupResp) Reset()         { *m = DeleteGroupResp{} }
func (m *DeleteGroupResp) String() string { return proto.CompactTextString(m) }
func (*DeleteGroupResp) ProtoMessage()    {}
func (*DeleteGroupResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{2}
}

func (m *DeleteGroupResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteGroupResp.Unmarshal(m, b)
}
func (m *DeleteGroupResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteGroupResp.Marshal(b, m, deterministic)
}
func (m *DeleteGroupResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteGroupResp.Merge(m, src)
}
func (m *DeleteGroupResp) XXX_Size() int {
	return xxx_messageInfo_DeleteGroupResp.Size(m)
}
func (m *DeleteGroupResp) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteGroupResp.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteGroupResp proto.InternalMessageInfo

func (m *DeleteGroupResp) GetResultMsg() *Result {
	if m != nil {
		return m.ResultMsg
	}
	return nil
}

type ImportDeviceReq struct {
	AccountId            int32    `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	DeviceImei           []string `protobuf:"bytes,2,rep,name=deviceImei,proto3" json:"deviceImei,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportDeviceReq) Reset()         { *m = ImportDeviceReq{} }
func (m *ImportDeviceReq) String() string { return proto.CompactTextString(m) }
func (*ImportDeviceReq) ProtoMessage()    {}
func (*ImportDeviceReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{3}
}

func (m *ImportDeviceReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportDeviceReq.Unmarshal(m, b)
}
func (m *ImportDeviceReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportDeviceReq.Marshal(b, m, deterministic)
}
func (m *ImportDeviceReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportDeviceReq.Merge(m, src)
}
func (m *ImportDeviceReq) XXX_Size() int {
	return xxx_messageInfo_ImportDeviceReq.Size(m)
}
func (m *ImportDeviceReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportDeviceReq.DiscardUnknown(m)
}

var xxx_messageInfo_ImportDeviceReq proto.InternalMessageInfo

func (m *ImportDeviceReq) GetAccountId() int32 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *ImportDeviceReq) GetDeviceImei() []string {
	if m != nil {
		return m.DeviceImei
	}
	return nil
}

type ImportDeviceResp struct {
	Result               *Result  `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportDeviceResp) Reset()         { *m = ImportDeviceResp{} }
func (m *ImportDeviceResp) String() string { return proto.CompactTextString(m) }
func (*ImportDeviceResp) ProtoMessage()    {}
func (*ImportDeviceResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{4}
}

func (m *ImportDeviceResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportDeviceResp.Unmarshal(m, b)
}
func (m *ImportDeviceResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportDeviceResp.Marshal(b, m, deterministic)
}
func (m *ImportDeviceResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportDeviceResp.Merge(m, src)
}
func (m *ImportDeviceResp) XXX_Size() int {
	return xxx_messageInfo_ImportDeviceResp.Size(m)
}
func (m *ImportDeviceResp) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportDeviceResp.DiscardUnknown(m)
}

var xxx_messageInfo_ImportDeviceResp proto.InternalMessageInfo

func (m *ImportDeviceResp) GetResult() *Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func init() {
	proto.RegisterType((*UpdateGroupReq)(nil), "talk_cloud.UpdateGroupReq")
	proto.RegisterType((*UpdateGroupResp)(nil), "talk_cloud.UpdateGroupResp")
	proto.RegisterType((*DeleteGroupResp)(nil), "talk_cloud.DeleteGroupResp")
	proto.RegisterType((*ImportDeviceReq)(nil), "talk_cloud.ImportDeviceReq")
	proto.RegisterType((*ImportDeviceResp)(nil), "talk_cloud.ImportDeviceResp")
}

func init() { proto.RegisterFile("talk_cloud_web.proto", fileDescriptor_bee75d1763415d2b) }

var fileDescriptor_bee75d1763415d2b = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x1b, 0xa3, 0x85, 0x4c, 0xc0, 0xea, 0x20, 0x12, 0x62, 0x91, 0x92, 0x53, 0xf1, 0x50,
	0xa5, 0x7a, 0x56, 0xd0, 0x82, 0xf6, 0x50, 0xc4, 0x15, 0xf1, 0x58, 0x9a, 0x64, 0x2c, 0xc5, 0xa4,
	0xbb, 0x66, 0x37, 0x8a, 0x0f, 0xe2, 0x4b, 0xfa, 0x14, 0x92, 0x34, 0xb0, 0xbb, 0x95, 0x7a, 0xf1,
	0xb8, 0x33, 0xff, 0xff, 0xcd, 0xce, 0xbf, 0x0b, 0x07, 0x6a, 0x96, 0xbd, 0x4e, 0x93, 0x8c, 0x97,
	0xe9, 0xf4, 0x83, 0xe2, 0x81, 0x28, 0xb8, 0xe2, 0x08, 0xba, 0x1a, 0x1e, 0x1a, 0x8a, 0x9c, 0xa7,
	0x94, 0xad, 0x34, 0xd1, 0x97, 0x03, 0xbb, 0x4f, 0x22, 0x9d, 0x29, 0xba, 0x2d, 0x78, 0x29, 0x18,
	0xbd, 0x61, 0x17, 0xbc, 0x94, 0xde, 0x17, 0x09, 0x8d, 0x53, 0x19, 0x6c, 0xf5, 0xdc, 0xbe, 0xcb,
	0x74, 0x01, 0x2f, 0xc0, 0x6f, 0x0e, 0xcb, 0x17, 0x2e, 0x03, 0xb7, 0xe7, 0xf6, 0xfd, 0x21, 0x0e,
	0x34, 0x7e, 0x30, 0xa1, 0x3c, 0xa6, 0x82, 0x99, 0x32, 0x3c, 0x05, 0x6f, 0x5e, 0xf1, 0xab, 0x53,
	0xb0, 0xdd, 0x73, 0xfa, 0xfe, 0x70, 0xdf, 0xf4, 0xac, 0x86, 0x6b, 0x4d, 0x74, 0x03, 0x1d, 0xeb,
	0x5a, 0x52, 0xe0, 0x19, 0x78, 0x05, 0xc9, 0x32, 0x53, 0x13, 0x39, 0x0f, 0x9c, 0x9a, 0x61, 0xcd,
	0x65, 0x75, 0x93, 0x69, 0x51, 0x05, 0x19, 0x51, 0x46, 0xff, 0x83, 0xdc, 0x43, 0x67, 0x9c, 0x0b,
	0x5e, 0xa8, 0x51, 0xbd, 0x4f, 0x93, 0xd0, 0x2c, 0x49, 0x78, 0xb9, 0x54, 0xe3, 0xb4, 0x86, 0xec,
	0x30, 0x5d, 0xc0, 0x63, 0x80, 0x66, 0xf5, 0x9c, 0x16, 0x75, 0x80, 0x1e, 0x33, 0x2a, 0xd1, 0x25,
	0xec, 0xd9, 0x40, 0x29, 0xf0, 0x04, 0xda, 0xab, 0x89, 0x7f, 0xdc, 0xa9, 0x51, 0x0c, 0xbf, 0x1d,
	0x80, 0x67, 0x8a, 0x1f, 0xa9, 0xa8, 0xec, 0x78, 0x07, 0xbe, 0x91, 0x14, 0x86, 0xa6, 0xd3, 0x7e,
	0xd9, 0xf0, 0x68, 0x63, 0x4f, 0x8a, 0xa8, 0x85, 0x57, 0xe0, 0x1b, 0x71, 0xe1, 0xef, 0x07, 0xb2,
	0x01, 0x6b, 0xd1, 0x46, 0x2d, 0x7c, 0x00, 0x34, 0x37, 0xbb, 0xfe, 0x64, 0x9c, 0x2b, 0xb4, 0x4c,
	0x6b, 0x51, 0x86, 0xdd, 0xcd, 0xcd, 0x0a, 0x19, 0xb7, 0xeb, 0x6f, 0x7a, 0xfe, 0x13, 0x00, 0x00,
	0xff, 0xff, 0xe0, 0x03, 0x97, 0xd9, 0xe2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WebServiceClient is the client API for WebService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServiceClient interface {
	UpdateGroup(ctx context.Context, in *UpdateGroupReq, opts ...grpc.CallOption) (*UpdateGroupResp, error)
	DeleteGroup(ctx context.Context, in *Group, opts ...grpc.CallOption) (*DeleteGroupResp, error)
	ImportDeviceByRoot(ctx context.Context, in *ImportDeviceReq, opts ...grpc.CallOption) (*ImportDeviceResp, error)
}

type webServiceClient struct {
	cc *grpc.ClientConn
}

func NewWebServiceClient(cc *grpc.ClientConn) WebServiceClient {
	return &webServiceClient{cc}
}

func (c *webServiceClient) UpdateGroup(ctx context.Context, in *UpdateGroupReq, opts ...grpc.CallOption) (*UpdateGroupResp, error) {
	out := new(UpdateGroupResp)
	err := c.cc.Invoke(ctx, "/talk_cloud.WebService/UpdateGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webServiceClient) DeleteGroup(ctx context.Context, in *Group, opts ...grpc.CallOption) (*DeleteGroupResp, error) {
	out := new(DeleteGroupResp)
	err := c.cc.Invoke(ctx, "/talk_cloud.WebService/DeleteGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webServiceClient) ImportDeviceByRoot(ctx context.Context, in *ImportDeviceReq, opts ...grpc.CallOption) (*ImportDeviceResp, error) {
	out := new(ImportDeviceResp)
	err := c.cc.Invoke(ctx, "/talk_cloud.WebService/ImportDeviceByRoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WebServiceServer is the server API for WebService service.
type WebServiceServer interface {
	UpdateGroup(context.Context, *UpdateGroupReq) (*UpdateGroupResp, error)
	DeleteGroup(context.Context, *Group) (*DeleteGroupResp, error)
	ImportDeviceByRoot(context.Context, *ImportDeviceReq) (*ImportDeviceResp, error)
}

// UnimplementedWebServiceServer can be embedded to have forward compatible implementations.
type UnimplementedWebServiceServer struct {
}

func (*UnimplementedWebServiceServer) UpdateGroup(ctx context.Context, req *UpdateGroupReq) (*UpdateGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGroup not implemented")
}
func (*UnimplementedWebServiceServer) DeleteGroup(ctx context.Context, req *Group) (*DeleteGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteGroup not implemented")
}
func (*UnimplementedWebServiceServer) ImportDeviceByRoot(ctx context.Context, req *ImportDeviceReq) (*ImportDeviceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ImportDeviceByRoot not implemented")
}

func RegisterWebServiceServer(s *grpc.Server, srv WebServiceServer) {
	s.RegisterService(&_WebService_serviceDesc, srv)
}

func _WebService_UpdateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).UpdateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talk_cloud.WebService/UpdateGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).UpdateGroup(ctx, req.(*UpdateGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebService_DeleteGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Group)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).DeleteGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talk_cloud.WebService/DeleteGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).DeleteGroup(ctx, req.(*Group))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebService_ImportDeviceByRoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ImportDeviceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).ImportDeviceByRoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talk_cloud.WebService/ImportDeviceByRoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).ImportDeviceByRoot(ctx, req.(*ImportDeviceReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _WebService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "talk_cloud.WebService",
	HandlerType: (*WebServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateGroup",
			Handler:    _WebService_UpdateGroup_Handler,
		},
		{
			MethodName: "DeleteGroup",
			Handler:    _WebService_DeleteGroup_Handler,
		},
		{
			MethodName: "ImportDeviceByRoot",
			Handler:    _WebService_ImportDeviceByRoot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "talk_cloud_web.proto",
}
