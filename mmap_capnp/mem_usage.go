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
	var mem runtime.MemStats

	bs := make([]byte, (num*tlogBlockSize)+100)
	buf := bytes.NewBuffer(bs)

	log.Println("---mem stats after allocating buffer that will be used by capnp list--")
	log.Println("---it is not go dependent, but capnp dependent---")
	runtime.ReadMemStats(&mem)
	printMemUsage(mem)
	startAlloc := mem.Alloc

	if err := writeList(num, buf); err != nil {
		log.Fatalf("failed to write to capnp list")
	}

	log.Println("---mem stats after write capnp lists to memory--")

	runtime.ReadMemStats(&mem)
	printMemUsage(mem)

	allocated := mem.Alloc - startAlloc
	fmt.Printf("----> there are %v bytes allocated\n", humanize.Comma(int64(allocated)))

	// make sure we can decode it
}

func checkMemUsageMap(num int) {
	var mem runtime.MemStats
	container := map[int][]byte{}
	for i := 0; i < num; i++ {
		container[i] = make([]byte, tlogBlockSize)
	}

	log.Println("---mem stats after allocating buffer that will be used by capnp list--")
	log.Println("---other than Go Map memory usage, other memory usage is needed by capnp format---")
	runtime.ReadMemStats(&mem)
	printMemUsage(mem)
	startAlloc := mem.Alloc

	for i := 0; i < num; i++ {
		buf := bytes.NewBuffer(container[i])
		writeBlock(buf, i, tlogBlockSize)
	}

	log.Println("---mem stats after write capnp lists to memory--")

	runtime.ReadMemStats(&mem)
	printMemUsage(mem)

	allocated := mem.Alloc - startAlloc
	fmt.Printf("----> there are %v bytes allocated\n", humanize.Comma(int64(allocated)))

}

func printMemUsage(mem runtime.MemStats) {
	fmt.Printf("alloc       : %v bytes\n", humanize.Comma(int64(mem.Alloc)))
	fmt.Printf("total alloc : %v bytes\n", humanize.Comma(int64(mem.TotalAlloc)))
	fmt.Printf("heap alloc  : %v bytes\n", humanize.Comma(int64(mem.HeapAlloc)))
	fmt.Printf("heap sys    : %v bytes\n", humanize.Comma(int64(mem.HeapSys)))
}
