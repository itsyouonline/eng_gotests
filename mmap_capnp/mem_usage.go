package main

import (
	"bytes"
	"math/big"
	"runtime"

	log "github.com/Sirupsen/logrus"

	"github.com/dustin/go-humanize"
	"zombiezen.com/go/capnproto2"
)

func tlogBlockSize() int {
	return dataLenInBlock()
}

func checkMemUsageList(num int) {
	log.Info("------- check memory usage of in-memory canpnp list -----")
	log.Infof("stored data size:%v bytes", dataLenInBlock())
	log.Infof("number of message:%v", humanize.Comma(int64(num)))

	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	if _, err := createList(num); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	var memList runtime.MemStats
	runtime.ReadMemStats(&memList)

	printMemDif(memStart, memList)
}

func checkMemUsageMap(num int) {
	log.Info("------- check memory usage of in-memory canpnp stored in Go map -----")
	log.Infof("stored data size:%v bytes", dataLenInBlock())
	log.Infof("number of message:%v", humanize.Comma(int64(num)))

	var memStart runtime.MemStats
	runtime.ReadMemStats(&memStart)

	container := make(map[int]*capnp.Message, num)

	// encode it
	for i := 0; i < num; i++ {
		_, msg, err := createBlock(i)
		if err != nil {
			log.Fatal(err)
		}
		container[i] = msg
	}
	var memList runtime.MemStats
	runtime.ReadMemStats(&memList)

	printMemDif(memStart, memList)
}

func checkMemUsageListEncoded(num int) {
	log.Info("------- check memory usage of in-memory encoded canpnp list -----")
	log.Infof("stored data size:%v bytes", dataLenInBlock())
	log.Infof("number of message:%v", humanize.Comma(int64(num)))

	// allocate memory
	bufSize := (num * tlogBlockSize()) + 100
	bs := make([]byte, bufSize)
	buf := bytes.NewBuffer(bs)

	log.Infof("buffer size:%v bytes -> it is not Go dependent, but capnp dependent\n", humanize.Comma(int64(bufSize)))

	var memBuf runtime.MemStats
	runtime.ReadMemStats(&memBuf)

	log.Info("total memory allocation:")
	log.Infof("\tbuffer: %v bytes",
		humanize.Comma(int64(memBuf.TotalAlloc)))
	log.Infof("\tbuffer overhead: %v bytes",
		humanize.Comma(int64(memBuf.TotalAlloc-uint64(bufSize))))

	buf.Truncate(0)
	if err := writeList(num, buf); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	var memWrite runtime.MemStats
	runtime.ReadMemStats(&memWrite)

	allocated := memWrite.TotalAlloc - memBuf.TotalAlloc
	log.Infof("\tcapnp encode:%v bytes", humanize.Comma(int64(allocated)))

	// decode it
	_, blocks, err := decodeAggBlocks(buf)
	if err != nil {
		log.Fatalf("failed to decode:%v", err)
	}

	var memRead runtime.MemStats
	runtime.ReadMemStats(&memRead)

	allocated = memRead.TotalAlloc - memWrite.TotalAlloc
	log.Infof("\tcapnp decode:%v bytes", humanize.Comma(int64(allocated)))

	// make sure we can really decode it
	for i := 0; i < num; i++ {
		block := blocks.At(i)
		checkBlockVal(&block, i)
	}
}

func checkMemUsageMapEncoded(num int) {
	log.Println("------- check memory usage of in-memory encoded canpnp stored in Go map -----")

	container := make(map[int][]byte, num)
	for i := 0; i < num; i++ {
		container[i] = make([]byte, tlogBlockSize())
	}

	bufSize := tlogBlockSize() * num
	log.Infof("buffer size:%v bytes -> it is not Go dependent, but capnp dependent", humanize.Comma(int64(bufSize)))

	var memMap runtime.MemStats
	runtime.ReadMemStats(&memMap)

	log.Info("total memory allocation:")
	log.Infof("\tmap of buffer: %v bytes",
		humanize.Comma(int64(memMap.TotalAlloc)))
	log.Infof("\tbuffer overhead: %v bytes",
		humanize.Comma(int64(memMap.TotalAlloc-uint64(bufSize))))

	// encode it
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		buf.Truncate(0)
		writeBlock(buf, i)
	}

	var memEncode runtime.MemStats
	runtime.ReadMemStats(&memEncode)

	allocated := memEncode.TotalAlloc - memMap.TotalAlloc
	log.Infof("\tcapnp encode:%v bytes", humanize.Comma(int64(allocated)))

	// decode it
	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		block, err := decodeBlock(buf)
		if err != nil {
			log.Fatalf("failed to decode block: %v :%v", i, err)
		}

		// check some of it, no need to check all
		// only make sure we did encode/decode corretlu
		if i < 10 {
			checkBlockVal(block, i)
		}
	}
}

func printMemDif(memStart, memEnd runtime.MemStats) {
	log.Info("memory stats:")
	log.Infof("\ttotal alloc : %v", formatMemValue(memEnd.TotalAlloc-memStart.TotalAlloc))
	log.Infof("\theap in use : %v", formatMemValue(memEnd.HeapInuse-memStart.HeapInuse))
	log.Infof("\tstack in use : %v", formatMemValue(memEnd.StackInuse-memStart.StackInuse))
}

func formatMemValue(val uint64) string {
	x := big.NewInt(int64(val))
	return humanize.BigBytes(x)
}
func printMemUsage(mem runtime.MemStats) {
	log.Infof("alloc       : %v bytes", humanize.Comma(int64(mem.Alloc)))
	log.Infof("total alloc : %v bytes", humanize.Comma(int64(mem.TotalAlloc)))
	log.Infof("heap alloc  : %v bytes", humanize.Comma(int64(mem.HeapAlloc)))
	log.Infof("heap sys    : %v bytes", humanize.Comma(int64(mem.HeapSys)))
}
