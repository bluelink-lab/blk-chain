// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: evm/types.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Whitelist struct {
	Hashes []string `protobuf:"bytes,1,rep,name=hashes,proto3" json:"hashes,omitempty" yaml:"hashes"`
}

func (m *Whitelist) Reset()         { *m = Whitelist{} }
func (m *Whitelist) String() string { return proto.CompactTextString(m) }
func (*Whitelist) ProtoMessage()    {}
func (*Whitelist) Descriptor() ([]byte, []int) {
	return fileDescriptor_6eba926c274d8fd0, []int{0}
}
func (m *Whitelist) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Whitelist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Whitelist.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Whitelist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Whitelist.Merge(m, src)
}
func (m *Whitelist) XXX_Size() int {
	return m.Size()
}
func (m *Whitelist) XXX_DiscardUnknown() {
	xxx_messageInfo_Whitelist.DiscardUnknown(m)
}

var xxx_messageInfo_Whitelist proto.InternalMessageInfo

func (m *Whitelist) GetHashes() []string {
	if m != nil {
		return m.Hashes
	}
	return nil
}

type DeferredInfo struct {
	TxIndex uint32                                 `protobuf:"varint,1,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	TxHash  []byte                                 `protobuf:"bytes,2,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	TxBloom []byte                                 `protobuf:"bytes,3,opt,name=tx_bloom,json=txBloom,proto3" json:"tx_bloom,omitempty"`
	Surplus github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=surplus,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"surplus"`
	Error   string                                 `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
}

func (m *DeferredInfo) Reset()         { *m = DeferredInfo{} }
func (m *DeferredInfo) String() string { return proto.CompactTextString(m) }
func (*DeferredInfo) ProtoMessage()    {}
func (*DeferredInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6eba926c274d8fd0, []int{1}
}
func (m *DeferredInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeferredInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeferredInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeferredInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeferredInfo.Merge(m, src)
}
func (m *DeferredInfo) XXX_Size() int {
	return m.Size()
}
func (m *DeferredInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_DeferredInfo.DiscardUnknown(m)
}

var xxx_messageInfo_DeferredInfo proto.InternalMessageInfo

func (m *DeferredInfo) GetTxIndex() uint32 {
	if m != nil {
		return m.TxIndex
	}
	return 0
}

func (m *DeferredInfo) GetTxHash() []byte {
	if m != nil {
		return m.TxHash
	}
	return nil
}

func (m *DeferredInfo) GetTxBloom() []byte {
	if m != nil {
		return m.TxBloom
	}
	return nil
}

func (m *DeferredInfo) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*Whitelist)(nil), "sheprotocol.blk-chain.evm.Whitelist")
	proto.RegisterType((*DeferredInfo)(nil), "sheprotocol.blk-chain.evm.DeferredInfo")
}

func init() { proto.RegisterFile("evm/types.proto", fileDescriptor_6eba926c274d8fd0) }

var fileDescriptor_6eba926c274d8fd0 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x31, 0x6b, 0xf2, 0x40,
	0x18, 0xc7, 0x73, 0xaf, 0xaf, 0xa6, 0x1e, 0x4a, 0x69, 0x10, 0x7a, 0xed, 0x10, 0x43, 0x86, 0x92,
	0x0e, 0x26, 0x43, 0xa1, 0x43, 0xc7, 0x50, 0xa8, 0xae, 0x59, 0x0a, 0x5d, 0x24, 0xc6, 0x47, 0x73,
	0x34, 0xc9, 0xc9, 0x3d, 0xa7, 0x9c, 0xdf, 0xa2, 0x5f, 0xa8, 0xbb, 0xa3, 0x63, 0xe9, 0x20, 0x45,
	0xbf, 0x41, 0x3f, 0x41, 0x49, 0x8c, 0xa5, 0xd3, 0x3d, 0xbf, 0xfb, 0x3f, 0xff, 0x3b, 0xf8, 0xd1,
	0x73, 0x58, 0xe5, 0x81, 0x5a, 0x2f, 0x00, 0xfd, 0x85, 0x14, 0x4a, 0x58, 0x0c, 0x81, 0x57, 0x53,
	0x22, 0x32, 0x1f, 0x81, 0x27, 0x69, 0xcc, 0x0b, 0x1f, 0x56, 0xf9, 0x75, 0x6f, 0x2e, 0xe6, 0xa2,
	0x8a, 0x82, 0x72, 0x3a, 0xee, 0xbb, 0xf7, 0xb4, 0xfd, 0x9c, 0x72, 0x05, 0x19, 0x47, 0x65, 0xdd,
	0xd2, 0x56, 0x1a, 0x63, 0x0a, 0xc8, 0x88, 0xd3, 0xf0, 0xda, 0xe1, 0xc5, 0xf7, 0xae, 0xdf, 0x5d,
	0xc7, 0x79, 0xf6, 0xe0, 0x1e, 0xef, 0xdd, 0xa8, 0x5e, 0x70, 0xdf, 0x09, 0xed, 0x3c, 0xc2, 0x0c,
	0xa4, 0x84, 0xe9, 0xa8, 0x98, 0x09, 0xeb, 0x8a, 0x9e, 0x29, 0x3d, 0xe6, 0xc5, 0x14, 0x34, 0x23,
	0x0e, 0xf1, 0xba, 0x91, 0xa9, 0xf4, 0xa8, 0x44, 0xeb, 0x92, 0x9a, 0x4a, 0x8f, 0xcb, 0x22, 0xfb,
	0xe7, 0x10, 0xaf, 0x13, 0xb5, 0x94, 0x1e, 0xc6, 0x98, 0xd6, 0x9d, 0x49, 0x26, 0x44, 0xce, 0x1a,
	0x55, 0x62, 0x2a, 0x1d, 0x96, 0x68, 0x0d, 0xa9, 0x89, 0x4b, 0xb9, 0xc8, 0x96, 0xc8, 0xfe, 0x3b,
	0xc4, 0x6b, 0x87, 0xfe, 0x66, 0xd7, 0x37, 0x3e, 0x77, 0xfd, 0x9b, 0x39, 0x57, 0xe9, 0x72, 0xe2,
	0x27, 0x22, 0x0f, 0x12, 0x81, 0xb9, 0xc0, 0xfa, 0x18, 0xe0, 0xf4, 0xb5, 0x56, 0x31, 0x2a, 0x54,
	0x74, 0xaa, 0x5b, 0x3d, 0xda, 0x04, 0x29, 0x85, 0x64, 0xcd, 0xf2, 0x9d, 0xe8, 0x08, 0xe1, 0xd3,
	0x66, 0x6f, 0x93, 0xed, 0xde, 0x26, 0x5f, 0x7b, 0x9b, 0xbc, 0x1d, 0x6c, 0x63, 0x7b, 0xb0, 0x8d,
	0x8f, 0x83, 0x6d, 0xbc, 0x0c, 0xfe, 0x7c, 0x80, 0xc0, 0x07, 0x27, 0x9b, 0x15, 0x54, 0x3a, 0x03,
	0x1d, 0xfc, 0x6a, 0x9f, 0xb4, 0xaa, 0xfc, 0xee, 0x27, 0x00, 0x00, 0xff, 0xff, 0xc8, 0xbe, 0x17,
	0x45, 0x8a, 0x01, 0x00, 0x00,
}

func (m *Whitelist) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Whitelist) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Whitelist) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hashes) > 0 {
		for iNdEx := len(m.Hashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Hashes[iNdEx])
			copy(dAtA[i:], m.Hashes[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Hashes[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *DeferredInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeferredInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeferredInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Error) > 0 {
		i -= len(m.Error)
		copy(dAtA[i:], m.Error)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Error)))
		i--
		dAtA[i] = 0x2a
	}
	{
		size := m.Surplus.Size()
		i -= size
		if _, err := m.Surplus.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.TxBloom) > 0 {
		i -= len(m.TxBloom)
		copy(dAtA[i:], m.TxBloom)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.TxBloom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TxHash) > 0 {
		i -= len(m.TxHash)
		copy(dAtA[i:], m.TxHash)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.TxHash)))
		i--
		dAtA[i] = 0x12
	}
	if m.TxIndex != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.TxIndex))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Whitelist) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Hashes) > 0 {
		for _, s := range m.Hashes {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *DeferredInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TxIndex != 0 {
		n += 1 + sovTypes(uint64(m.TxIndex))
	}
	l = len(m.TxHash)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.TxBloom)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Surplus.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.Error)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Whitelist) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Whitelist: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Whitelist: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hashes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hashes = append(m.Hashes, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *DeferredInfo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DeferredInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeferredInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxIndex", wireType)
			}
			m.TxIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TxIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHash = append(m.TxHash[:0], dAtA[iNdEx:postIndex]...)
			if m.TxHash == nil {
				m.TxHash = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxBloom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxBloom = append(m.TxBloom[:0], dAtA[iNdEx:postIndex]...)
			if m.TxBloom == nil {
				m.TxBloom = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Surplus", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Surplus.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Error", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Error = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
