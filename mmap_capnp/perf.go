package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
)

func perfLoadFile(num int, isMmap bool) {
	log.Println("------  load & walk capnp message")

	// create the file
	f, err := createCapnpFile("/tmp/capnp_perf", num)
	if err != nil {
		log.Fatalf("failed to create capnp file:%v", err)
	}

	defer func() {
		name := f.Name()
		f.Close()
		os.Remove(name)
	}()

	// decode it back
	r, err := getReader(f, isMmap)
	if err != nil {
		log.Fatalf("failed to create reader:%v", err)
	}

	start := time.Now()

	_, blocks, err := decodeAggBlocks(r)
	if err != nil {
		log.Fatalf("failed to decode:%v", err)
	}
	for i := 0; i < blocks.Len(); i++ {
		_ = blocks.At(i)
	}

	fmt.Printf("number of messages:%v\n", humanize.Comma(int64(num)))
	fmt.Printf("mmap : %v\n", isMmap)
	fmt.Printf("file size:%v bytes\n", humanize.Comma(getFileSize(f)))
	fmt.Printf("time:%v seconds\n", time.Since(start).Seconds())
}

func getReader(f *os.File, isMmap bool) (io.Reader, error) {
	if !isMmap {
		f.Seek(0, 0)
		return bufio.NewReader(f), nil
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(getFileSize(f)), syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func createCapnpFile(filename string, num int) (*os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	w := bufio.NewWriter(f)

	return f, writeList(num, w)
}

func getFileSize(f *os.File) int64 {
	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("failed to get file size of '%v' err:%v", f.Name, err)
	}
	return fi.Size()
}
