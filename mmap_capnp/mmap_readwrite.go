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

		// set the buf
		start := i * blockSize
		end := start + blockSize
		buf := bytes.NewBuffer(data[start:end])

		// write block
		if err := writeBlock(buf, i, blockSize); err != nil {
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
