// Code generated by protoc-gen-go. DO NOT EDIT.
// source: plugin.proto

package plugin

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
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

type Validator struct {
	Omitempty            bool     `protobuf:"varint,1,opt,name=omitempty,proto3" json:"omitempty,omitempty"`
	Gte                  float64  `protobuf:"fixed64,2,opt,name=gte,proto3" json:"gte,omitempty"`
	Gt                   float64  `protobuf:"fixed64,3,opt,name=gt,proto3" json:"gt,omitempty"`
	Lte                  float64  `protobuf:"fixed64,4,opt,name=lte,proto3" json:"lte,omitempty"`
	Lt                   float64  `protobuf:"fixed64,5,opt,name=lt,proto3" json:"lt,omitempty"`
	Eq                   float64  `protobuf:"fixed64,7,opt,name=eq,proto3" json:"eq,omitempty"`
	StringEq             string   `protobuf:"bytes,8,opt,name=stringEq,proto3" json:"stringEq,omitempty"`
	In                   string   `protobuf:"bytes,9,opt,name=in,proto3" json:"in,omitempty"`
	Regexp               string   `protobuf:"bytes,10,opt,name=regexp,proto3" json:"regexp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Validator) Reset()         { *m = Validator{} }
func (m *Validator) String() string { return proto.CompactTextString(m) }
func (*Validator) ProtoMessage()    {}
func (*Validator) Descriptor() ([]byte, []int) {
	return fileDescriptor_22a625af4bc1cc87, []int{0}
}

func (m *Validator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Validator.Unmarshal(m, b)
}
func (m *Validator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Validator.Marshal(b, m, deterministic)
}
func (m *Validator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Validator.Merge(m, src)
}
func (m *Validator) XXX_Size() int {
	return xxx_messageInfo_Validator.Size(m)
}
func (m *Validator) XXX_DiscardUnknown() {
	xxx_messageInfo_Validator.DiscardUnknown(m)
}

var xxx_messageInfo_Validator proto.InternalMessageInfo

func (m *Validator) GetOmitempty() bool {
	if m != nil {
		return m.Omitempty
	}
	return false
}

func (m *Validator) GetGte() float64 {
	if m != nil {
		return m.Gte
	}
	return 0
}

func (m *Validator) GetGt() float64 {
	if m != nil {
		return m.Gt
	}
	return 0
}

func (m *Validator) GetLte() float64 {
	if m != nil {
		return m.Lte
	}
	return 0
}

func (m *Validator) GetLt() float64 {
	if m != nil {
		return m.Lt
	}
	return 0
}

func (m *Validator) GetEq() float64 {
	if m != nil {
		return m.Eq
	}
	return 0
}

func (m *Validator) GetStringEq() string {
	if m != nil {
		return m.StringEq
	}
	return ""
}

func (m *Validator) GetIn() string {
	if m != nil {
		return m.In
	}
	return ""
}

func (m *Validator) GetRegexp() string {
	if m != nil {
		return m.Regexp
	}
	return ""
}

var E_Validator = &proto.ExtensionDesc{
	ExtendedType:  (*descriptorpb.FieldOptions)(nil),
	ExtensionType: (*Validator)(nil),
	Field:         50002,
	Name:          "plugin.validator",
	Tag:           "bytes,50002,opt,name=validator",
	Filename:      "plugin.proto",
}

func init() {
	proto.RegisterType((*Validator)(nil), "plugin.Validator")
	proto.RegisterExtension(E_Validator)
}

func init() { proto.RegisterFile("plugin.proto", fileDescriptor_22a625af4bc1cc87) }

var fileDescriptor_22a625af4bc1cc87 = []byte{
	// 246 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0x41, 0x4a, 0xc4, 0x30,
	0x14, 0x86, 0x49, 0x47, 0x6b, 0x13, 0x45, 0x34, 0x0b, 0x79, 0x0c, 0x0a, 0xc5, 0x55, 0x57, 0x19,
	0xd0, 0x9d, 0x7b, 0xdd, 0x2a, 0x5d, 0xb8, 0x9f, 0xb1, 0xcf, 0x10, 0xc8, 0x24, 0x69, 0xfa, 0x46,
	0xf4, 0x02, 0xde, 0xcb, 0x2b, 0x78, 0x22, 0x49, 0x32, 0x76, 0x76, 0xfd, 0x3e, 0xbe, 0x86, 0xfc,
	0x11, 0x67, 0xc1, 0xee, 0xb4, 0x71, 0x2a, 0x44, 0x4f, 0x5e, 0xd6, 0x85, 0x96, 0xad, 0xf6, 0x5e,
	0x5b, 0x5c, 0x65, 0xbb, 0xd9, 0xbd, 0xaf, 0x06, 0x9c, 0xde, 0xa2, 0x09, 0xe4, 0x63, 0x29, 0x6f,
	0x7f, 0x98, 0xe0, 0xaf, 0x6b, 0x6b, 0x86, 0x35, 0xf9, 0x28, 0xaf, 0x05, 0xf7, 0x5b, 0x43, 0xb8,
	0x0d, 0xf4, 0x05, 0xac, 0x65, 0x5d, 0xd3, 0x1f, 0x84, 0xbc, 0x10, 0x0b, 0x4d, 0x08, 0x55, 0xcb,
	0x3a, 0xd6, 0xa7, 0x4f, 0x79, 0x2e, 0x2a, 0x4d, 0xb0, 0xc8, 0xa2, 0xd2, 0x94, 0x0a, 0x4b, 0x08,
	0x47, 0xa5, 0xb0, 0xa5, 0xb0, 0x04, 0xc7, 0xa5, 0xb0, 0x94, 0x18, 0x47, 0x38, 0x29, 0x8c, 0xa3,
	0x5c, 0x8a, 0x66, 0xa2, 0x68, 0x9c, 0x7e, 0x1c, 0xa1, 0x69, 0x59, 0xc7, 0xfb, 0x99, 0x53, 0x6b,
	0x1c, 0xf0, 0x6c, 0x2b, 0xe3, 0xe4, 0x95, 0xa8, 0x23, 0x6a, 0xfc, 0x0c, 0x20, 0xb2, 0xdb, 0xd3,
	0xc3, 0x8b, 0xe0, 0x1f, 0xf3, 0x84, 0x1b, 0x55, 0x36, 0xab, 0xff, 0xcd, 0xea, 0xc9, 0xa0, 0x1d,
	0x9e, 0x03, 0x19, 0xef, 0x26, 0xf8, 0xfd, 0x4e, 0xb7, 0x3d, 0xbd, 0xbb, 0x54, 0xfb, 0x07, 0x9b,
	0xc7, 0xf7, 0x87, 0x43, 0x36, 0x75, 0xfe, 0xf9, 0xfe, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x78, 0x07,
	0x49, 0xb4, 0x56, 0x01, 0x00, 0x00,
}
