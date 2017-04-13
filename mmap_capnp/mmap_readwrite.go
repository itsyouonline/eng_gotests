package main

import (
	"bytes"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

//  num times:
// 	- encode canpnp message
//  - add it to mmap'ed file
//  - read some written capnp message
// all messages are written to the same file
// every messages is appended, we don't rewrite the whole file
// IMPORTANT: this only works because our messages have a fixed size, so we canpnp
// calculate the start and endpoint of every message
func writeOneReadOne(num int) error {
	log.Info("====== {write one capnp doc, read one capnp doc} =========")
	log.Infof("%v messages", num)

	blockSize := dataLenInBlock() + 30 /* 30 extra space needed by capnp. TODO : find a way to count it automatically*/

	// create mem mapped file
	f, data, err := createMemMap(num * blockSize)
	if err != nil {
		return err
	}
	// defers execute in reverse order so Munmap is called first, then the file gets closed
	defer f.Close()
	defer syscall.Munmap(data)

	for i := 0; i < num; i++ {

		// set the buf
		// calculate the offset of the message start and endpoint in the file, then
		// create a new slice of the target region
		start := i * blockSize
		end := start + blockSize
		buf := bytes.NewBuffer(data[start:end])

		// write block to the slice
		if err := writeBlock(buf, i); err != nil {
			return err
		}

		if i < 5 { // too early to read
			continue
		}

		// decode some block
		// decide which block we want to read, then calculate its start and end offset
		// create a new slice of the target region. this slice holds the contents of our
		// block
		readIdx := i / 2
		start = readIdx * blockSize
		end = start + blockSize
		buf = bytes.NewBuffer(data[start:end])

		// decode the block
		decodedBlock, err := decodeBlock(buf)
		if err != nil {
			log.Infof("failed to decode block:%v", err)
			return err
		}

		// check it's values. checkBlockVal logs a fatal error and causes the program
		// to exit if the expected value is not found in the block
		checkBlockVal(decodedBlock, readIdx)
	}
	log.Println("all blocks written, all checks on written blocks succeeded")
	return nil
}

// - encode num capnp message to capnp list & write it to mem mapped file
// - read some data
// generate num messages. messages are stored in a capnp list of another struct,
// which is then encoded and written to the file as a whole. this woul allow working
// with different sized blocks, however changing any one block requires an entire
// file rewrite
func writeListRead(num int) error {
	log.Infof("======== create %v capnp messages to store in capnp list and write it to mem mapped file ==========", num)

	// calc the size of our mmaped file
	size := 100 + (num * dataLenInBlock())
	// create mem mapped file
	f, data, err := createMemMap(size)
	if err != nil {
		return err
	}
	// defer the file close first then Munmap, so mumpa gets called first, data
	// is written tot the file, and then the file is closed
	defer f.Close()
	defer syscall.Munmap(data)
	buf := bytes.NewBuffer(data)

	// create & write capnp msg to mem mapped file
	buf.Truncate(0)
	if err := writeList(num, buf); err != nil {
		return err
	}

	log.Info("all blocks written to the list and encoded.")

	// read it again to verify the content
	_, decodedBlockList, err := decodeAggBlocks(buf)
	if err != nil {
		log.Infof("failed to decode block list:%v", err)
		return err
	}

	for i := 0; i < num; i++ {
		block := decodedBlockList.At(i)
		// do some checking
		checkBlockVal(&block, i)
	}

	log.Info("list decoded again and all blocks successfully verified.")

	return nil
}
