// Code generated by capnpc-go.

package main

// AUTO GENERATED - DO NOT EDIT

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type TlogBlock struct{ capnp.Struct }

// TlogBlock_TypeID is the unique identifier for the type TlogBlock.
const TlogBlock_TypeID = 0x8cf178de3c82d431

func NewTlogBlock(s *capnp.Segment) (TlogBlock, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return TlogBlock{st}, err
}

func NewRootTlogBlock(s *capnp.Segment) (TlogBlock, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return TlogBlock{st}, err
}

func ReadRootTlogBlock(msg *capnp.Message) (TlogBlock, error) {
	root, err := msg.RootPtr()
	return TlogBlock{root.Struct()}, err
}

func (s TlogBlock) String() string {
	str, _ := text.Marshal(0x8cf178de3c82d431, s.Struct)
	return str
}

func (s TlogBlock) Sequence() uint64 {
	return s.Struct.Uint64(0)
}

func (s TlogBlock) SetSequence(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s TlogBlock) Text() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s TlogBlock) HasText() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TlogBlock) TextBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s TlogBlock) SetText(v string) error {
	return s.Struct.SetText(0, v)
}

// TlogBlock_List is a list of TlogBlock.
type TlogBlock_List struct{ capnp.List }

// NewTlogBlock creates a new list of TlogBlock.
func NewTlogBlock_List(s *capnp.Segment, sz int32) (TlogBlock_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return TlogBlock_List{l}, err
}

func (s TlogBlock_List) At(i int) TlogBlock { return TlogBlock{s.List.Struct(i)} }

func (s TlogBlock_List) Set(i int, v TlogBlock) error { return s.List.SetStruct(i, v.Struct) }

// TlogBlock_Promise is a wrapper for a TlogBlock promised by a client call.
type TlogBlock_Promise struct{ *capnp.Pipeline }

func (p TlogBlock_Promise) Struct() (TlogBlock, error) {
	s, err := p.Pipeline.Struct()
	return TlogBlock{s}, err
}

type TlogAggregation struct{ capnp.Struct }

// TlogAggregation_TypeID is the unique identifier for the type TlogAggregation.
const TlogAggregation_TypeID = 0xe46ab5b4b619e094

func NewTlogAggregation(s *capnp.Segment) (TlogAggregation, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return TlogAggregation{st}, err
}

func NewRootTlogAggregation(s *capnp.Segment) (TlogAggregation, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return TlogAggregation{st}, err
}

func ReadRootTlogAggregation(msg *capnp.Message) (TlogAggregation, error) {
	root, err := msg.RootPtr()
	return TlogAggregation{root.Struct()}, err
}

func (s TlogAggregation) String() string {
	str, _ := text.Marshal(0xe46ab5b4b619e094, s.Struct)
	return str
}

func (s TlogAggregation) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s TlogAggregation) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s TlogAggregation) SetName(v string) error {
	return s.Struct.SetText(0, v)
}

func (s TlogAggregation) Size() uint64 {
	return s.Struct.Uint64(0)
}

func (s TlogAggregation) SetSize(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s TlogAggregation) Timestamp() uint64 {
	return s.Struct.Uint64(8)
}

func (s TlogAggregation) SetTimestamp(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s TlogAggregation) VolumeId() uint32 {
	return s.Struct.Uint32(16)
}

func (s TlogAggregation) SetVolumeId(v uint32) {
	s.Struct.SetUint32(16, v)
}

func (s TlogAggregation) Blocks() (TlogBlock_List, error) {
	p, err := s.Struct.Ptr(1)
	return TlogBlock_List{List: p.List()}, err
}

func (s TlogAggregation) HasBlocks() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) SetBlocks(v TlogBlock_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewBlocks sets the blocks field to a newly
// allocated TlogBlock_List, preferring placement in s's segment.
func (s TlogAggregation) NewBlocks(n int32) (TlogBlock_List, error) {
	l, err := NewTlogBlock_List(s.Struct.Segment(), n)
	if err != nil {
		return TlogBlock_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s TlogAggregation) Prev() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return []byte(p.Data()), err
}

func (s TlogAggregation) HasPrev() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s TlogAggregation) SetPrev(v []byte) error {
	return s.Struct.SetData(2, v)
}

// TlogAggregation_List is a list of TlogAggregation.
type TlogAggregation_List struct{ capnp.List }

// NewTlogAggregation creates a new list of TlogAggregation.
func NewTlogAggregation_List(s *capnp.Segment, sz int32) (TlogAggregation_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3}, sz)
	return TlogAggregation_List{l}, err
}

func (s TlogAggregation_List) At(i int) TlogAggregation { return TlogAggregation{s.List.Struct(i)} }

func (s TlogAggregation_List) Set(i int, v TlogAggregation) error {
	return s.List.SetStruct(i, v.Struct)
}

// TlogAggregation_Promise is a wrapper for a TlogAggregation promised by a client call.
type TlogAggregation_Promise struct{ *capnp.Pipeline }

func (p TlogAggregation_Promise) Struct() (TlogAggregation, error) {
	s, err := p.Pipeline.Struct()
	return TlogAggregation{s}, err
}

const schema_f4533cbae6e08506 = "x\xdaT\x91\xbdk\x14Q\x14\xc5\xcf\xb9o&\xbb\x81" +
	"\xc4\xc9\x8b\x03\xc6\xff@E\x85%\x16\x12\x02~tZ" +
	"\xe5\xa9 V2\x8e\x8fq\xcc|\x99\x99\xc4\xc5R\xb0" +
	"\xb3\xdcB\x0ba\x05\x0b\x85m\x84\xd5\xc2\xcaF\xb0\xb0" +
	"\x12kYD\xec\x17l,d\xe4\x15\xbbk\xaa\x07?" +
	"\xee\xbb\xbfs\xb8k_/J\xcf?&\x80\xd9\xf0\x97" +
	"\xda\xde\xb7\xc7\xdb\xdf\xfb\xd3\xa70\xebd\xbb\xf4d\xf2" +
	"\xeb\xc3\xf6\xf5\xdf\xf0\xd9\x01\xf4\xf4\x87\xfe\xeb\xde?\x0f" +
	"\xc1v09\xfe~\xfc\xee\xfeO7\xa9\xfe\x9bT\x1d" +
	"`\xf3\x16\xd7y4u\x9f6-o\x12g\xda&+" +
	"\x93\xdbu|Ol\x1e\x9d\x8d\xa3\xaa\xa8\xb6nde" +
	"r9\xeb\x94\xf1\xee\x0ei\xba\xca\x03<\x02\xfa\xe4U" +
	"\xc0\x9cP4\xe7\x84dH\xc7z\xa7\x00sZ\xd1\x9c" +
	"\x17\xb6\xb5}\xb0o\x8b\xd8\x02\xe02\x84\xcb`\xd0\xd8" +
	"~\xc3\x15\x08W\xc0\xb9L\x1d\x96]J\x92=\x9bD" +
	"MZ\xb2p\xca\x8d\xb9\xf2\xb9[?P4\xc3\x85\xf2" +
	"\x85c\xcf\x14\xcd+\xa1\x16\x86\x14@\xbf\xbc\x06\x98\xa1" +
	"\xa2\x19\x09\xb5\xf2B*@\xbfq\x81_+\x9a\xb1P" +
	"{\x0c\xe9\x01\xfa\xed\x16`F\x8a\xe6\xb3P\xfb\x12\xd2" +
	"\x07\xf4'\xb7\xf3\xa3\xa2\xf9\"\x0c\x8a(\xb7\xb3\xc8A" +
	"\x9d>\xb2\xb32m\x93\xe6\xb6n\xa2\x1c\xac\xe6\xec\xa0" +
	"\xcc\xf6s{\xe5\xae+\xdd\x85\xb0\x0b^\xb8\x93\x95\xf1" +
	"n\xcd#\xe0\x8e\"\xd7\x16\xd7\x03\x1d\x0c\xaa={\xc0" +
	"U\x08W\xc1\x7f\x01\x00\x00\xff\xff7\x1fn\xfe"

func init() {
	schemas.Register(schema_f4533cbae6e08506,
		0x8cf178de3c82d431,
		0xe46ab5b4b619e094)
}
