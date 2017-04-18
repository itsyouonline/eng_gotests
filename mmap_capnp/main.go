package main

import (
	"flag"
	"os"
	"runtime/pprof"

	log "github.com/Sirupsen/logrus"
)

// declare the vars to hold the flags
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
	cpuProf        bool
	heapProf       bool
)

func main() {
	// declare our flags
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
	flag.BoolVar(&cpuProf, "cpu-prof", false, "cpu profiling")
	flag.BoolVar(&heapProf, "heap-prof", false, "heap profiling")

	// and parse them
	flag.Parse()

	// we should run at least one test
	if !mmapList && !mmapOne && !memList && !memMap && !memListEncoded && !memMapEncoded && !loadFile && !loadFileMmap {
		log.Info("please specify test to perform")
		log.Info("run with '-h' option to see all available tests")
		return
	}

	if cpuProf {
		f, err := os.Create("app.cpuprof")
		if err != nil {
			log.Fatalf("failed to create profiling file: %v", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// get the amount of tlog blocks we want to use from the package lvl variable
	num := optNum

	// run the enabled tests

	if mmapList {
		if err := writeListRead(num); err != nil {
			log.Infof("err = %v", err)
		}
	}

	if mmapOne {
		if err := writeOneReadOne(num); err != nil {
			log.Infof("writeOneReadOn err = %v", err)
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
		checkMemUsageMapEncoded(num)
	}

	if loadFile {
		perfLoadFile(num, false)
	}
	if loadFileMmap {
		perfLoadFile(num, true)
	}

	if heapProf {
		f, err := os.Create("app.mprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
