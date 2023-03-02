// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// NOTE: Protocol buffers used variable-length encoding, so even though uint64
// is arguably much bigger than we'd need, it's simpler to just use that. We
// do, however, use signed ints for disk usage and limit because those fields
// are optional (and should not be reported if they are < 0).
type PublishRequest struct {
	UnixTime             int64    `protobuf:"varint,1,opt,name=unix_time,json=unixTime,proto3" json:"unix_time,omitempty"`
	Running              bool     `protobuf:"varint,2,opt,name=running,proto3" json:"running,omitempty"`
	MilliCpuUsage        uint64   `protobuf:"varint,3,opt,name=milli_cpu_usage,json=milliCpuUsage,proto3" json:"milli_cpu_usage,omitempty"`
	MemoryTotalMb        uint64   `protobuf:"varint,4,opt,name=memory_total_mb,json=memoryTotalMb,proto3" json:"memory_total_mb,omitempty"`
	MemoryRssMb          uint64   `protobuf:"varint,5,opt,name=memory_rss_mb,json=memoryRssMb,proto3" json:"memory_rss_mb,omitempty"`
	MemoryLimitMb        uint64   `protobuf:"varint,6,opt,name=memory_limit_mb,json=memoryLimitMb,proto3" json:"memory_limit_mb,omitempty"`
	DiskUsageMb          int64    `protobuf:"zigzag64,7,opt,name=disk_usage_mb,json=diskUsageMb,proto3" json:"disk_usage_mb,omitempty"`
	DiskLimitMb          int64    `protobuf:"zigzag64,8,opt,name=disk_limit_mb,json=diskLimitMb,proto3" json:"disk_limit_mb,omitempty"`
	DiskReadKbps         uint64   `protobuf:"varint,9,opt,name=disk_read_kbps,json=diskReadKbps,proto3" json:"disk_read_kbps,omitempty"`
	DiskWriteKbps        uint64   `protobuf:"varint,10,opt,name=disk_write_kbps,json=diskWriteKbps,proto3" json:"disk_write_kbps,omitempty"`
	DiskReadIops         uint64   `protobuf:"varint,11,opt,name=disk_read_iops,json=diskReadIops,proto3" json:"disk_read_iops,omitempty"`
	DiskWriteIops        uint64   `protobuf:"varint,12,opt,name=disk_write_iops,json=diskWriteIops,proto3" json:"disk_write_iops,omitempty"`
	PidsCurrent          uint64   `protobuf:"varint,13,opt,name=pids_current,json=pidsCurrent,proto3" json:"pids_current,omitempty"`
	PidsLimit            uint64   `protobuf:"varint,14,opt,name=pids_limit,json=pidsLimit,proto3" json:"pids_limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublishRequest) Reset()         { *m = PublishRequest{} }
func (m *PublishRequest) String() string { return proto.CompactTextString(m) }
func (*PublishRequest) ProtoMessage()    {}
func (*PublishRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *PublishRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublishRequest.Unmarshal(m, b)
}
func (m *PublishRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublishRequest.Marshal(b, m, deterministic)
}
func (m *PublishRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublishRequest.Merge(m, src)
}
func (m *PublishRequest) XXX_Size() int {
	return xxx_messageInfo_PublishRequest.Size(m)
}
func (m *PublishRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PublishRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PublishRequest proto.InternalMessageInfo

func (m *PublishRequest) GetUnixTime() int64 {
	if m != nil {
		return m.UnixTime
	}
	return 0
}

func (m *PublishRequest) GetRunning() bool {
	if m != nil {
		return m.Running
	}
	return false
}

func (m *PublishRequest) GetMilliCpuUsage() uint64 {
	if m != nil {
		return m.MilliCpuUsage
	}
	return 0
}

func (m *PublishRequest) GetMemoryTotalMb() uint64 {
	if m != nil {
		return m.MemoryTotalMb
	}
	return 0
}

func (m *PublishRequest) GetMemoryRssMb() uint64 {
	if m != nil {
		return m.MemoryRssMb
	}
	return 0
}

func (m *PublishRequest) GetMemoryLimitMb() uint64 {
	if m != nil {
		return m.MemoryLimitMb
	}
	return 0
}

func (m *PublishRequest) GetDiskUsageMb() int64 {
	if m != nil {
		return m.DiskUsageMb
	}
	return 0
}

func (m *PublishRequest) GetDiskLimitMb() int64 {
	if m != nil {
		return m.DiskLimitMb
	}
	return 0
}

func (m *PublishRequest) GetDiskReadKbps() uint64 {
	if m != nil {
		return m.DiskReadKbps
	}
	return 0
}

func (m *PublishRequest) GetDiskWriteKbps() uint64 {
	if m != nil {
		return m.DiskWriteKbps
	}
	return 0
}

func (m *PublishRequest) GetDiskReadIops() uint64 {
	if m != nil {
		return m.DiskReadIops
	}
	return 0
}

func (m *PublishRequest) GetDiskWriteIops() uint64 {
	if m != nil {
		return m.DiskWriteIops
	}
	return 0
}

func (m *PublishRequest) GetPidsCurrent() uint64 {
	if m != nil {
		return m.PidsCurrent
	}
	return 0
}

func (m *PublishRequest) GetPidsLimit() uint64 {
	if m != nil {
		return m.PidsLimit
	}
	return 0
}

type PublishResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PublishResponse) Reset()         { *m = PublishResponse{} }
func (m *PublishResponse) String() string { return proto.CompactTextString(m) }
func (*PublishResponse) ProtoMessage()    {}
func (*PublishResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *PublishResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PublishResponse.Unmarshal(m, b)
}
func (m *PublishResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PublishResponse.Marshal(b, m, deterministic)
}
func (m *PublishResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PublishResponse.Merge(m, src)
}
func (m *PublishResponse) XXX_Size() int {
	return xxx_messageInfo_PublishResponse.Size(m)
}
func (m *PublishResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PublishResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PublishResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*PublishRequest)(nil), "PublishRequest")
	proto.RegisterType((*PublishResponse)(nil), "PublishResponse")
}

func init() {
	proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c)
}

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x4b, 0x4b, 0xeb, 0x40,
	0x14, 0x80, 0x6f, 0x6e, 0x1f, 0x69, 0xa6, 0xaf, 0x7b, 0x67, 0x35, 0x28, 0x42, 0x0c, 0x22, 0x59,
	0x65, 0xa1, 0x4b, 0xdd, 0x68, 0x57, 0xa2, 0x05, 0x09, 0x15, 0xc1, 0xcd, 0x90, 0x34, 0x43, 0x1c,
	0x9a, 0x64, 0xc6, 0x79, 0xa0, 0xfe, 0x04, 0xff, 0xb5, 0xcc, 0x49, 0x5b, 0xdb, 0xba, 0xcc, 0xc7,
	0x97, 0x2f, 0xe7, 0x84, 0x83, 0x82, 0x4c, 0xf2, 0x44, 0x2a, 0x61, 0x44, 0xf4, 0xd5, 0x45, 0x93,
	0x47, 0x9b, 0x57, 0x5c, 0xbf, 0xa6, 0xec, 0xcd, 0x32, 0x6d, 0xf0, 0x31, 0x0a, 0x6c, 0xc3, 0x3f,
	0xa8, 0xe1, 0x35, 0x23, 0x5e, 0xe8, 0xc5, 0x9d, 0x74, 0xe0, 0xc0, 0x82, 0xd7, 0x0c, 0x13, 0xe4,
	0x2b, 0xdb, 0x34, 0xbc, 0x29, 0xc9, 0xdf, 0xd0, 0x8b, 0x07, 0xe9, 0xe6, 0x11, 0x9f, 0xa3, 0x69,
	0xcd, 0xab, 0x8a, 0xd3, 0xa5, 0xb4, 0xd4, 0xea, 0xac, 0x64, 0xa4, 0x13, 0x7a, 0x71, 0x37, 0x1d,
	0x03, 0x9e, 0x49, 0xfb, 0xe4, 0x20, 0x78, 0xac, 0x16, 0xea, 0x93, 0x1a, 0x61, 0xb2, 0x8a, 0xd6,
	0x39, 0xe9, 0xae, 0x3d, 0xc0, 0x0b, 0x47, 0xe7, 0x39, 0x8e, 0xd0, 0x1a, 0x50, 0xa5, 0xb5, 0xb3,
	0x7a, 0x60, 0x0d, 0x5b, 0x98, 0x6a, 0x3d, 0xcf, 0x77, 0x5a, 0x15, 0xaf, 0xb9, 0x71, 0x56, 0x7f,
	0xb7, 0xf5, 0xe0, 0x68, 0xdb, 0x2a, 0xb8, 0x5e, 0xb5, 0x63, 0x39, 0xcb, 0x0f, 0xbd, 0x18, 0xa7,
	0x43, 0x07, 0x61, 0xaa, 0x1d, 0x67, 0x5b, 0x1a, 0xfc, 0x38, 0x9b, 0xce, 0x19, 0x9a, 0x80, 0xa3,
	0x58, 0x56, 0xd0, 0x55, 0x2e, 0x35, 0x09, 0xe0, 0x73, 0x23, 0x47, 0x53, 0x96, 0x15, 0xf7, 0xb9,
	0xd4, 0x6e, 0x2a, 0xb0, 0xde, 0x15, 0x37, 0xac, 0xd5, 0x50, 0x3b, 0x95, 0xc3, 0xcf, 0x8e, 0x82,
	0xb7, 0x57, 0xe3, 0x42, 0x6a, 0x32, 0xdc, 0xaf, 0xdd, 0x89, 0x5f, 0x35, 0xd0, 0x46, 0x07, 0x35,
	0xf0, 0x4e, 0xd1, 0x48, 0xf2, 0x42, 0xd3, 0xa5, 0x55, 0x8a, 0x35, 0x86, 0x8c, 0xdb, 0xdf, 0xe5,
	0xd8, 0xac, 0x45, 0xf8, 0x04, 0x21, 0x50, 0x60, 0x45, 0x32, 0x01, 0x21, 0x70, 0x04, 0xf6, 0x8b,
	0xfe, 0xa3, 0xe9, 0xf6, 0x14, 0xb4, 0x14, 0x8d, 0x66, 0x17, 0xd7, 0x08, 0xdd, 0x94, 0xa5, 0x62,
	0x65, 0x66, 0x84, 0xc2, 0x09, 0xf2, 0xd7, 0x02, 0x9e, 0x26, 0xfb, 0x57, 0x73, 0xf4, 0x2f, 0x39,
	0x78, 0x37, 0xfa, 0x73, 0xeb, 0xbf, 0xf4, 0x92, 0xab, 0x4c, 0xf2, 0xbc, 0x0f, 0xc7, 0x76, 0xf9,
	0x1d, 0x00, 0x00, 0xff, 0xff, 0x1b, 0xdd, 0x67, 0xe0, 0x79, 0x02, 0x00, 0x00,
}
