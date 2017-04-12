package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"

	"github.com/dustin/go-humanize"
)

const (
	tlogBlockSize = 20
)

func checkMemUsageList(num int) {
	log.Println("------- check memory usage of in-memory canpnp list -----")

	// allocate memory
	bufSize := (num * tlogBlockSize) + 100
	bs := make([]byte, bufSize)
	buf := bytes.NewBuffer(bs)

	fmt.Printf("buffer size:%v bytes -> it is not Go dependent, but capnp dependent\n", humanize.Comma(int64(bufSize)))

	var memBuf runtime.MemStats
	runtime.ReadMemStats(&memBuf)

	fmt.Printf("total memory allocation:\n")
	fmt.Printf("\tbuffer: %v bytes\n",
		humanize.Comma(int64(memBuf.TotalAlloc)))
	fmt.Printf("\tbuffer overhead: %v bytes\n",
		humanize.Comma(int64(memBuf.TotalAlloc-uint64(bufSize))))

	if err := writeList(num, buf); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	allocated := memWrite.TotalAlloc - memBuf.TotalAlloc
	fmt.Printf("\tcapnp encode:%v bytes\n", humanize.Comma(int64(allocated)))

	// decode it
	_, blocks, err := decodeAggBlocks(buf)
	if err != nil {
		log.Fatalf("failed to decode:%v", err)
	}

	var memRead runtime.MemStats
	runtime.ReadMemStats(&memRead)

	allocated = memRead.TotalAlloc - memWrite.TotalAlloc
	fmt.Printf("\tcapnp decode:%v bytes\n", humanize.Comma(int64(allocated)))

	// make sure we can really decode it
	for i := 0; i < num; i++ {
		block := blocks.At(i)
		checkBlockVal(&block, i)
	}
}

func checkMemUsageMap(num int) {
	log.Println("------- check memory usage of in-memory canpnp stored in Go map -----")

	container := make(map[int][]byte, num)
	for i := 0; i < num; i++ {
		container[i] = make([]byte, tlogBlockSize)
	}

	bufSize := tlogBlockSize * num
	fmt.Printf("buffer size:%v bytes -> it is not Go dependent, but capnp dependent\n", humanize.Comma(int64(bufSize)))

	var memMap runtime.MemStats
	runtime.ReadMemStats(&memMap)

	fmt.Printf("total memory allocation:\n")
	fmt.Printf("\tmap of buffer: %v bytes\n",
		humanize.Comma(int64(memMap.TotalAlloc)))
	fmt.Printf("\tbuffer overhead: %v bytes\n",
		humanize.Comma(int64(memMap.TotalAlloc-uint64(bufSize))))

	// encode it
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		buf.Truncate(0)
		writeBlock(buf, i, tlogBlockSize)
	}

	var memEncode runtime.MemStats
	runtime.ReadMemStats(&memEncode)

	allocated := memEncode.TotalAlloc - memMap.TotalAlloc
	fmt.Printf("\tcapnp encode:%v bytes\n", humanize.Comma(int64(allocated)))

	// decode it
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		block, err := decodeBlock(buf)
		if err != nil {
			log.Fatalf("failed to decode block: %v :%v\n", i, err)
		}

		// check some of it, no need to check all
		// only make sure we did encode/decode corretlu
		if i < 10 {
			checkBlockVal(block, i)
		}
	}
}

func printMemUsage(mem runtime.MemStats) {
	fmt.Printf("alloc       : %v bytes\n", humanize.Comma(int64(mem.Alloc)))
	fmt.Printf("total alloc : %v bytes\n", humanize.Comma(int64(mem.TotalAlloc)))
	fmt.Printf("heap alloc  : %v bytes\n", humanize.Comma(int64(mem.HeapAlloc)))
	fmt.Printf("heap sys    : %v bytes\n", humanize.Comma(int64(mem.HeapSys)))
}
