package main

import (
	"bytes"
	"log"
	"os"
	"syscall"

	"zombiezen.com/go/capnproto2"
)

func main() {
	if err := writeBulkRead(1000 * 1000); err != nil {
		log.Printf("err = %v\n", err)
	}
}

// - encode 1M capnp message & write it to mem mapped file
// - read some data
func writeBulkRead(num int) error {
	log.Println("create 1M capnp messages and write it to mem mapped file")

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
		block.SetVolumeId(1)
		block.SetSequence(uint64(i))
		block.SetLba(uint64(i))
		block.SetTimestamp(uint64(i))

	}
	agg.SetSize(uint64(num))

	size := num * 40 // TODO : how we estimate the size of generated capnp message?

	// create mem mapped file
	f, err := os.Create("/tmp/capnp_mmap")
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(int64(size-1), 0); err != nil {
		return err
	}
	_, err = f.Write([]byte(" "))
	if err != nil {
		return err
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return err
	}
	defer syscall.Munmap(data)

	buf := bytes.NewBuffer(data)

	buf.Write([]byte("hello, just a test"))
	log.Printf("write to mem mapped file seems work")
	buf.Truncate(0)

	// write capnp msg to mem mapped file
	if err := capnp.NewEncoder(buf).Encode(aggMsg); err != nil {
		return err
	}
	log.Println("writing capnp message seems work")

	// read it again
	msg, err := capnp.NewDecoder(buf).Decode()
	decodedAgg, err := ReadRootTlogAggregation(msg)
	if err != nil {
		log.Printf("failed to decode message:%v", err)
		return err
	}
	name, err := decodedAgg.Name()
	if err != nil {
		log.Printf("failed to get agg name:%v\n", err)
	}

	log.Printf("agg name = %v\n", name)

	decodedBlockList, err := decodedAgg.Blocks()
	if err != nil {
		log.Printf("failed to decode block list:%v", err)
		return err
	}

	for i := 0; i < num; i++ {
		block := decodedBlockList.At(i)
		// do some checking
		if block.Sequence() != uint64(i) {
			log.Fatalf("invalid sequence. expected:%v, got:%v", i, block.Timestamp())
		}
		if block.Timestamp() != uint64(i) {
			log.Fatalf("invalid timestamp. expected:%v, got:%v", i, block.Timestamp())
		}
	}

	log.Println("all is good")

	return nil
}
