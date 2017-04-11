package main

import (
	"bytes"
	"log"

	"zombiezen.com/go/capnproto2"
)

// write single tlog block to a buffer
func writeBlock(buf *bytes.Buffer, i, blockSize int) error {
	// create block
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}
	block, err := NewRootTlogBlock(seg)
	if err != nil {
		return err
	}
	setBlockVal(&block, i)

	// add it to mmap'ed file
	buf.Truncate(0)

	return capnp.NewEncoder(buf).Encode(msg)
}

// write tlog blocks to capnp list
func writeList(num int, buf *bytes.Buffer) error {
	buf.Truncate(0)

	// create the capnp aggregation object
	aggMsg, aggSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return err
	}
	agg, err := NewRootTlogAggregation(aggSeg)
	if err != nil {
		return err
	}

	agg.SetName("the 1 M message")
	agg.SetSize(0)
	blockList, err := agg.NewBlocks(int32(num))
	if err != nil {
		return err
	}

	// add blocks
	for i := 0; i < blockList.Len(); i++ {
		block := blockList.At(i)
		setBlockVal(&block, i)
	}

	agg.SetSize(uint64(num))

	// write capnp msg to mem mapped file
	return capnp.NewEncoder(buf).Encode(aggMsg)
}

func setBlockVal(block *TlogBlock, val int) {
	block.SetVolumeId(uint32(val))
	block.SetSequence(uint64(val))
	block.SetLba(uint64(val))
	block.SetTimestamp(uint64(val))
}

func checkBlockVal(block *TlogBlock, val int) {
	if block.VolumeId() != uint32(val) {
		log.Fatalf("invalid volume id. expected:%v, got:%v", val, block.VolumeId)
	}
	if block.Sequence() != uint64(val) {
		log.Fatalf("invalid sequence. expected:%v, got:%v", val, block.Sequence())
	}
	if block.Timestamp() != uint64(val) {
		log.Fatalf("invalid timestamp. expected:%v, got:%v", val, block.Timestamp())
	}

}
