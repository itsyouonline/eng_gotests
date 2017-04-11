package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	mmapList bool
	mmapOne  bool
	memList  bool
)

func main() {
	flag.BoolVar(&mmapList, "mmap-list", false, "write & read capnp list from mmap'ed file")
	flag.BoolVar(&mmapOne, "mmap-one", false, "write one  and read one capnp from mmap'ed file")
	flag.BoolVar(&memList, "mem-list", false, "check memory usage of capnp list")

	flag.Parse()

	if !mmapList && !mmapOne && !memList {
		fmt.Println("please specify test to perform")
		fmt.Println("run with '-h' option to see all available tests")
		return
	}

	num := 1000 * 1000

	if mmapList {
		if err := writeListRead(num); err != nil {
			log.Printf("err = %v\n", err)
		}
	}

	if mmapOne {
		if err := writeOneReadOne(num, 70); err != nil {
			log.Printf("writeOneReadOn err = %v\n", err)
		}
	}

	if memList {
		checkMemUsageList(num)
	}
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
	return num * 20
}
