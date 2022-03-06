// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/veritas/veritas.proto

package veritas

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

type MessageType int32

const (
	MessageType_Approve MessageType = 0
	MessageType_Abort   MessageType = 1
)

var MessageType_name = map[int32]string{
	0: "Approve",
	1: "Abort",
}

var MessageType_value = map[string]int32{
	"Approve": 0,
	"Abort":   1,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{0}
}

type VerifyRequest struct {
	Signature            string   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyRequest) Reset()         { *m = VerifyRequest{} }
func (m *VerifyRequest) String() string { return proto.CompactTextString(m) }
func (*VerifyRequest) ProtoMessage()    {}
func (*VerifyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{0}
}

func (m *VerifyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyRequest.Unmarshal(m, b)
}
func (m *VerifyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyRequest.Marshal(b, m, deterministic)
}
func (m *VerifyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyRequest.Merge(m, src)
}
func (m *VerifyRequest) XXX_Size() int {
	return xxx_messageInfo_VerifyRequest.Size(m)
}
func (m *VerifyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyRequest proto.InternalMessageInfo

func (m *VerifyRequest) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *VerifyRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type VerifyResponse struct {
	Value                 string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Version               int64    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	RootDigest            []byte   `protobuf:"bytes,3,opt,name=root_digest,json=rootDigest,proto3" json:"root_digest,omitempty"`
	SideNodes             [][]byte `protobuf:"bytes,4,rep,name=side_nodes,json=sideNodes,proto3" json:"side_nodes,omitempty"`
	NonMembershipLeafData []byte   `protobuf:"bytes,5,opt,name=non_membership_leaf_data,json=nonMembershipLeafData,proto3" json:"non_membership_leaf_data,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *VerifyResponse) Reset()         { *m = VerifyResponse{} }
func (m *VerifyResponse) String() string { return proto.CompactTextString(m) }
func (*VerifyResponse) ProtoMessage()    {}
func (*VerifyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{1}
}

func (m *VerifyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyResponse.Unmarshal(m, b)
}
func (m *VerifyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyResponse.Marshal(b, m, deterministic)
}
func (m *VerifyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyResponse.Merge(m, src)
}
func (m *VerifyResponse) XXX_Size() int {
	return xxx_messageInfo_VerifyResponse.Size(m)
}
func (m *VerifyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyResponse proto.InternalMessageInfo

func (m *VerifyResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *VerifyResponse) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *VerifyResponse) GetRootDigest() []byte {
	if m != nil {
		return m.RootDigest
	}
	return nil
}

func (m *VerifyResponse) GetSideNodes() [][]byte {
	if m != nil {
		return m.SideNodes
	}
	return nil
}

func (m *VerifyResponse) GetNonMembershipLeafData() []byte {
	if m != nil {
		return m.NonMembershipLeafData
	}
	return nil
}

type SharedLog struct {
	Seq                  int64         `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
	Sets                 []*SetRequest `protobuf:"bytes,2,rep,name=sets,proto3" json:"sets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SharedLog) Reset()         { *m = SharedLog{} }
func (m *SharedLog) String() string { return proto.CompactTextString(m) }
func (*SharedLog) ProtoMessage()    {}
func (*SharedLog) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{2}
}

func (m *SharedLog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SharedLog.Unmarshal(m, b)
}
func (m *SharedLog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SharedLog.Marshal(b, m, deterministic)
}
func (m *SharedLog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SharedLog.Merge(m, src)
}
func (m *SharedLog) XXX_Size() int {
	return xxx_messageInfo_SharedLog.Size(m)
}
func (m *SharedLog) XXX_DiscardUnknown() {
	xxx_messageInfo_SharedLog.DiscardUnknown(m)
}

var xxx_messageInfo_SharedLog proto.InternalMessageInfo

func (m *SharedLog) GetSeq() int64 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *SharedLog) GetSets() []*SetRequest {
	if m != nil {
		return m.Sets
	}
	return nil
}

type Block struct {
	Txs                  []*SharedLog `protobuf:"bytes,1,rep,name=txs,proto3" json:"txs,omitempty"`
	Type                 MessageType  `protobuf:"varint,2,opt,name=type,proto3,enum=controller.MessageType" json:"type,omitempty"`
	Signature            string       `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	BlkId                string       `protobuf:"bytes,4,opt,name=blkId,proto3" json:"blkId,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Block) Reset()         { *m = Block{} }
func (m *Block) String() string { return proto.CompactTextString(m) }
func (*Block) ProtoMessage()    {}
func (*Block) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{3}
}

func (m *Block) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Block.Unmarshal(m, b)
}
func (m *Block) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Block.Marshal(b, m, deterministic)
}
func (m *Block) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Block.Merge(m, src)
}
func (m *Block) XXX_Size() int {
	return xxx_messageInfo_Block.Size(m)
}
func (m *Block) XXX_DiscardUnknown() {
	xxx_messageInfo_Block.DiscardUnknown(m)
}

var xxx_messageInfo_Block proto.InternalMessageInfo

func (m *Block) GetTxs() []*SharedLog {
	if m != nil {
		return m.Txs
	}
	return nil
}

func (m *Block) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_Approve
}

func (m *Block) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *Block) GetBlkId() string {
	if m != nil {
		return m.BlkId
	}
	return ""
}

type GetRequest struct {
	Signature            string   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{4}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *GetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Version              int64    `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{5}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *GetResponse) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

type SetRequest struct {
	Signature            string   `protobuf:"bytes,1,opt,name=signature,proto3" json:"signature,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	Version              int64    `protobuf:"varint,4,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetRequest) Reset()         { *m = SetRequest{} }
func (m *SetRequest) String() string { return proto.CompactTextString(m) }
func (*SetRequest) ProtoMessage()    {}
func (*SetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{6}
}

func (m *SetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetRequest.Unmarshal(m, b)
}
func (m *SetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetRequest.Marshal(b, m, deterministic)
}
func (m *SetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetRequest.Merge(m, src)
}
func (m *SetRequest) XXX_Size() int {
	return xxx_messageInfo_SetRequest.Size(m)
}
func (m *SetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetRequest proto.InternalMessageInfo

func (m *SetRequest) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *SetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SetRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *SetRequest) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

type SetResponse struct {
	Txid                 string   `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetResponse) Reset()         { *m = SetResponse{} }
func (m *SetResponse) String() string { return proto.CompactTextString(m) }
func (*SetResponse) ProtoMessage()    {}
func (*SetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb98f398e219d6a8, []int{7}
}

func (m *SetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetResponse.Unmarshal(m, b)
}
func (m *SetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetResponse.Marshal(b, m, deterministic)
}
func (m *SetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetResponse.Merge(m, src)
}
func (m *SetResponse) XXX_Size() int {
	return xxx_messageInfo_SetResponse.Size(m)
}
func (m *SetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetResponse proto.InternalMessageInfo

func (m *SetResponse) GetTxid() string {
	if m != nil {
		return m.Txid
	}
	return ""
}

func init() {
	proto.RegisterEnum("controller.MessageType", MessageType_name, MessageType_value)
	proto.RegisterType((*VerifyRequest)(nil), "controller.VerifyRequest")
	proto.RegisterType((*VerifyResponse)(nil), "controller.VerifyResponse")
	proto.RegisterType((*SharedLog)(nil), "controller.SharedLog")
	proto.RegisterType((*Block)(nil), "controller.Block")
	proto.RegisterType((*GetRequest)(nil), "controller.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "controller.GetResponse")
	proto.RegisterType((*SetRequest)(nil), "controller.SetRequest")
	proto.RegisterType((*SetResponse)(nil), "controller.SetResponse")
}

func init() { proto.RegisterFile("proto/veritas/veritas.proto", fileDescriptor_bb98f398e219d6a8) }

var fileDescriptor_bb98f398e219d6a8 = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0xc5, 0x5d, 0xbb, 0x55, 0xc6, 0x6d, 0x89, 0x46, 0x2d, 0x5d, 0x02, 0x88, 0x60, 0x09, 0x11,
	0x15, 0x29, 0x48, 0x41, 0x2a, 0x17, 0x50, 0xd5, 0xaa, 0x52, 0x55, 0xa9, 0xe5, 0x60, 0x23, 0x0e,
	0x5c, 0xac, 0x4d, 0x3d, 0x49, 0xad, 0xb8, 0x5e, 0x77, 0x77, 0x13, 0x35, 0x9f, 0xc1, 0xd7, 0x20,
	0xfe, 0x0e, 0x79, 0x13, 0x37, 0x76, 0xa1, 0x07, 0x7a, 0xca, 0xcc, 0xdb, 0x79, 0xf3, 0x26, 0x6f,
	0x3c, 0xf0, 0xa2, 0x50, 0xd2, 0xc8, 0x0f, 0x33, 0x52, 0xa9, 0x11, 0xba, 0xfa, 0xed, 0x5b, 0x14,
	0xe1, 0x52, 0xe6, 0x46, 0xc9, 0x2c, 0x23, 0x15, 0x1c, 0xc2, 0xd6, 0x77, 0x52, 0xe9, 0x68, 0x1e,
	0xd2, 0xcd, 0x94, 0xb4, 0xc1, 0x97, 0xd0, 0xd2, 0xe9, 0x38, 0x17, 0x66, 0xaa, 0x88, 0x3b, 0x5d,
	0xa7, 0xd7, 0x0a, 0x57, 0x00, 0xb6, 0x81, 0x4d, 0x68, 0xce, 0xd7, 0x2c, 0x5e, 0x86, 0xc1, 0x6f,
	0x07, 0xb6, 0xab, 0x0e, 0xba, 0x90, 0xb9, 0x26, 0xdc, 0x01, 0x6f, 0x26, 0xb2, 0x69, 0x45, 0x5f,
	0x24, 0xc8, 0x61, 0x63, 0x46, 0x4a, 0xa7, 0x32, 0xb7, 0x74, 0x16, 0x56, 0x29, 0xbe, 0x06, 0x5f,
	0x49, 0x69, 0xe2, 0x24, 0x1d, 0x93, 0x36, 0x9c, 0x75, 0x9d, 0xde, 0x66, 0x08, 0x25, 0x74, 0x62,
	0x11, 0x7c, 0x05, 0xa0, 0xd3, 0x84, 0xe2, 0x5c, 0x26, 0xa4, 0xb9, 0xdb, 0x65, 0xbd, 0xcd, 0x72,
	0xa8, 0x84, 0xbe, 0x96, 0x00, 0x7e, 0x02, 0x9e, 0xcb, 0x3c, 0xbe, 0xa6, 0xeb, 0x21, 0x29, 0x7d,
	0x95, 0x16, 0x71, 0x46, 0x62, 0x14, 0x27, 0xc2, 0x08, 0xee, 0xd9, 0x66, 0xbb, 0xb9, 0xcc, 0x2f,
	0xee, 0x9e, 0xcf, 0x49, 0x8c, 0x4e, 0x84, 0x11, 0xc1, 0x19, 0xb4, 0xa2, 0x2b, 0xa1, 0x28, 0x39,
	0x97, 0xe3, 0xf2, 0xaf, 0x69, 0xba, 0xb1, 0x33, 0xb3, 0xb0, 0x0c, 0x71, 0x1f, 0x5c, 0x4d, 0x46,
	0xf3, 0xb5, 0x2e, 0xeb, 0xf9, 0x83, 0x67, 0xfd, 0x95, 0x6d, 0xfd, 0x88, 0xcc, 0xd2, 0xb0, 0xd0,
	0xd6, 0x04, 0x3f, 0x1d, 0xf0, 0x8e, 0x33, 0x79, 0x39, 0xc1, 0x77, 0xc0, 0xcc, 0xad, 0xe6, 0x8e,
	0x25, 0xed, 0x36, 0x48, 0x95, 0x56, 0x58, 0x56, 0xe0, 0x7b, 0x70, 0xcd, 0xbc, 0x20, 0xeb, 0xc6,
	0xf6, 0x60, 0xaf, 0x5e, 0x79, 0x41, 0x5a, 0x8b, 0x31, 0x7d, 0x9b, 0x17, 0x14, 0xda, 0xa2, 0xe6,
	0x5a, 0xd8, 0xfd, 0xb5, 0xec, 0x80, 0x37, 0xcc, 0x26, 0x67, 0x09, 0x77, 0x17, 0x8e, 0xdb, 0x24,
	0xf8, 0x0c, 0x70, 0x7a, 0x37, 0xe7, 0x7f, 0x2f, 0xf6, 0x0b, 0xf8, 0x96, 0xfd, 0xb8, 0xa5, 0x06,
	0x19, 0x40, 0xf4, 0x68, 0xf1, 0x95, 0x1a, 0x7b, 0x40, 0xcd, 0x6d, 0xaa, 0xbd, 0x01, 0x3f, 0xaa,
	0x0d, 0x8b, 0xe0, 0x9a, 0xdb, 0x34, 0x59, 0x2a, 0xd9, 0x78, 0xff, 0x2d, 0xf8, 0x35, 0x5b, 0xd1,
	0x87, 0x8d, 0xa3, 0xa2, 0x50, 0x72, 0x46, 0xed, 0x27, 0xd8, 0x02, 0xef, 0x68, 0x28, 0x95, 0x69,
	0x3b, 0x83, 0x5f, 0x0e, 0xb8, 0xe5, 0x67, 0x85, 0x07, 0xc0, 0x4e, 0xc9, 0x60, 0x63, 0xed, 0x2b,
	0x3b, 0x3b, 0x7b, 0x7f, 0xe1, 0x4b, 0xed, 0x03, 0x60, 0xd1, 0x7d, 0x5e, 0xf4, 0x00, 0xaf, 0x3e,
	0xf3, 0x21, 0xac, 0x2f, 0xee, 0x08, 0x9f, 0xd7, 0x4b, 0x1a, 0xd7, 0xd9, 0xe9, 0xfc, 0xeb, 0x69,
	0xd1, 0xe0, 0xf8, 0xe9, 0x8f, 0xad, 0xc6, 0xd5, 0x0f, 0xd7, 0x6d, 0xfa, 0xf1, 0x4f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xaa, 0x55, 0xdc, 0x62, 0x0d, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Verify(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error)
}

type nodeClient struct {
	cc *grpc.ClientConn
}

func NewNodeClient(cc *grpc.ClientConn) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/controller.Node/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, "/controller.Node/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Verify(ctx context.Context, in *VerifyRequest, opts ...grpc.CallOption) (*VerifyResponse, error) {
	out := new(VerifyResponse)
	err := c.cc.Invoke(ctx, "/controller.Node/Verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
type NodeServer interface {
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Verify(context.Context, *VerifyRequest) (*VerifyResponse, error)
}

// UnimplementedNodeServer can be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (*UnimplementedNodeServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedNodeServer) Set(ctx context.Context, req *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (*UnimplementedNodeServer) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Verify not implemented")
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.Node/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.Node/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Verify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Verify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/controller.Node/Verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Verify(ctx, req.(*VerifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "controller.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Node_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _Node_Set_Handler,
		},
		{
			MethodName: "Verify",
			Handler:    _Node_Verify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/veritas/veritas.proto",
}
