package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/tarantool/go-tarantool"
)

const (
	dataLen = 200
)

func main() {
	var user string
	var passwd string
	var addr string
	var num int
	var numClient int

	flag.StringVar(&user, "user", "", "tarantool username")
	flag.StringVar(&passwd, "passwd", "", "tarantool password")
	flag.StringVar(&addr, "addr", "127.0.0.1:3301", "tarantool server address (default: 127.0.0.1:3301)")
	flag.IntVar(&num, "num", 1000*1000, "number of messages to store")
	flag.IntVar(&numClient, "num-client", runtime.NumCPU(), "number of clients (default:number of cpu cores)")

	flag.Parse()

	fmt.Println("test data:")
	fmt.Printf("\t number of messages:%v\n", num)
	fmt.Printf("\t data len = %v bytes\n", dataLen)
	fmt.Printf("\t number of clients = %v\n", numClient)

	opts := tarantool.Opts{
		User: user,
		Pass: passwd,
	}

	space := "xxx"
	index := "id"

	data := make([]byte, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = 10
	}

	start := time.Now()

	idChan := make(chan int, numClient)
	closeChan := make(chan bool)
	go func() {
		for i := 0; i < num; i++ {
			idChan <- i
		}
		close(closeChan)
	}()

	var wg sync.WaitGroup
	for i := 0; i < numClient; i++ {
		wg.Add(1)
		go func() {
			conn, err := tarantool.Connect(addr, opts)
			if err != nil {
				log.Fatalf("Connection refused: %s", err)
			}

			defer wg.Done()
			for {
				select {
				case idx := <-idChan:
					resp, err := conn.Insert(space, []interface{}{idx, data})
					if err != nil {
						log.Fatalf("Code:%v, err:%v", resp.Code, err)
					}
				case <-closeChan:
					return
				}
			}
		}()
	}
	wg.Wait()

	fmt.Printf("\t time needed:%v seconds\n", time.Since(start).Seconds())

	// test getting some data we have stored
	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		log.Fatalf("Connection refused: %s", err)
	}

	resp, err := conn.Select(space, index, 0, 1, tarantool.IterEq, []interface{}{0})
	if err != nil {
		log.Fatalf("Code:%v, err:%v", resp.Code, err)
	}
}
