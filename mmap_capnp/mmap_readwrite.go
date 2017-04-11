package main

import (
	"bytes"
	"log"
	"syscall"

	"zombiezen.com/go/capnproto2"
)

// for 0..num
// 	- encode canpnp message
//  - add it to mmap'ed file
//  - read some written capnp message
func writeOneReadOne(num, blockSize int) error {
	log.Printf("====== for i...%v {write one capnp doc, read one capnp doc} ========= \n", num)
	// create mem mapped file
	f, data, err := createMemMap(num * blockSize)
	if err != nil {
		return err
	}
	defer f.Close()
	defer syscall.Munmap(data)

	for i := 0; i < num; i++ {
		// write block
		if err := writeBlock(data, i, blockSize); err != nil {
			return err
		}

		if i < 5 { // too early to read
			continue
		}

		// decode some block
		readIdx := i / 2

		decodedBlock, err := decodeBlock(data, readIdx, blockSize)
		if err != nil {
			log.Printf("failed to decode block:%v\n", err)
			return err
		}

		// check it's values
		checkBlockVal(decodedBlock, readIdx)
	}
	log.Println("all  good")
	return nil
}

func decodeBlock(data []byte, i, blockSize int) (*TlogBlock, error) {
	start := i * blockSize
	end := start + blockSize
	buf := bytes.NewBuffer(data[start:end])

	msg, err := capnp.NewDecoder(buf).Decode()
	if err != nil {
		return nil, err
	}

	block, err := ReadRootTlogBlock(msg)
	return &block, err
}
func writeBlock(data []byte, i, blockSize int) error {
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
	start := i * blockSize
	end := start + blockSize
	buf := bytes.NewBuffer(data[start:end])
	buf.Truncate(0)

	return capnp.NewEncoder(buf).Encode(msg)
}

// - encode 1M capnp message to capnp list & write it to mem mapped file

// - read some data
func writeListRead(num int) error {
	log.Printf("======== create %v capnp messages to capnp list and write it to mem mapped file ==========\n", num)

	// create mem mapped file
	size := countMemSize(num)

	f, data, err := createMemMap(size)
	if err != nil {
		return err
	}
	defer f.Close()
	defer syscall.Munmap(data)
	buf := bytes.NewBuffer(data)

	// create & write capnp msg to mem mapped file
	if err := writeList(num, buf); err != nil {
		return err
	}

	// read it again to verify the content
	_, decodedBlockList, err := decodeAggBlocks(buf)
	if err != nil {
		log.Printf("failed to decode block list:%v", err)
		return err
	}

	for i := 0; i < num; i++ {
		block := decodedBlockList.At(i)
		// do some checking
		checkBlockVal(&block, i)
	}

	log.Println("all  good")

	return nil
}

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

func decodeAggBlocks(buf *bytes.Buffer) (*TlogAggregation, *TlogBlock_List, error) {
	msg, err := capnp.NewDecoder(buf).Decode()
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
