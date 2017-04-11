package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	mmapBulk bool
	mmapOne  bool
)

func main() {
	flag.BoolVar(&mmapBulk, "mmap-bulk", false, "write bulk and read bulk")
	flag.BoolVar(&mmapOne, "mmap-one", false, "write one and read one")

	flag.Parse()

	if !mmapBulk && !mmapOne {
		fmt.Println("please specify test to perform")
		fmt.Println("run with '-h' option to see all available tests")
		return
	}

	num := 1000 * 1000

	if mmapBulk {
		if err := writeBulkRead(num); err != nil {
			log.Printf("err = %v\n", err)
		}
	}

	if mmapOne {
		if err := writeOneReadOne(num, 70); err != nil {
			log.Printf("writeOneReadOn err = %v\n", err)
		}
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
	return num * 40
}
