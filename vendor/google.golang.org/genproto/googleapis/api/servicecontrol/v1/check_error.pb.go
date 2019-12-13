// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/servicecontrol/v1/check_error.proto

package servicecontrol // import "google.golang.org/genproto/googleapis/api/servicecontrol/v1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Error codes for Check responses.
type CheckError_Code int32

const (
	// This is never used in `CheckResponse`.
	CheckError_ERROR_CODE_UNSPECIFIED CheckError_Code = 0
	// The consumer's project id was not found.
	// Same as [google.rpc.Code.NOT_FOUND][].
	CheckError_NOT_FOUND CheckError_Code = 5
	// The consumer doesn't have access to the specified resource.
	// Same as [google.rpc.Code.PERMISSION_DENIED][].
	CheckError_PERMISSION_DENIED CheckError_Code = 7
	// Quota check failed. Same as [google.rpc.Code.RESOURCE_EXHAUSTED][].
	CheckError_RESOURCE_EXHAUSTED CheckError_Code = 8
	// The consumer hasn't activated the service.
	CheckError_SERVICE_NOT_ACTIVATED CheckError_Code = 104
	// The consumer cannot access the service because billing is disabled.
	CheckError_BILLING_DISABLED CheckError_Code = 107
	// The consumer's project has been marked as deleted (soft deletion).
	CheckError_PROJECT_DELETED CheckError_Code = 108
	// The consumer's project number or id does not represent a valid project.
	CheckError_PROJECT_INVALID CheckError_Code = 114
	// The IP address of the consumer is invalid for the specific consumer
	// project.
	CheckError_IP_ADDRESS_BLOCKED CheckError_Code = 109
	// The referer address of the consumer request is invalid for the specific
	// consumer project.
	CheckError_REFERER_BLOCKED CheckError_Code = 110
	// The client application of the consumer request is invalid for the
	// specific consumer project.
	CheckError_CLIENT_APP_BLOCKED CheckError_Code = 111
	// The consumer's API key is invalid.
	CheckError_API_KEY_INVALID CheckError_Code = 105
	// The consumer's API Key has expired.
	CheckError_API_KEY_EXPIRED CheckError_Code = 112
	// The consumer's API Key was not found in config record.
	CheckError_API_KEY_NOT_FOUND CheckError_Code = 113
	// The backend server for looking up project id/number is unavailable.
	CheckError_NAMESPACE_LOOKUP_UNAVAILABLE CheckError_Code = 300
	// The backend server for checking service status is unavailable.
	CheckError_SERVICE_STATUS_UNAVAILABLE CheckError_Code = 301
	// The backend server for checking billing status is unavailable.
	CheckError_BILLING_STATUS_UNAVAILABLE CheckError_Code = 302
)

var CheckError_Code_name = map[int32]string{
	0:   "ERROR_CODE_UNSPECIFIED",
	5:   "NOT_FOUND",
	7:   "PERMISSION_DENIED",
	8:   "RESOURCE_EXHAUSTED",
	104: "SERVICE_NOT_ACTIVATED",
	107: "BILLING_DISABLED",
	108: "PROJECT_DELETED",
	114: "PROJECT_INVALID",
	109: "IP_ADDRESS_BLOCKED",
	110: "REFERER_BLOCKED",
	111: "CLIENT_APP_BLOCKED",
	105: "API_KEY_INVALID",
	112: "API_KEY_EXPIRED",
	113: "API_KEY_NOT_FOUND",
	300: "NAMESPACE_LOOKUP_UNAVAILABLE",
	301: "SERVICE_STATUS_UNAVAILABLE",
	302: "BILLING_STATUS_UNAVAILABLE",
}
var CheckError_Code_value = map[string]int32{
	"ERROR_CODE_UNSPECIFIED":       0,
	"NOT_FOUND":                    5,
	"PERMISSION_DENIED":            7,
	"RESOURCE_EXHAUSTED":           8,
	"SERVICE_NOT_ACTIVATED":        104,
	"BILLING_DISABLED":             107,
	"PROJECT_DELETED":              108,
	"PROJECT_INVALID":              114,
	"IP_ADDRESS_BLOCKED":           109,
	"REFERER_BLOCKED":              110,
	"CLIENT_APP_BLOCKED":           111,
	"API_KEY_INVALID":              105,
	"API_KEY_EXPIRED":              112,
	"API_KEY_NOT_FOUND":            113,
	"NAMESPACE_LOOKUP_UNAVAILABLE": 300,
	"SERVICE_STATUS_UNAVAILABLE":   301,
	"BILLING_STATUS_UNAVAILABLE":   302,
}

func (x CheckError_Code) String() string {
	return proto.EnumName(CheckError_Code_name, int32(x))
}
func (CheckError_Code) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_check_error_1be1bde99e60d1aa, []int{0, 0}
}

// Defines the errors to be returned in
// [google.api.servicecontrol.v1.CheckResponse.check_errors][google.api.servicecontrol.v1.CheckResponse.check_errors].
type CheckError struct {
	// The error code.
	Code CheckError_Code `protobuf:"varint,1,opt,name=code,proto3,enum=google.api.servicecontrol.v1.CheckError_Code" json:"code,omitempty"`
	// Free-form text providing details on the error cause of the error.
	Detail               string   `protobuf:"bytes,2,opt,name=detail,proto3" json:"detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckError) Reset()         { *m = CheckError{} }
func (m *CheckError) String() string { return proto.CompactTextString(m) }
func (*CheckError) ProtoMessage()    {}
func (*CheckError) Descriptor() ([]byte, []int) {
	return fileDescriptor_check_error_1be1bde99e60d1aa, []int{0}
}
func (m *CheckError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckError.Unmarshal(m, b)
}
func (m *CheckError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckError.Marshal(b, m, deterministic)
}
func (dst *CheckError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckError.Merge(dst, src)
}
func (m *CheckError) XXX_Size() int {
	return xxx_messageInfo_CheckError.Size(m)
}
func (m *CheckError) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckError.DiscardUnknown(m)
}

var xxx_messageInfo_CheckError proto.InternalMessageInfo

func (m *CheckError) GetCode() CheckError_Code {
	if m != nil {
		return m.Code
	}
	return CheckError_ERROR_CODE_UNSPECIFIED
}

func (m *CheckError) GetDetail() string {
	if m != nil {
		return m.Detail
	}
	return ""
}

func init() {
	proto.RegisterType((*CheckError)(nil), "google.api.servicecontrol.v1.CheckError")
	proto.RegisterEnum("google.api.servicecontrol.v1.CheckError_Code", CheckError_Code_name, CheckError_Code_value)
}

func init() {
	proto.RegisterFile("google/api/servicecontrol/v1/check_error.proto", fileDescriptor_check_error_1be1bde99e60d1aa)
}

var fileDescriptor_check_error_1be1bde99e60d1aa = []byte{
	// 484 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xdd, 0x6e, 0xd3, 0x3e,
	0x18, 0xc6, 0xff, 0xe9, 0xbf, 0x0c, 0x66, 0x09, 0x16, 0x0c, 0xab, 0x46, 0x55, 0x89, 0xb2, 0xa3,
	0x9d, 0x90, 0x68, 0x70, 0xc8, 0x91, 0x6b, 0xbf, 0x05, 0xaf, 0x59, 0x62, 0xd9, 0x49, 0x35, 0x38,
	0xb1, 0x42, 0x1a, 0x65, 0xd1, 0xba, 0xb8, 0xa4, 0x51, 0xaf, 0x80, 0x0b, 0xe0, 0x2a, 0x38, 0x02,
	0xae, 0x8f, 0x43, 0xe4, 0x76, 0xfd, 0x92, 0xa6, 0x1d, 0xfa, 0x79, 0x7f, 0xcf, 0x63, 0xbd, 0x1f,
	0xc8, 0x2b, 0x8c, 0x29, 0xa6, 0xb9, 0x9f, 0xce, 0x4a, 0x7f, 0x9e, 0xd7, 0x8b, 0x32, 0xcb, 0x33,
	0x53, 0x35, 0xb5, 0x99, 0xfa, 0x8b, 0x73, 0x3f, 0xbb, 0xce, 0xb3, 0x1b, 0x9d, 0xd7, 0xb5, 0xa9,
	0xbd, 0x59, 0x6d, 0x1a, 0x83, 0x7b, 0x2b, 0xde, 0x4b, 0x67, 0xa5, 0xb7, 0xcf, 0x7b, 0x8b, 0xf3,
	0x6e, 0x6f, 0x27, 0x2d, 0xad, 0x2a, 0xd3, 0xa4, 0x4d, 0x69, 0xaa, 0xf9, 0xca, 0x7b, 0xfa, 0xa3,
	0x8d, 0x10, 0xb5, 0x89, 0x60, 0x03, 0x31, 0x41, 0xed, 0xcc, 0x4c, 0xf2, 0x13, 0xa7, 0xef, 0x9c,
	0x3d, 0x7b, 0xf7, 0xd6, 0x7b, 0x28, 0xd9, 0xdb, 0xfa, 0x3c, 0x6a, 0x26, 0xb9, 0x5c, 0x5a, 0x71,
	0x07, 0x1d, 0x4c, 0xf2, 0x26, 0x2d, 0xa7, 0x27, 0xad, 0xbe, 0x73, 0x76, 0x28, 0xef, 0x5e, 0xa7,
	0x3f, 0xff, 0x47, 0x6d, 0x8b, 0xe1, 0x2e, 0xea, 0x80, 0x94, 0x91, 0xd4, 0x34, 0x62, 0xa0, 0x93,
	0x50, 0x09, 0xa0, 0x7c, 0xc8, 0x81, 0xb9, 0xff, 0xe1, 0xa7, 0xe8, 0x30, 0x8c, 0x62, 0x3d, 0x8c,
	0x92, 0x90, 0xb9, 0x8f, 0xf0, 0x31, 0x7a, 0x2e, 0x40, 0x5e, 0x72, 0xa5, 0x78, 0x14, 0x6a, 0x06,
	0xa1, 0xa5, 0x1e, 0xe3, 0x0e, 0xc2, 0x12, 0x54, 0x94, 0x48, 0x0a, 0x1a, 0xae, 0x3e, 0x91, 0x44,
	0xc5, 0xc0, 0xdc, 0x27, 0xf8, 0x15, 0x3a, 0x56, 0x20, 0xc7, 0x9c, 0x82, 0xb6, 0x29, 0x84, 0xc6,
	0x7c, 0x4c, 0x6c, 0xe9, 0x1a, 0xbf, 0x44, 0xee, 0x80, 0x07, 0x01, 0x0f, 0x3f, 0x6a, 0xc6, 0x15,
	0x19, 0x04, 0xc0, 0xdc, 0x1b, 0xfc, 0x02, 0x1d, 0x09, 0x19, 0x5d, 0x00, 0x8d, 0x35, 0x83, 0x00,
	0x2c, 0x3a, 0xdd, 0x15, 0x79, 0x38, 0x26, 0x01, 0x67, 0x6e, 0x6d, 0xbf, 0xe4, 0x42, 0x13, 0xc6,
	0x24, 0x28, 0xa5, 0x07, 0x41, 0x44, 0x47, 0xc0, 0xdc, 0x5b, 0x0b, 0x4b, 0x18, 0x82, 0x04, 0xb9,
	0x11, 0x2b, 0x0b, 0xd3, 0x80, 0x43, 0x18, 0x6b, 0x22, 0xc4, 0x46, 0x37, 0x16, 0x26, 0x82, 0xeb,
	0x11, 0x7c, 0xde, 0x24, 0x97, 0xbb, 0x22, 0x5c, 0x09, 0x2e, 0x81, 0xb9, 0x33, 0xdb, 0xf8, 0x5a,
	0xdc, 0xce, 0xe3, 0x1b, 0x7e, 0x83, 0x7a, 0x21, 0xb9, 0x04, 0x25, 0x08, 0x05, 0x1d, 0x44, 0xd1,
	0x28, 0x11, 0x3a, 0x09, 0xc9, 0x98, 0xf0, 0xc0, 0xb6, 0xe4, 0xfe, 0x6a, 0xe1, 0xd7, 0xa8, 0xbb,
	0x9e, 0x81, 0x8a, 0x49, 0x9c, 0xa8, 0x3d, 0xe0, 0xf7, 0x12, 0x58, 0x4f, 0xe2, 0x1e, 0xe0, 0x4f,
	0x6b, 0xf0, 0xdd, 0x41, 0xfd, 0xcc, 0xdc, 0x3e, 0xb8, 0xfb, 0xc1, 0xd1, 0x76, 0xf9, 0xc2, 0x1e,
	0x92, 0x70, 0xbe, 0x5c, 0xdc, 0x19, 0x0a, 0x33, 0x4d, 0xab, 0xc2, 0x33, 0x75, 0xe1, 0x17, 0x79,
	0xb5, 0x3c, 0x33, 0x7f, 0x55, 0x4a, 0x67, 0xe5, 0xfc, 0xfe, 0xab, 0xfe, 0xb0, 0xaf, 0xfc, 0x75,
	0x9c, 0xaf, 0x07, 0x4b, 0xe7, 0xfb, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x65, 0x26, 0xbf,
	0x0e, 0x03, 0x00, 0x00,
}
