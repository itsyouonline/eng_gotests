package main

import (
	"bytes"
	"math/big"
	"runtime"

	log "github.com/Sirupsen/logrus"

	"github.com/dustin/go-humanize"
	"zombiezen.com/go/capnproto2"
)

// tlogBlockSize returns the size in bytes of a tlog Block for this run
func tlogBlockSize() int {
	return dataLenInBlock()
}

// checkMemUsageList checks the memory usage of holding the designated amount of messages
// in memory in a capnp list (wrapped in a TlogAggregation)
func checkMemUsageList(num int) {
	log.Info("------- check memory usage of in-memory canpnp list -----")
	log.Infof("stored data size: %v bytes", dataLenInBlock())
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// record the amount of memory we are currently using
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// create the capnp list with the defined amount of blocks
	if _, err := createList(num); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	// record the amount of memory now that we created our capnp list
	var memList runtime.MemStats
	runtime.ReadMemStats(&memList)

	// print our memory usage.
	printMemDif(memStart, memList)
}

// checkMemUsageSlice checks the memory usage by holding the designated amount of messages
// in memory in a slice. Note that the messages aren't encoded
func checkMemUsageSlice(num int) {
	log.Info("------- check memory usage of in-memory canpnp stored in Go slice -----")
	log.Infof("stored data size: %v bytes", dataLenInBlock())
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// save the current memory usage
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// create our map
	container := make([]*capnp.Message, num)

	var memSlice runtime.MemStats
	runtime.ReadMemStats(&memSlice)

	log.Info("total memory allocation for our slice (overhead):")
	log.Infof("\tslice of buffer: %v bytes",
		humanize.Comma(int64(memSlice.HeapInuse-memStart.HeapInuse)))

	// generate the messages holding the blocks
	for i := 0; i < num; i++ {
		_, msg, err := createBlock(i)
		if err != nil {
			log.Fatal(err)
		}
		container[i] = msg
	}
	// get the memory usage now that we have filled our slice
	var memList runtime.MemStats
	runtime.ReadMemStats(&memList)
	// and calculate and print the difference
	log.Info("additional memory usage after filling our slice with encoded messages: ")
	printMemDif(memSlice, memList)
}

// checkMemUsageListEncoded checks the memory used by holding an encoded capnp list
// in memory. the list contains our actual items
func checkMemUsageListEncoded(num int) {
	log.Info("------- check memory usage of in-memory encoded canpnp list -----")
	log.Infof("stored data size: %v bytes", dataLenInBlock())
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// record base memory usage
	var memBase runtime.MemStats
	runtime.ReadMemStats(&memBase)

	// create our buffer
	buf := new(bytes.Buffer)

	// memory usage after allocation of our buffer
	var memBuf runtime.MemStats
	runtime.ReadMemStats(&memBuf)

	if err := writeList(num, buf); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	// memory usage after encodig and writing to our buffer
	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	log.Info("\tmemory usage of our encoded capnp list filled with blocks: ")
	printMemDif(memBuf, memWrite)

	// read and decode the buffer
	_, blocks, err := decodeAggBlocks(buf)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}

	// verify that we correctly decoded the blocks
	for i := 0; i < num; i++ {
		block := blocks.At(i)
		checkBlockVal(&block, i)
	}
}

// checkMemUsageSliceEncoded checks the memory used for holding encoded capnp messages in
// a go slice
func checkMemUsageSliceEncoded(num int) {
	log.Println("------- check memory usage of in-memory encoded canpnp stored in Go slice -----")
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// record the memory we already used
	var baseMem runtime.MemStats
	runtime.ReadMemStats(&baseMem)

	// initialize our slice
	container := make([][]byte, num)

	// record the memory usage after our slice allocation
	var memSlice runtime.MemStats
	runtime.ReadMemStats(&memSlice)

	log.Info("total memory allocation for our slice (overhead):")
	log.Infof("\tslice of buffer: %v bytes",
		humanize.Comma(int64(memSlice.HeapInuse-baseMem.HeapInuse)))

	// create a new buffer, this buffer will be reused for each block we encode
	buf := new(bytes.Buffer)
	// fill our slice with encoded blocks
	for i := 0; i < num; i++ {
		// make sure the buffer is empty
		buf.Reset()
		// write the block in the buffer
		writeBlock(buf, i)
		// set the byteslice with the encoded block in the slice
		// since we want to reuse the buffer and thus the underlying slice, we need to allocate
		// a new slice and copy the buffer contents
		container[i] = make([]byte, buf.Len(), buf.Len())
		copy(container[i], buf.Bytes())
	}

	// record memory usage now that we encoded the blocks. also includes memory used
	// by the encoded blocks
	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	log.Info("additional memory usage after filling our slice with encoded messages: ")
	printMemDif(memSlice, memWrite)

	// decode our blocks and verify that they are correct.
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		block, err := decodeBlock(buf)
		if err != nil {
			log.Fatalf("failed to decode block: %v: %v", i, err)
		}

		// check some of it, no need to check all
		// only make sure we did encode/decode correctly
		if i < 10 {
			checkBlockVal(block, i)
		}
	}
}

// print mem diff prints the difference in used memory between two MemStats.
func printMemDif(memStart, memEnd runtime.MemStats) {
	log.Info("memory stats:")
	log.Infof("\ttotal alloc : %v", formatMemValue(memEnd.TotalAlloc-memStart.TotalAlloc))
	log.Infof("\theap in use : %v", formatMemValue(memEnd.HeapInuse-memStart.HeapInuse))
	log.Infof("\tstack in use : %v", formatMemValue(memEnd.StackInuse-memStart.StackInuse))
}

// formatMemValue turns large numbers into nicer looking strings
func formatMemValue(val uint64) string {
	x := big.NewInt(int64(val))
	return humanize.BigBytes(x)
}

// printMemUsage prints the memory usage from a MemStat
func printMemUsage(mem runtime.MemStats) {
	log.Infof("alloc       : %v bytes", humanize.Comma(int64(mem.Alloc)))
	log.Infof("total alloc : %v bytes", humanize.Comma(int64(mem.TotalAlloc)))
	log.Infof("heap alloc  : %v bytes", humanize.Comma(int64(mem.HeapAlloc)))
	log.Infof("heap sys    : %v bytes", humanize.Comma(int64(mem.HeapSys)))
}
