// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package protobufs

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
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x4d, 0x4f, 0xdb, 0x30,
	0x18, 0x80, 0x97, 0xb5, 0x6b, 0x1b, 0xa7, 0x1f, 0x9b, 0x4f, 0xd6, 0xa6, 0x49, 0x59, 0x34, 0x4d,
	0x39, 0xe5, 0x30, 0x8e, 0x70, 0x81, 0x9e, 0x10, 0x54, 0x42, 0x51, 0x11, 0x12, 0x17, 0x2b, 0x6e,
	0x4c, 0xb0, 0x9a, 0xc4, 0xc6, 0x1f, 0x02, 0x7e, 0x02, 0xff, 0x1a, 0xf9, 0x4d, 0x5b, 0xda, 0x72,
	0x4b, 0x1e, 0x3d, 0x79, 0xf2, 0xda, 0x7a, 0x51, 0x58, 0x28, 0x91, 0x29, 0x2d, 0xad, 0x4c, 0xde,
	0xfa, 0x68, 0x7a, 0xe3, 0x58, 0x2d, 0xcc, 0x63, 0xce, 0x9f, 0x1c, 0x37, 0x16, 0xff, 0x42, 0xa1,
	0x6b, 0xc5, 0x0b, 0xb5, 0xa2, 0xe1, 0x24, 0x88, 0x83, 0xb4, 0x97, 0x8f, 0x3c, 0x58, 0x8a, 0x86,
	0x63, 0x82, 0x86, 0xda, 0xb5, 0xad, 0x68, 0x2b, 0xf2, 0x35, 0x0e, 0xd2, 0x51, 0xbe, 0x7d, 0xc5,
	0xff, 0xd0, 0xac, 0x11, 0x75, 0x2d, 0xe8, 0x4a, 0x39, 0xea, 0x4c, 0x51, 0x71, 0xd2, 0x8b, 0x83,
	0xb4, 0x9f, 0x4f, 0x00, 0xcf, 0x95, 0xbb, 0xf5, 0x10, 0x3c, 0xde, 0x48, 0xfd, 0x4a, 0xad, 0xb4,
	0x45, 0x4d, 0x1b, 0x46, 0xfa, 0x1b, 0x0f, 0xf0, 0xd2, 0xd3, 0x05, 0xc3, 0x09, 0xda, 0x00, 0xaa,
	0x8d, 0xf1, 0xd6, 0x37, 0xb0, 0xa2, 0x0e, 0xe6, 0xc6, 0x2c, 0xd8, 0x5e, 0xab, 0x16, 0x8d, 0xb0,
	0xde, 0x1a, 0xec, 0xb7, 0xae, 0x3d, 0xed, 0x5a, 0xa5, 0x30, 0xeb, 0x6e, 0x2c, 0x6f, 0x0d, 0xe3,
	0x20, 0xc5, 0x79, 0xe4, 0x21, 0x4c, 0xb5, 0xe7, 0xec, 0x4a, 0xa3, 0x0f, 0x67, 0xdb, 0xf9, 0x8b,
	0xa6, 0xe0, 0x68, 0x5e, 0x94, 0x74, 0xcd, 0x94, 0x21, 0x21, 0xfc, 0x6e, 0xec, 0x69, 0xce, 0x8b,
	0xf2, 0x8a, 0x29, 0xe3, 0xa7, 0x02, 0xeb, 0x59, 0x0b, 0xcb, 0x3b, 0x0d, 0x75, 0x53, 0x79, 0x7c,
	0xe7, 0x29, 0x78, 0x07, 0x35, 0x21, 0x95, 0x21, 0xd1, 0x61, 0xed, 0x52, 0x7e, 0xaa, 0x81, 0x36,
	0x3e, 0xaa, 0x81, 0xf7, 0x07, 0x8d, 0x95, 0x28, 0x0d, 0x5d, 0x39, 0xad, 0x79, 0x6b, 0xc9, 0xa4,
	0xbb, 0x2e, 0xcf, 0xe6, 0x1d, 0xc2, 0xbf, 0x11, 0x02, 0x05, 0x8e, 0x48, 0xa6, 0x20, 0x84, 0x9e,
	0xc0, 0xf9, 0x92, 0x1f, 0x68, 0xb6, 0x5b, 0x05, 0xa3, 0x64, 0x6b, 0xf8, 0xff, 0x33, 0x84, 0xce,
	0xab, 0x4a, 0xf3, 0xaa, 0xb0, 0x52, 0xe3, 0x0c, 0x0d, 0x37, 0x02, 0x9e, 0x65, 0x87, 0x5b, 0xf3,
	0xf3, 0x7b, 0x76, 0xf4, 0x6d, 0xf2, 0xe5, 0x62, 0x72, 0x1f, 0x65, 0xa7, 0xb0, 0x67, 0xcc, 0x3d,
	0x18, 0x36, 0x80, 0xc7, 0x93, 0xf7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x23, 0x67, 0xff, 0xdf, 0x7f,
	0x02, 0x00, 0x00,
}
