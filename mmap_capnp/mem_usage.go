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

// checkMemUsageMap checks the memory usage by holding the designated amount of messages
// in memory in a map. Note that the messages aren't encoded
func checkMemUsageMap(num int) {
	log.Info("------- check memory usage of in-memory canpnp stored in Go map -----")
	log.Infof("stored data size:%v bytes", dataLenInBlock())
	log.Infof("number of message:%v", humanize.Comma(int64(num)))

	// save the current memory usage
	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	// create our map
	container := make(map[int]*capnp.Message, num)

	// generate the blocks
	for i := 0; i < num; i++ {
		_, msg, err := createBlock(i)
		if err != nil {
			log.Fatal(err)
		}
		container[i] = msg
	}
	// get the memory usage now that we have filled our map
	var memList runtime.MemStats
	runtime.ReadMemStats(&memList)
	// and calculate and print the difference
	printMemDif(memStart, memList)
}

// checkMemUsageListEncoded checks the memory used to encode and decode a capnp list
// the reported memory is not used by the list itself, only by the write/encode and
// read/decode functions
func checkMemUsageListEncoded(num int) {
	log.Info("------- check memory usage of in-memory encoded canpnp list -----")
	log.Infof("stored data size: %v bytes", dataLenInBlock())
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// record base memory usage
	var memBase runtime.MemStats
	runtime.ReadMemStats(&memBase)

	// allocate memory
	bufSize := (num * tlogBlockSize()) + 100
	bs := make([]byte, bufSize)
	buf := bytes.NewBuffer(bs)

	log.Infof("buffer size: %v bytes -> it is not Go dependent, but capnp dependent\n", humanize.Comma(int64(bufSize)))

	// memory usage after allocation of our buffer
	var memBuf runtime.MemStats
	runtime.ReadMemStats(&memBuf)

	log.Info("total memory allocation:")
	log.Infof("\tbuffer: %v bytes",
		humanize.Comma(int64(memBuf.TotalAlloc-memBase.TotalAlloc)))
	log.Infof("\tbuffer overhead: %v bytes",
		humanize.Comma(int64(memBuf.TotalAlloc-memBase.TotalAlloc-uint64(bufSize))))

	buf.Truncate(0)
	if err := writeList(num, buf); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	// memory usage after encodig and writing to our buffer
	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	allocated := memWrite.TotalAlloc - memBuf.TotalAlloc
	log.Infof("\tmemory allocated while encoding and writing the capnp list: %v bytes", humanize.Comma(int64(allocated)))

	// read and decode the buffer
	_, blocks, err := decodeAggBlocks(buf)
	if err != nil {
		log.Fatalf("failed to decode:%v", err)
	}

	// memory usage after reading and decoding the buffer
	var memRead runtime.MemStats
	runtime.ReadMemStats(&memRead)

	allocated = memRead.TotalAlloc - memWrite.TotalAlloc
	log.Infof("\tmemory allocated while reading and decoding the capnp list: %v bytes", humanize.Comma(int64(allocated)))

	// verify that we correctly decoded the blocks
	for i := 0; i < num; i++ {
		block := blocks.At(i)
		checkBlockVal(&block, i)
	}
}

// checkMemUsageMapEncoded checks the memory used for encoding / decoding the designated amount
// of capnp messages stored in a map. reported memory is not used by the messages and map,
// only by the encoding and decoding funcions
func checkMemUsageMapEncoded(num int) {
	log.Println("------- check memory usage of in-memory encoded canpnp stored in Go map -----")
	log.Infof("number of message: %v", humanize.Comma(int64(num)))

	// record the memory we already used
	var baseMem runtime.MemStats
	runtime.ReadMemStats(&baseMem)

	container := make(map[int][]byte, num)
	for i := 0; i < num; i++ {
		container[i] = make([]byte, tlogBlockSize()+30)
	}

	// record the memory usage after our map allocation
	var memMap runtime.MemStats
	runtime.ReadMemStats(&memMap)

	bufSize := (tlogBlockSize() + 30) * num
	log.Infof("buffer size: %v bytes -> it is not Go dependent, but capnp dependent", humanize.Comma(int64(bufSize)))

	log.Info("total memory allocation:")
	log.Infof("\tmap of buffer: %v bytes",
		humanize.Comma(int64(memMap.TotalAlloc-baseMem.TotalAlloc)))
	log.Infof("\tbuffer overhead: %v bytes",
		humanize.Comma(int64(memMap.TotalAlloc-baseMem.TotalAlloc-uint64(bufSize))))

	// fill our map with encoded blocks
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		buf.Truncate(0)
		writeBlock(buf, i)
	}

	// record memory usage now that we encoded the blocks. This memory usage does
	// not include the blocks since we already allocated space for them when we
	// initialized the map
	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	allocated := memWrite.TotalAlloc - memMap.TotalAlloc
	log.Infof("\tmemory allocated while encoding capnp messages in a map: %v bytes", humanize.Comma(int64(allocated)))

	// decode our blocks and verify that they are correct.
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		block, err := decodeBlock(buf)
		if err != nil {
			log.Fatalf("failed to decode block: %v: %v", i, err)
		}

		// check some of it, no need to check all
		// only make sure we did encode/decode corretly
		if i < 10 {
			checkBlockVal(block, i)
		}
	}

	// memory usage after reading and decoding the buffer
	var memRead runtime.MemStats
	runtime.ReadMemStats(&memRead)

	allocated = memRead.TotalAlloc - memWrite.TotalAlloc
	log.Infof("\tmemory allocated while decoding capnp messages in a map: %v bytes", humanize.Comma(int64(allocated)))

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
