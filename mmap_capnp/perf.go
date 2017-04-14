package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/dustin/go-humanize"
)

// perfLoadFile checks the performance for loading a capnp file from disk, using both
// mmap and just getting a reader for the file.
func perfLoadFile(num int, isMmap bool) {
	log.Println("------  load file & walk over capnp message  ------")

	// create the file
	f, err := createCapnpFile("/tmp/capnp_perf", num)
	if err != nil {
		log.Fatalf("failed to create capnp file: %v", err)
	}

	// make sure we close and remove the file when we are done
	defer func() {
		name := f.Name()
		f.Close()
		os.Remove(name)
	}()

	start := time.Now()

	// decode it back
	r, err := getReader(f, isMmap)
	if err != nil {
		log.Fatalf("failed to create reader:%v", err)
	}

	// decode tlog aggregations from the reader
	_, blocks, err := decodeAggBlocks(r)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}
	// walk over the tlog blocks in the list
	for i := 0; i < blocks.Len(); i++ {
		_ = blocks.At(i)
	}
	// get the time it took
	duration := time.Since(start).Seconds()

	log.Infof("number of messages: %v", humanize.Comma(int64(num)))
	log.Infof("mmap: %v", isMmap)
	log.Infof("file size: %v bytes", humanize.Comma(getFileSize(f)))
	log.Infof("time: %v seconds", duration)
}

// getReader gets an io.Reader for the file. if a mmap reader is requested, the file will
// be mmaped first and then a reader for said mmap is created and returned.
func getReader(f *os.File, isMmap bool) (io.Reader, error) {
	// if we don't need an mmap, set the pointer to the beginning of the file and make a
	// new bufio reader
	if !isMmap {
		f.Seek(0, 0)
		return bufio.NewReader(f), nil
	}

	// else create the mmap first
	data, err := syscall.Mmap(int(f.Fd()), 0, int(getFileSize(f)), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	// create a new buffer from the mmap
	return bytes.NewBuffer(data), nil
}

// createCapnpFile creates a new file with the requested name. if the file is
// created successfully, an attempt is made to populate it with a tlog aggregation
// with the requested amount of tlog blocks. the caller must close the returned file
// pointer when they are done.
func createCapnpFile(filename string, num int) (*os.File, error) {
	// create the file, override existing file if any
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	// get a new buffered writer to the file
	w := bufio.NewWriter(f)

	// return the file pointer, and a possible error while writing the list
	return f, writeList(num, w)
}

// getFileSize returns the size of a file. If there is an error, the program exits
func getFileSize(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("failed to get file size of '%v' err:%v", f.Name(), err)
	}
	return fi.Size()
}
