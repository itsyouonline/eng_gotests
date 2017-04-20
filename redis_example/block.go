package main

import (
	capnp "zombiezen.com/go/capnproto2"
)

// generateSequenceBlock creates a new Capnp message with tlogblock.
func generateSequenceBlock(seq int) (*TlogBlock, *capnp.Message, error) {
	segmentBuf := make([]byte, 0, segmentBufferSize())
	// create a message
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(segmentBuf))
	if err != nil {
		return nil, nil, err
	}
	// create the block
	block, err := NewRootTlogBlock(seg)
	if err != nil {
		return nil, nil, err
	}
	//set the block value
	setBlockVal(&block, seq)
	return &block, msg, nil
}

func setBlockVal(block *TlogBlock, val int) {
	block.SetSequence(uint64(val))

	if dataSize > 0 {
		b := make([]byte, dataSize)
		block.SetText(string(b))
	}
}

func segmentBufferSize() int {
	return 8*((dataLenInBlock()/8)+1) + 20
}

func dataLenInBlock() int {
	len := 16
	len += dataSize
	return len
}
