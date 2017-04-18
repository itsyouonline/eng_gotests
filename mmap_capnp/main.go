package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var (
	mmapList       bool
	mmapOne        bool
	memList        bool
	memMap         bool
	memListEncoded bool
	memMapEncoded  bool
	loadFile       bool
	loadFileMmap   bool
	optDataLen     int
	optNum         int
)

func main() {
	flag.BoolVar(&mmapList, "mmap-list", false, "write & read capnp list from mmap'ed file")
	flag.BoolVar(&mmapOne, "mmap-one", false, "write one  and read one capnp from mmap'ed file")
	flag.BoolVar(&memList, "mem-list", false, "check memory usage of capnp list")
	flag.BoolVar(&memMap, "mem-map", false, "check memory usage of capnp stored in Go map")
	flag.BoolVar(&memListEncoded, "mem-list-encoded", false, "check memory usage of encoded capnp list")
	flag.BoolVar(&memMapEncoded, "mem-map-encoded", false, "check memory usage of encoded capnp in Go map")
	flag.BoolVar(&loadFile, "load-file", false, "load plain file")
	flag.BoolVar(&loadFileMmap, "load-file-mmap", false, "load mmap'ed file")
	flag.IntVar(&optDataLen, "data-len", 0, "number of bytes of data to add to the capnp message(default = 0)")
	flag.IntVar(&optNum, "num", 1000*1000, "number of messages (default = 1M)")

	flag.Parse()

	if !mmapList && !mmapOne && !memList && !memMap && !memListEncoded && !memMapEncoded && !loadFile && !loadFileMmap {
		fmt.Println("please specify test to perform")
		fmt.Println("run with '-h' option to see all available tests")
		return
	}

	num := optNum

	if mmapList {
		if err := writeListRead(num); err != nil {
			log.Printf("err = %v\n", err)
		}
	}

	if mmapOne {
		if err := writeOneReadOne(num); err != nil {
			log.Printf("writeOneReadOn err = %v\n", err)
		}
	}

	if memList {
		checkMemUsageList(num)
	}
	if memMap {
		checkMemUsageMap(num)
	}
	if memListEncoded {
		checkMemUsageListEncoded(num)
	}
	if memMapEncoded {
		checkMemUsageMapEncoded(num / 1000)
	}

	if loadFile {
		perfLoadFile(num, false)
	}
	if loadFileMmap {
		perfLoadFile(num, true)
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
