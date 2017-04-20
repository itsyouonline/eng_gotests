package main

import (
	"flag"
	"runtime"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/tarantool/go-tarantool"
)

// test with 200 bytes worth of data
const (
	dataLen = 200
)

func main() {
	// variables for flags
	var user string
	var passwd string
	var addr string
	var num int
	var numClient int

	// flag declarations
	flag.StringVar(&user, "user", "", "tarantool username")
	flag.StringVar(&passwd, "passwd", "", "tarantool password")
	flag.StringVar(&addr, "addr", "127.0.0.1:3302", "tarantool server address")
	flag.IntVar(&num, "num", 1000*1000, "number of messages to store")
	flag.IntVar(&numClient, "num-client", runtime.NumCPU(), "number of clients")

	// parse the flags
	flag.Parse()

	// print some info about what is going to happen
	log.Info("test data:")
	log.Infof("\t number of messages: %v", num)
	log.Infof("\t data len = %v bytes", dataLen)
	log.Infof("\t number of clients = %v", numClient)

	// options struct for the tarantool library
	opts := tarantool.Opts{
		User: user,
		Pass: passwd,
	}

	space := "xxx"
	index := "id"

	// some junk data
	data := make([]byte, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = 10
	}

	// record the current time
	start := time.Now()

	// open a buffered channel to send ints...
	idChan := make(chan int, numClient)
	// and an unbuffered channel to send bools
	closeChan := make(chan bool)
	// launch a closure in a seperate goroutine
	go func() {
		for i := 0; i < num; i++ {
			// send ints on the channel. this blocks if the channel buffer is full
			idChan <- i
		}
		// close the `closeChan` channel, this will cause listerners to receive a zero value
		// if they attempt to receive something
		close(closeChan)
	}()

	var wg sync.WaitGroup
	for i := 0; i < numClient; i++ {
		// for every client goroutine, increase the wait group counter.
		wg.Add(1)
		// launch a new goroutine
		go func() {
			// connect to the db
			conn, err := tarantool.Connect(addr, opts)
			if err != nil {
				log.Fatalf("Connection refused: %s", err)
			}
			// when the function and thus the goroutine exits, notify the waitgroup
			defer wg.Done()
			// infinite loop
			for {
				// pull items from the channels.
				// try to get an item from idChan. if there currently isn't an item there,
				// try and get one from the closeChan. Repeat forever
				select {
				case idx := <-idChan:
					// store the data
					_, err := conn.Insert(space, []interface{}{idx, data})
					if err != nil {
						log.Fatalf("err: %v", err)
					}
				// if we receive something here, return, causing the function and goroutine
				// to exit. our defered wg.Done() is called to notify the wait group. since we
				// accept anything here, the zero value from listening after close has been called
				// on the channel effectively causes all our goroutines to exit after the channel
				// is closed.
				case <-closeChan:
					return
				}
			}
		}()
	}
	// block until the specified amount of wg.Done() calls have occureed.
	wg.Wait()

	// log the time it took
	log.Infof("\t time needed: %v seconds\n", time.Since(start).Seconds())

	// test getting some data we have stored
	// get a new connection
	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		log.Fatalf("Connection refused: %s", err)
	}

	// and get some data from the connection
	resp, err := conn.Select(space, index, 0, 1, tarantool.IterEq, []interface{}{0})
	if err != nil {
		log.Fatalf("Code: %v, err: %v", resp.Code, err)
	}
}
