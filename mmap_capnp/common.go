package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"zombiezen.com/go/capnproto2"
)

func createBlock(i int) (*TlogBlock, *capnp.Message, error) {
	// create block
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, nil, err
	}
	block, err := NewRootTlogBlock(seg)
	if err != nil {
		return nil, nil, err
	}
	setBlockVal(&block, i)
	return &block, msg, nil
}

// write single tlog block to a buffer
func writeBlock(buf *bytes.Buffer, i int) error {
	_, msg, err := createBlock(i)
	if err != nil {
		return err
	}
	// add it to mmap'ed file
	buf.Truncate(0)

	return capnp.NewEncoder(buf).Encode(msg)
}

func decodeBlock(buf *bytes.Buffer) (*TlogBlock, error) {
	msg, err := capnp.NewDecoder(buf).Decode()
	if err != nil {
		log.Printf("decode failed :%v\n", err)
		return nil, err
	}

	block, err := ReadRootTlogBlock(msg)
	if err != nil {
		log.Printf("decode failed to read root block:%v\n", err)
	}
	return &block, err
}

func createList(num int) (*capnp.Message, error) {
	// create the capnp aggregation object
	aggMsg, aggSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	agg, err := NewRootTlogAggregation(aggSeg)
	if err != nil {
		return nil, err
	}

	agg.SetName("the 1 M message")
	agg.SetSize(0)
	blockList, err := agg.NewBlocks(int32(num))
	if err != nil {
		return nil, err
	}

	// add blocks
	for i := 0; i < blockList.Len(); i++ {
		block := blockList.At(i)
		setBlockVal(&block, i)
	}

	agg.SetSize(uint64(num))
	return aggMsg, nil
}

// write tlog blocks to capnp list
func writeList(num int, w io.Writer) error {
	fmt.Println("create capnp messages aggregation...")
	aggMsg, err := createList(num)
	if err != nil {
		return err
	}

	// write capnp msg
	fmt.Println("write capnp messages to file/memory...")
	return capnp.NewEncoder(w).Encode(aggMsg)
}

// decode tlog aggregation
func decodeAggBlocks(r io.Reader) (*TlogAggregation, *TlogBlock_List, error) {
	msg, err := capnp.NewDecoder(r).Decode()
	if err != nil {
		return nil, nil, err
	}
	agg, err := ReadRootTlogAggregation(msg)
	if err != nil {
		return nil, nil, err
	}

	blocks, err := agg.Blocks()
	return &agg, &blocks, err
}
func setBlockVal(block *TlogBlock, val int) {
	block.SetSequence(uint64(val))

	if optDataLen > 0 {
		b := make([]byte, optDataLen) //yes, we make allocation here!
		block.SetText(string(b))
	}
}

func checkBlockVal(block *TlogBlock, val int) {
	if block.Sequence() != uint64(val) {
		log.Fatalf("invalid sequence. expected:%v, got:%v", val, block.Sequence())
	}
}

func dataLenInBlock() int {
	len := 8 /* sequence */ + 4 /* capnp overhead */
	len += optDataLen           /* text */
	return len
}
