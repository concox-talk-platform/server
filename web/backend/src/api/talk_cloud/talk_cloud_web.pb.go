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

type UpdDInfoReq struct {
	DeviceInfo           *DeviceUpdate `protobuf:"bytes,1,opt,name=deviceInfo,proto3" json:"deviceInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *UpdDInfoReq) Reset()         { *m = UpdDInfoReq{} }
func (m *UpdDInfoReq) String() string { return proto.CompactTextString(m) }
func (*UpdDInfoReq) ProtoMessage()    {}
func (*UpdDInfoReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{0}
}

func (m *UpdDInfoReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdDInfoReq.Unmarshal(m, b)
}
func (m *UpdDInfoReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdDInfoReq.Marshal(b, m, deterministic)
}
func (m *UpdDInfoReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdDInfoReq.Merge(m, src)
}
func (m *UpdDInfoReq) XXX_Size() int {
	return xxx_messageInfo_UpdDInfoReq.Size(m)
}
func (m *UpdDInfoReq) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdDInfoReq.DiscardUnknown(m)
}

var xxx_messageInfo_UpdDInfoReq proto.InternalMessageInfo

func (m *UpdDInfoReq) GetDeviceInfo() *DeviceUpdate {
	if m != nil {
		return m.DeviceInfo
	}
	return nil
}

type UpdDInfoResp struct {
	Res                  *Result  `protobuf:"bytes,1,opt,name=res,proto3" json:"res,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdDInfoResp) Reset()         { *m = UpdDInfoResp{} }
func (m *UpdDInfoResp) String() string { return proto.CompactTextString(m) }
func (*UpdDInfoResp) ProtoMessage()    {}
func (*UpdDInfoResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{1}
}

func (m *UpdDInfoResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdDInfoResp.Unmarshal(m, b)
}
func (m *UpdDInfoResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdDInfoResp.Marshal(b, m, deterministic)
}
func (m *UpdDInfoResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdDInfoResp.Merge(m, src)
}
func (m *UpdDInfoResp) XXX_Size() int {
	return xxx_messageInfo_UpdDInfoResp.Size(m)
}
func (m *UpdDInfoResp) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdDInfoResp.DiscardUnknown(m)
}

var xxx_messageInfo_UpdDInfoResp proto.InternalMessageInfo

func (m *UpdDInfoResp) GetRes() *Result {
	if m != nil {
		return m.Res
	}
	return nil
}

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
	return fileDescriptor_bee75d1763415d2b, []int{2}
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
	return fileDescriptor_bee75d1763415d2b, []int{3}
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
	return fileDescriptor_bee75d1763415d2b, []int{4}
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
	AccountId            int32         `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	Devices              []*DeviceInfo `protobuf:"bytes,2,rep,name=Devices,proto3" json:"Devices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *ImportDeviceReq) Reset()         { *m = ImportDeviceReq{} }
func (m *ImportDeviceReq) String() string { return proto.CompactTextString(m) }
func (*ImportDeviceReq) ProtoMessage()    {}
func (*ImportDeviceReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{5}
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

func (m *ImportDeviceReq) GetDevices() []*DeviceInfo {
	if m != nil {
		return m.Devices
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
	return fileDescriptor_bee75d1763415d2b, []int{6}
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

type DeviceInfo struct {
	Id                   int32    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	IMei                 string   `protobuf:"bytes,2,opt,name=IMei,proto3" json:"IMei,omitempty"`
	UserName             string   `protobuf:"bytes,3,opt,name=UserName,proto3" json:"UserName,omitempty"`
	PassWord             string   `protobuf:"bytes,4,opt,name=PassWord,proto3" json:"PassWord,omitempty"`
	AccountId            int32    `protobuf:"varint,5,opt,name=AccountId,proto3" json:"AccountId,omitempty"`
	CreateTime           string   `protobuf:"bytes,6,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	LLTime               string   `protobuf:"bytes,7,opt,name=LLTime,proto3" json:"LLTime,omitempty"`
	ChangeTime           string   `protobuf:"bytes,8,opt,name=ChangeTime,proto3" json:"ChangeTime,omitempty"`
	LocalTime            uint64   `protobuf:"varint,9,opt,name=LocalTime,proto3" json:"LocalTime,omitempty"`
	Speed                float32  `protobuf:"fixed32,10,opt,name=Speed,proto3" json:"Speed,omitempty"`
	Course               float32  `protobuf:"fixed32,11,opt,name=Course,proto3" json:"Course,omitempty"`
	DeviceType           string   `protobuf:"bytes,12,opt,name=DeviceType,proto3" json:"DeviceType,omitempty"`
	ActiveTime           string   `protobuf:"bytes,13,opt,name=ActiveTime,proto3" json:"ActiveTime,omitempty"`
	SaleTime             string   `protobuf:"bytes,14,opt,name=SaleTime,proto3" json:"SaleTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeviceInfo) Reset()         { *m = DeviceInfo{} }
func (m *DeviceInfo) String() string { return proto.CompactTextString(m) }
func (*DeviceInfo) ProtoMessage()    {}
func (*DeviceInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{7}
}

func (m *DeviceInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceInfo.Unmarshal(m, b)
}
func (m *DeviceInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceInfo.Marshal(b, m, deterministic)
}
func (m *DeviceInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceInfo.Merge(m, src)
}
func (m *DeviceInfo) XXX_Size() int {
	return xxx_messageInfo_DeviceInfo.Size(m)
}
func (m *DeviceInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceInfo proto.InternalMessageInfo

func (m *DeviceInfo) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *DeviceInfo) GetIMei() string {
	if m != nil {
		return m.IMei
	}
	return ""
}

func (m *DeviceInfo) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *DeviceInfo) GetPassWord() string {
	if m != nil {
		return m.PassWord
	}
	return ""
}

func (m *DeviceInfo) GetAccountId() int32 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *DeviceInfo) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

func (m *DeviceInfo) GetLLTime() string {
	if m != nil {
		return m.LLTime
	}
	return ""
}

func (m *DeviceInfo) GetChangeTime() string {
	if m != nil {
		return m.ChangeTime
	}
	return ""
}

func (m *DeviceInfo) GetLocalTime() uint64 {
	if m != nil {
		return m.LocalTime
	}
	return 0
}

func (m *DeviceInfo) GetSpeed() float32 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *DeviceInfo) GetCourse() float32 {
	if m != nil {
		return m.Course
	}
	return 0
}

func (m *DeviceInfo) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *DeviceInfo) GetActiveTime() string {
	if m != nil {
		return m.ActiveTime
	}
	return ""
}

func (m *DeviceInfo) GetSaleTime() string {
	if m != nil {
		return m.SaleTime
	}
	return ""
}

type DeviceUpdate struct {
	Id                   int32    `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	IMei                 string   `protobuf:"bytes,2,opt,name=IMei,proto3" json:"IMei,omitempty"`
	NickName             string   `protobuf:"bytes,3,opt,name=NickName,proto3" json:"NickName,omitempty"`
	LoginId              int32    `protobuf:"varint,4,opt,name=LoginId,proto3" json:"LoginId,omitempty"`
	CreateTime           string   `protobuf:"bytes,6,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	LLTime               string   `protobuf:"bytes,7,opt,name=LLTime,proto3" json:"LLTime,omitempty"`
	ChangeTime           string   `protobuf:"bytes,8,opt,name=ChangeTime,proto3" json:"ChangeTime,omitempty"`
	LocalTime            uint64   `protobuf:"varint,9,opt,name=LocalTime,proto3" json:"LocalTime,omitempty"`
	DeviceType           string   `protobuf:"bytes,12,opt,name=DeviceType,proto3" json:"DeviceType,omitempty"`
	ActiveTime           string   `protobuf:"bytes,13,opt,name=ActiveTime,proto3" json:"ActiveTime,omitempty"`
	SaleTime             string   `protobuf:"bytes,14,opt,name=SaleTime,proto3" json:"SaleTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeviceUpdate) Reset()         { *m = DeviceUpdate{} }
func (m *DeviceUpdate) String() string { return proto.CompactTextString(m) }
func (*DeviceUpdate) ProtoMessage()    {}
func (*DeviceUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_bee75d1763415d2b, []int{8}
}

func (m *DeviceUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceUpdate.Unmarshal(m, b)
}
func (m *DeviceUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceUpdate.Marshal(b, m, deterministic)
}
func (m *DeviceUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceUpdate.Merge(m, src)
}
func (m *DeviceUpdate) XXX_Size() int {
	return xxx_messageInfo_DeviceUpdate.Size(m)
}
func (m *DeviceUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceUpdate proto.InternalMessageInfo

func (m *DeviceUpdate) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *DeviceUpdate) GetIMei() string {
	if m != nil {
		return m.IMei
	}
	return ""
}

func (m *DeviceUpdate) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *DeviceUpdate) GetLoginId() int32 {
	if m != nil {
		return m.LoginId
	}
	return 0
}

func (m *DeviceUpdate) GetCreateTime() string {
	if m != nil {
		return m.CreateTime
	}
	return ""
}

func (m *DeviceUpdate) GetLLTime() string {
	if m != nil {
		return m.LLTime
	}
	return ""
}

func (m *DeviceUpdate) GetChangeTime() string {
	if m != nil {
		return m.ChangeTime
	}
	return ""
}

func (m *DeviceUpdate) GetLocalTime() uint64 {
	if m != nil {
		return m.LocalTime
	}
	return 0
}

func (m *DeviceUpdate) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *DeviceUpdate) GetActiveTime() string {
	if m != nil {
		return m.ActiveTime
	}
	return ""
}

func (m *DeviceUpdate) GetSaleTime() string {
	if m != nil {
		return m.SaleTime
	}
	return ""
}

func init() {
	proto.RegisterType((*UpdDInfoReq)(nil), "talk_cloud.UpdDInfoReq")
	proto.RegisterType((*UpdDInfoResp)(nil), "talk_cloud.UpdDInfoResp")
	proto.RegisterType((*UpdateGroupReq)(nil), "talk_cloud.UpdateGroupReq")
	proto.RegisterType((*UpdateGroupResp)(nil), "talk_cloud.UpdateGroupResp")
	proto.RegisterType((*DeleteGroupResp)(nil), "talk_cloud.DeleteGroupResp")
	proto.RegisterType((*ImportDeviceReq)(nil), "talk_cloud.ImportDeviceReq")
	proto.RegisterType((*ImportDeviceResp)(nil), "talk_cloud.ImportDeviceResp")
	proto.RegisterType((*DeviceInfo)(nil), "talk_cloud.DeviceInfo")
	proto.RegisterType((*DeviceUpdate)(nil), "talk_cloud.DeviceUpdate")
}

func init() { proto.RegisterFile("talk_cloud_web.proto", fileDescriptor_bee75d1763415d2b) }

var fileDescriptor_bee75d1763415d2b = []byte{
	// 615 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x55, 0x5f, 0x4f, 0x13, 0x41,
	0x10, 0xa7, 0x77, 0xb4, 0xd0, 0x29, 0x02, 0x4e, 0x08, 0x6e, 0x2a, 0x31, 0xcd, 0xc5, 0x87, 0xc6,
	0x07, 0x24, 0xc8, 0x83, 0x4f, 0x1a, 0x84, 0x04, 0x9b, 0x14, 0xa2, 0x0b, 0x0d, 0x8f, 0xe4, 0x7a,
	0x37, 0xd6, 0x86, 0x6b, 0xf7, 0xbc, 0xbd, 0x62, 0xf8, 0x02, 0x7e, 0x03, 0x3f, 0x89, 0x6f, 0x7e,
	0x3a, 0xb3, 0xbb, 0x77, 0xb7, 0xdb, 0x22, 0xc4, 0xc4, 0x18, 0xdf, 0x3a, 0xbf, 0xdf, 0xcc, 0x6f,
	0xfe, 0x75, 0xe7, 0x60, 0x2b, 0x0f, 0x93, 0xeb, 0xab, 0x28, 0x11, 0xb3, 0xf8, 0xea, 0x2b, 0x0d,
	0x77, 0xd3, 0x4c, 0xe4, 0x02, 0xc1, 0xa2, 0xed, 0x6d, 0xc7, 0x63, 0x22, 0x62, 0x4a, 0x8c, 0x4f,
	0x70, 0x02, 0xad, 0x41, 0x1a, 0x1f, 0xf7, 0xa6, 0x9f, 0x04, 0xa7, 0x2f, 0xf8, 0x1a, 0x20, 0xa6,
	0x9b, 0x71, 0x44, 0x0a, 0x60, 0xb5, 0x4e, 0xad, 0xdb, 0xda, 0x67, 0xbb, 0x36, 0x76, 0xf7, 0x58,
	0xb3, 0x83, 0x34, 0x0e, 0x73, 0xe2, 0x8e, 0x6f, 0x70, 0x00, 0x6b, 0x56, 0x48, 0xa6, 0xf8, 0x1c,
	0xfc, 0x8c, 0x64, 0x21, 0x81, 0xae, 0x04, 0x27, 0x39, 0x4b, 0x72, 0xae, 0xe8, 0xe0, 0x7b, 0x0d,
	0xd6, 0x8d, 0xd8, 0x49, 0x26, 0x66, 0xa9, 0x2a, 0x61, 0x07, 0x9a, 0x85, 0x6c, 0x2c, 0x99, 0xd7,
	0xf1, 0xbb, 0x3e, 0xb7, 0x00, 0x1e, 0x40, 0xcb, 0x26, 0x95, 0xcc, 0xef, 0xf8, 0x8b, 0xf2, 0xa7,
	0x34, 0x19, 0x52, 0xc6, 0x5d, 0x37, 0x7c, 0x09, 0xcd, 0x91, 0xd2, 0xd7, 0x5d, 0x2d, 0xeb, 0x92,
	0x1e, 0xbb, 0x31, 0x26, 0xb9, 0xf5, 0x09, 0x8e, 0x60, 0x63, 0xae, 0x2c, 0x99, 0xe2, 0x1e, 0x34,
	0x33, 0x5d, 0xf9, 0xa9, 0x1c, 0x3d, 0xd0, 0x96, 0x75, 0x52, 0x22, 0xc7, 0x94, 0xd0, 0xdf, 0x89,
	0x84, 0xb0, 0xd1, 0x9b, 0xa4, 0x22, 0xcb, 0xcd, 0xe4, 0x8b, 0x09, 0x85, 0x51, 0x24, 0x66, 0xd3,
	0xbc, 0x17, 0x6b, 0x91, 0x3a, 0xb7, 0x00, 0xee, 0xc1, 0x8a, 0x71, 0x35, 0xd3, 0x6b, 0xed, 0x6f,
	0xdf, 0xdd, 0x9f, 0xde, 0x52, 0xe9, 0x16, 0xbc, 0x81, 0xcd, 0xf9, 0x14, 0x32, 0xc5, 0x17, 0xd0,
	0x30, 0x35, 0x3c, 0x50, 0x65, 0xe1, 0x11, 0x7c, 0xf3, 0x01, 0xac, 0x2e, 0xae, 0x83, 0x57, 0xd5,
	0xe5, 0xf5, 0x62, 0x44, 0x58, 0xee, 0x9d, 0xd2, 0x98, 0x79, 0x9d, 0x5a, 0xb7, 0xc9, 0xf5, 0x6f,
	0x6c, 0xc3, 0xea, 0x40, 0x52, 0x76, 0x16, 0x4e, 0x88, 0xf9, 0x1a, 0xaf, 0x6c, 0xc5, 0x7d, 0x08,
	0xa5, 0xbc, 0x14, 0x59, 0xac, 0x77, 0xd5, 0xe4, 0x95, 0xad, 0x5a, 0x3f, 0xac, 0x5a, 0xaf, 0x9b,
	0xd6, 0x2b, 0x00, 0x9f, 0x01, 0x1c, 0x65, 0x14, 0xe6, 0x74, 0x31, 0x9e, 0x10, 0x6b, 0xe8, 0x58,
	0x07, 0xc1, 0x6d, 0x68, 0xf4, 0xfb, 0x9a, 0x5b, 0xd1, 0x5c, 0x61, 0xe9, 0xb8, 0xcf, 0xe1, 0x74,
	0x64, 0xe2, 0x56, 0x8b, 0xb8, 0x0a, 0x51, 0x59, 0xfb, 0x22, 0x0a, 0x13, 0x4d, 0x37, 0x3b, 0xb5,
	0xee, 0x32, 0xb7, 0x00, 0x6e, 0x41, 0xfd, 0x3c, 0x25, 0x8a, 0x19, 0x74, 0x6a, 0x5d, 0x8f, 0x1b,
	0x43, 0xe5, 0x3a, 0x12, 0xb3, 0x4c, 0x12, 0x6b, 0x69, 0xb8, 0xb0, 0x54, 0x2e, 0x33, 0xab, 0x8b,
	0xdb, 0x94, 0xd8, 0x9a, 0xc9, 0x65, 0x11, 0xc5, 0x1f, 0x46, 0xf9, 0xf8, 0xc6, 0xd4, 0xf2, 0xc8,
	0xf0, 0x16, 0x51, 0xd3, 0x39, 0x0f, 0x13, 0xc3, 0xae, 0x9b, 0xe9, 0x94, 0x76, 0xf0, 0xd3, 0x83,
	0x35, 0xf7, 0x81, 0xfe, 0xe9, 0x2a, 0xce, 0xc6, 0xd1, 0xb5, 0xbb, 0x8a, 0xd2, 0x46, 0x06, 0x2b,
	0x7d, 0x31, 0x1a, 0x4f, 0x7b, 0x66, 0x13, 0x75, 0x5e, 0x9a, 0xff, 0x69, 0xd4, 0xff, 0x70, 0x78,
	0xfb, 0x3f, 0x3c, 0x80, 0x4b, 0x1a, 0x9e, 0x53, 0xa6, 0xe4, 0xf0, 0xbd, 0x3e, 0x8c, 0xe5, 0x05,
	0xc0, 0xb6, 0xfb, 0xff, 0x9f, 0xbf, 0x58, 0xed, 0xa7, 0xf7, 0x72, 0x32, 0x0d, 0x96, 0xf0, 0x2d,
	0xb4, 0x9c, 0x33, 0x80, 0x77, 0x0f, 0xcf, 0xbc, 0xc0, 0xc2, 0xc9, 0x08, 0x96, 0xf0, 0x23, 0xa0,
	0xfb, 0x3e, 0xdf, 0xdd, 0x72, 0x21, 0x72, 0x9c, 0x0b, 0x5a, 0x38, 0x11, 0xed, 0x9d, 0xfb, 0x49,
	0x2d, 0x79, 0x02, 0x9b, 0xa6, 0x50, 0xe7, 0xdd, 0x3e, 0x59, 0x68, 0xa3, 0xfc, 0x28, 0xb4, 0xd9,
	0xef, 0x09, 0x25, 0x34, 0x6c, 0xe8, 0xcf, 0xc8, 0xab, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xcc,
	0x02, 0x2f, 0x1f, 0x82, 0x06, 0x00, 0x00,
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
	UpdateDeviceInfo(ctx context.Context, in *UpdDInfoReq, opts ...grpc.CallOption) (*UpdDInfoResp, error)
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

func (c *webServiceClient) UpdateDeviceInfo(ctx context.Context, in *UpdDInfoReq, opts ...grpc.CallOption) (*UpdDInfoResp, error) {
	out := new(UpdDInfoResp)
	err := c.cc.Invoke(ctx, "/talk_cloud.WebService/UpdateDeviceInfo", in, out, opts...)
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
	UpdateDeviceInfo(context.Context, *UpdDInfoReq) (*UpdDInfoResp, error)
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
func (*UnimplementedWebServiceServer) UpdateDeviceInfo(ctx context.Context, req *UpdDInfoReq) (*UpdDInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDeviceInfo not implemented")
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

func _WebService_UpdateDeviceInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdDInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServiceServer).UpdateDeviceInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/talk_cloud.WebService/UpdateDeviceInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServiceServer).UpdateDeviceInfo(ctx, req.(*UpdDInfoReq))
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
		{
			MethodName: "UpdateDeviceInfo",
			Handler:    _WebService_UpdateDeviceInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "talk_cloud_web.proto",
}