package main

import (
	"flag"
	"os"
	"runtime/pprof"

	log "github.com/Sirupsen/logrus"
)

// declare the vars to hold the flags
var (
	mmapList        bool
	mmapOne         bool
	memList         bool
	memSlice        bool
	memListEncoded  bool
	memSliceEncoded bool
	loadFile        bool
	loadFileMmap    bool
	optDataLen      int
	optNum          int
	cpuProf         bool
	heapProf        bool
)

func main() {
	// declare our flags
	flag.BoolVar(&mmapList, "mmap-list", false, "write & read capnp list from mmap'ed file")
	flag.BoolVar(&mmapOne, "mmap-one", false, "write one  and read one capnp from mmap'ed file")
	flag.BoolVar(&memList, "mem-list", false, "check memory usage of capnp list")
	flag.BoolVar(&memSlice, "mem-slice", false, "check memory usage of capnp stored in Go slice")
	flag.BoolVar(&memListEncoded, "mem-list-encoded", false, "check memory usage of encoded capnp list")
	flag.BoolVar(&memSliceEncoded, "mem-slice-encoded", false, "check memory usage of encoded capnp in Go slice")
	flag.BoolVar(&loadFile, "load-file", false, "load plain file")
	flag.BoolVar(&loadFileMmap, "load-file-mmap", false, "load mmap'ed file")
	flag.IntVar(&optDataLen, "data-len", 0, "number of bytes of data to add to the capnp message (default 0)")
	flag.IntVar(&optNum, "num", 1000*1000, "number of messages")
	flag.BoolVar(&cpuProf, "cpu-prof", false, "cpu profiling")
	flag.BoolVar(&heapProf, "heap-prof", false, "heap profiling")

	// and parse them
	flag.Parse()

	// we should run at least one test
	if !mmapList && !mmapOne && !memList && !memSlice && !memListEncoded && !memSliceEncoded && !loadFile && !loadFileMmap {
		log.Info("please specify test to perform")
		log.Info("run with '-h' option to see all available tests")
		return
	}

	// enable cpu profiling if the plag is set
	if cpuProf {
		// dump profile in this file to later use with `go tool pprof`
		f, err := os.Create("app.cpuprof")
		if err != nil {
			// Exit if we cant create the profile
			log.Fatalf("failed to create profiling file: %v", err)
		}
		// close file when we are done
		defer f.Close()
		// start profile
		pprof.StartCPUProfile(f)
		// make sure we stop the profiling
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

	if memSlice {
		checkMemUsageSlice(num)
	}

	if memListEncoded {
		checkMemUsageListEncoded(num)
	}

	if memSliceEncoded {
		checkMemUsageSliceEncoded(num)
	}

	if loadFile {
		perfLoadFile(num, false)
	}

	if loadFileMmap {
		perfLoadFile(num, true)
	}

	// dump the heap profile if the flag is set
	if heapProf {
		// create file to dump the profile
		f, err := os.Create("app.mprof")
		if err != nil {
			// exit if the file can't be created. not that the tests will already be done
			// since this is the last statement in the main function
			log.Fatalf("failed to create profiling file: %v", err)
		}
		// dump the heap profile
		pprof.WriteHeapProfile(f)
		// close the file
		f.Close()
	}
}
