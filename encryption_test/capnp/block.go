package capnp

import (
	"bytes"
	"crypto/rand"
	"hash/crc32"
	mathrand "math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	capnp "zombiezen.com/go/capnproto2"
)

var dataSize int

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// SetDataSize sets the data size to embed in a tlog block
func SetDataSize(ds int) {
	dataSize = ds
	log.Info("Current amount of actual data in a single block: ", dataLenInBlock())
}

// GenerateBlocks makes a slice of the desired length filled with encoded capnp messages.
// In this case, said messages contain tlog blocks
func GenerateBlocks(amount int) [][]byte {
	buf := make([][]byte, amount)
	writebuf := bytes.NewBuffer(nil)
	for i := range buf {
		err := writeBlock(writebuf, i)
		if err != nil {
			log.Fatalf("Failed to write block %v: %v", i, err)
		}
		buf[i] = writebuf.Bytes()
		writebuf = bytes.NewBuffer(nil)
	}
	return buf
}

// createBlock creates a new tlog block and sets the sequence
func createBlock(i int) (*TlogBlock, *capnp.Message, error) {
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
	setBlockVal(&block, 0, i)
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

// setBlockVal sets the sequence of a TlogBlock and sets the text field if the data
// len has been specified to be bigger than 0 for this run
func setBlockVal(block *TlogBlock, vid int, seq int) {
	block.SetVolumeId(uint32(vid))
	block.SetSequence(uint64(seq))
	block.SetCrc32(uint32(0))
	block.SetSize(uint32(dataSize))
	block.SetLba(mathrand.Uint64())
	if dataSize > 0 {
		b := make([]byte, dataSize)
		_, err := rand.Read(b)
		if err != nil {
			log.Fatal("Failed to read random data for block creation: ", err)
		}
		block.SetData(b)
		block.SetCrc32(uint32(crc32.ChecksumIEEE(b)))
	}
	block.SetTimestamp(uint64(time.Now().UnixNano()))
}

// GenerateList generates and encodes a new tlog aggregation
func GenerateList(amount int) []byte {
	list, err := createList(amount)
	if err != nil {
		log.Fatal("Failed to create list: ", err)
	}
	buf := bytes.NewBuffer(nil)
	err = capnp.NewEncoder(buf).Encode(list)
	if err != nil {
		log.Fatal("Failed to encode list: ", err)
	}
	return buf.Bytes()
}

// createList creates the defined amount of TlogBlocks, stored in a capnp list
// wrapped in a TlogAggregation.
func createList(num int) (*capnp.Message, error) {
	// create an empty byte slice and estimate its total size to increase performance
	buf := make([]byte, 0, 100+num*segmentBufferSize())
	// create a new capnp message
	aggMsg, aggSeg, err := capnp.NewMessage(capnp.SingleSegment(buf))
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
	agg.SetName("test aggregation")
	agg.SetSize(0)
	agg.SetTimestamp(uint64(time.Now().Unix()))
	agg.SetVolumeId(0)
	// create the TlogBlock list
	blockList, err := agg.NewBlocks(int32(num))
	if err != nil {
		return nil, err
	}

	// add blocks to the new list
	for i := 0; i < blockList.Len(); i++ {
		block := blockList.At(i)
		setBlockVal(&block, 0, i)
	}

	// update the size of the aggregation object now that we've added the blocks
	agg.SetSize(uint64(num))
	return aggMsg, nil
}

// dataLenInBlock returns the length of the data in a unenoded capnp tlog block
func dataLenInBlock() int {
	len := 36 /* fixed size fields */ + 4 /* capnp overhead */
	len += dataSize                       /* text */
	return len
}

func segmentBufferSize() int {
	return 8*((dataLenInBlock()/8)+1) + 20
}
