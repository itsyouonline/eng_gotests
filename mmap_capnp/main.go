package main

import (
	"bytes"
	"log"
	"os"
	"syscall"

	"zombiezen.com/go/capnproto2"
)

func main() {
	num := 1000 * 1000
	if err := writeBulkRead(num); err != nil {
		log.Printf("err = %v\n", err)
	}

	if err := writeOneReadOne(num, 70); err != nil {
		log.Printf("writeOneReadOn err = %v\n", err)
	}
}

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

// - encode 1M capnp message & write it to mem mapped file
// - read some data
func writeBulkRead(num int) error {
	log.Printf("======== create %v capnp messages and write it to mem mapped file ==========\n", num)

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

	size := countMemSize(num)

	// create mem mapped file
	f, data, err := createMemMap(size)
	if err != nil {
		return err
	}
	defer f.Close()
	defer syscall.Munmap(data)
	buf := bytes.NewBuffer(data)
	buf.Truncate(0)

	// write capnp msg to mem mapped file
	if err := capnp.NewEncoder(buf).Encode(aggMsg); err != nil {
		return err
	}

	// read it again
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

// create mmap'ed file with given size
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

func countMemSize(num int) int {
	return num * 40
}
