// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: share/pb/share.proto

package pb

import (
	fmt "fmt"
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

type MerkleProof struct {
	Start    int64    `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	End      int64    `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
	Nodes    [][]byte `protobuf:"bytes,3,rep,name=nodes,proto3" json:"nodes,omitempty"`
	LeafHash []byte   `protobuf:"bytes,4,opt,name=leaf_hash,json=leafHash,proto3" json:"leaf_hash,omitempty"`
}

func (m *MerkleProof) Reset()         { *m = MerkleProof{} }
func (m *MerkleProof) String() string { return proto.CompactTextString(m) }
func (*MerkleProof) ProtoMessage()    {}
func (*MerkleProof) Descriptor() ([]byte, []int) {
	return fileDescriptor_349346f22540ed4a, []int{0}
}
func (m *MerkleProof) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerkleProof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerkleProof.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerkleProof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerkleProof.Merge(m, src)
}
func (m *MerkleProof) XXX_Size() int {
	return m.Size()
}
func (m *MerkleProof) XXX_DiscardUnknown() {
	xxx_messageInfo_MerkleProof.DiscardUnknown(m)
}

var xxx_messageInfo_MerkleProof proto.InternalMessageInfo

func (m *MerkleProof) GetStart() int64 {
	if m != nil {
		return m.Start
	}
	return 0
}

func (m *MerkleProof) GetEnd() int64 {
	if m != nil {
		return m.End
	}
	return 0
}

func (m *MerkleProof) GetNodes() [][]byte {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func (m *MerkleProof) GetLeafHash() []byte {
	if m != nil {
		return m.LeafHash
	}
	return nil
}

type Share struct {
	Data  []byte       `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Proof *MerkleProof `protobuf:"bytes,2,opt,name=Proof,proto3" json:"Proof,omitempty"`
}

func (m *Share) Reset()         { *m = Share{} }
func (m *Share) String() string { return proto.CompactTextString(m) }
func (*Share) ProtoMessage()    {}
func (*Share) Descriptor() ([]byte, []int) {
	return fileDescriptor_349346f22540ed4a, []int{1}
}
func (m *Share) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Share) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Share.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Share) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Share.Merge(m, src)
}
func (m *Share) XXX_Size() int {
	return m.Size()
}
func (m *Share) XXX_DiscardUnknown() {
	xxx_messageInfo_Share.DiscardUnknown(m)
}

var xxx_messageInfo_Share proto.InternalMessageInfo

func (m *Share) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Share) GetProof() *MerkleProof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func init() {
	proto.RegisterType((*MerkleProof)(nil), "ipld.pb.MerkleProof")
	proto.RegisterType((*Share)(nil), "ipld.pb.Share")
}

func init() { proto.RegisterFile("share/pb/share.proto", fileDescriptor_349346f22540ed4a) }

var fileDescriptor_349346f22540ed4a = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x4f, 0x4a, 0xc4, 0x30,
	0x1c, 0x85, 0x1b, 0x3b, 0xf5, 0x4f, 0xa6, 0x0b, 0x09, 0x5d, 0x14, 0x84, 0x50, 0x66, 0x55, 0x44,
	0x13, 0xd0, 0x1b, 0x88, 0xa8, 0x1b, 0x41, 0xe2, 0xce, 0x8d, 0x24, 0xd3, 0xcc, 0xa4, 0x18, 0x27,
	0x21, 0x89, 0xf7, 0xf0, 0x58, 0x2e, 0xbb, 0x74, 0x29, 0xed, 0x45, 0x24, 0xa9, 0xc8, 0xec, 0xbe,
	0xdf, 0xe3, 0x91, 0xf7, 0x11, 0x58, 0x79, 0xc5, 0x9d, 0xa4, 0x56, 0xd0, 0x04, 0xc4, 0x3a, 0x13,
	0x0c, 0x3a, 0xea, 0xad, 0xee, 0x88, 0x15, 0x2b, 0x05, 0x97, 0x8f, 0xd2, 0xbd, 0x69, 0xf9, 0xe4,
	0x8c, 0xd9, 0xa0, 0x0a, 0x16, 0x3e, 0x70, 0x17, 0x6a, 0xd0, 0x80, 0x36, 0x67, 0xf3, 0x81, 0x4e,
	0x61, 0x2e, 0x77, 0x5d, 0x7d, 0x90, 0xb2, 0x88, 0xb1, 0xb7, 0x33, 0x9d, 0xf4, 0x75, 0xde, 0xe4,
	0x6d, 0xc9, 0xe6, 0x03, 0x9d, 0xc1, 0x13, 0x2d, 0xf9, 0xe6, 0x55, 0x71, 0xaf, 0xea, 0x45, 0x03,
	0xda, 0x92, 0x1d, 0xc7, 0xe0, 0x81, 0x7b, 0xb5, 0xba, 0x87, 0xc5, 0x73, 0x34, 0x40, 0x08, 0x2e,
	0x6e, 0x79, 0xe0, 0x69, 0xa2, 0x64, 0x89, 0xd1, 0x39, 0x2c, 0x92, 0x40, 0xda, 0x58, 0x5e, 0x55,
	0xe4, 0xcf, 0x8f, 0xec, 0xc9, 0xb1, 0xb9, 0x72, 0x73, 0xf7, 0x35, 0x62, 0x30, 0x8c, 0x18, 0xfc,
	0x8c, 0x18, 0x7c, 0x4e, 0x38, 0x1b, 0x26, 0x9c, 0x7d, 0x4f, 0x38, 0x7b, 0xb9, 0xd8, 0xf6, 0x41,
	0x7d, 0x08, 0xb2, 0x36, 0xef, 0x74, 0x2d, 0xb5, 0xf4, 0xa1, 0xe7, 0xc6, 0x6d, 0xff, 0xf9, 0x32,
	0x6a, 0xd2, 0xf8, 0x34, 0xb5, 0x42, 0x1c, 0xa6, 0xaf, 0xb8, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff,
	0xee, 0x51, 0xeb, 0x90, 0x22, 0x01, 0x00, 0x00,
}

func (m *MerkleProof) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerkleProof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerkleProof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.LeafHash) > 0 {
		i -= len(m.LeafHash)
		copy(dAtA[i:], m.LeafHash)
		i = encodeVarintShare(dAtA, i, uint64(len(m.LeafHash)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Nodes) > 0 {
		for iNdEx := len(m.Nodes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Nodes[iNdEx])
			copy(dAtA[i:], m.Nodes[iNdEx])
			i = encodeVarintShare(dAtA, i, uint64(len(m.Nodes[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.End != 0 {
		i = encodeVarintShare(dAtA, i, uint64(m.End))
		i--
		dAtA[i] = 0x10
	}
	if m.Start != 0 {
		i = encodeVarintShare(dAtA, i, uint64(m.Start))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Share) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Share) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Share) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Proof != nil {
		{
			size, err := m.Proof.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintShare(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintShare(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintShare(dAtA []byte, offset int, v uint64) int {
	offset -= sovShare(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MerkleProof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Start != 0 {
		n += 1 + sovShare(uint64(m.Start))
	}
	if m.End != 0 {
		n += 1 + sovShare(uint64(m.End))
	}
	if len(m.Nodes) > 0 {
		for _, b := range m.Nodes {
			l = len(b)
			n += 1 + l + sovShare(uint64(l))
		}
	}
	l = len(m.LeafHash)
	if l > 0 {
		n += 1 + l + sovShare(uint64(l))
	}
	return n
}

func (m *Share) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovShare(uint64(l))
	}
	if m.Proof != nil {
		l = m.Proof.Size()
		n += 1 + l + sovShare(uint64(l))
	}
	return n
}

func sovShare(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozShare(x uint64) (n int) {
	return sovShare(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MerkleProof) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: MerkleProof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerkleProof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Start", wireType)
			}
			m.Start = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Start |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field End", wireType)
			}
			m.End = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.End |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nodes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Nodes = append(m.Nodes, make([]byte, postIndex-iNdEx))
			copy(m.Nodes[len(m.Nodes)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeafHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LeafHash = append(m.LeafHash[:0], dAtA[iNdEx:postIndex]...)
			if m.LeafHash == nil {
				m.LeafHash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func (m *Share) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: Share: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Share: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proof", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Proof == nil {
				m.Proof = &MerkleProof{}
			}
			if err := m.Proof.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func skipShare(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowShare
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
					return 0, ErrIntOverflowShare
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
					return 0, ErrIntOverflowShare
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
				return 0, ErrInvalidLengthShare
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupShare
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthShare
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthShare        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowShare          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupShare = fmt.Errorf("proto: unexpected end of group")
)
