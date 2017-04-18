package main

import (
	"bytes"
	"io"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"

	"zombiezen.com/go/capnproto2"
)

// createMemMap creates a file with the given size. it then creates an mmap of said
// file. both the mmaped file and a pointer to the disk file are returned. it is the
// callers responsibility to unmap and close the file when they are done
func createMemMap(size int) (*os.File, []byte, error) {
	// create mem mapped file
	f, err := os.Create("/tmp/capnp_mmap")
	if err != nil {
		return nil, nil, err
	}

	if _, err := f.Seek(int64(size-1), 0); err != nil {
		f.Close()
		return nil, nil, err
	}
	_, err = f.Write([]byte(" "))
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return f, data, nil
}

// createBlock creates a new tlog block and sets the sequence
func createBlock(i int) (*TlogBlock, *capnp.Message, error) {
	segmentBuf := make([]byte, 0, segmentBufferSize())
	// create block
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
	setBlockVal(&block, i)
	return &block, msg, nil
}

// writeBlock writes a single tlog block to a buffer. any existing data in the
// buffer is removed
func writeBlock(buf *bytes.Buffer, i int) error {
	_, msg, err := createBlock(i)
	if err != nil {
		return err
	}
	// clear the buffer
	buf.Truncate(0)
	// encode and write the block in the buffer
	return capnp.NewEncoder(buf).Encode(msg)
}

// decodeBlock reads and decodes a single tlog block from a buffer. the buffer
// must not contain any excessive data
func decodeBlock(buf *bytes.Buffer) (*TlogBlock, error) {
	// read the buffer and decode the message
	msg, err := capnp.NewDecoder(buf).Decode()
	if err != nil {
		log.Infof("decode failed: %v", err)
		return nil, err
	}

	// get the block from the message
	block, err := ReadRootTlogBlock(msg)
	if err != nil {
		log.Warnf("decode failed to read root block: %v", err)
	}
	return &block, err
}

// createList creates the defined amount of TlogBlocks, stored in a capnp list
// wrapped in a TlogAggregation.
func createList(num int) (*capnp.Message, error) {
	// create a new capnp message
	aggMsg, aggSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	// create the capnp aggregation object
	agg, err := NewRootTlogAggregation(aggSeg)
	if err != nil {
		return nil, err
	}

	// set some data fields of the aggregation, this doesn't have a real function but
	// allows for easy validation
	agg.SetName("the 1M messages")
	agg.SetSize(0)
	// create the TlogBlock list
	blockList, err := agg.NewBlocks(int32(num))
	if err != nil {
		return nil, err
	}

	// add blocks to the new list
	for i := 0; i < blockList.Len(); i++ {
		block := blockList.At(i)
		setBlockVal(&block, i)
	}

	// update the size of the aggregation object now that we've added the blocks
	agg.SetSize(uint64(num))
	return aggMsg, nil
}

// writeList creates a new TlogAggregation, writes the designated amount of tlogBlockSize
// to the aggregations list, encodes it, and writes it to the provided writer.
func writeList(num int, w io.Writer) error {
	log.Info("create capnp messages aggregation...")
	aggMsg, err := createList(num)
	if err != nil {
		return err
	}

	// write capnp msg
	log.Info("write capnp messages to file/memory...")
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

// setBlockVal sets the sequence of a TlogBlock and sets the text field if the data
// len has been specified to be bigger than 0 for this run
func setBlockVal(block *TlogBlock, val int) {
	block.SetSequence(uint64(val))

	if optDataLen > 0 {
		b := make([]byte, optDataLen) //yes, we make allocation here!
		block.SetText(string(b))
	}
}

// checkBlockVal verifies that the sequence set in the block . if it is not correct,
// the program prints the expected and actual value before exiting
func checkBlockVal(block *TlogBlock, val int) {
	if block.Sequence() != uint64(val) {
		log.Fatalf("invalid sequence. expected: %v, got: %v", val, block.Sequence())
	}
}

// dataLenInBlock returns the length of the data in a unenoded capnp tlog block
func dataLenInBlock() int {
	len := 8 /* sequence */ + 4 /* capnp overhead */
	len += optDataLen           /* text */
	return len
}

func segmentBufferSize() int {
	return 8*((dataLenInBlock()/8)+1) + 10
}
