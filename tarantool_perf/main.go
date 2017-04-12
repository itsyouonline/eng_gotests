package main

import (
	"flag"
	"fmt"
	"log"
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

	flag.StringVar(&user, "user", "", "tarantool username")
	flag.StringVar(&passwd, "passwd", "", "tarantool password")
	flag.StringVar(&addr, "addr", "127.0.0.1:3301", "tarantool server address (default: 127.0.0.1:3301)")
	flag.IntVar(&num, "num", 1000*1000, "number of messages to store")

	flag.Parse()

	fmt.Println("test data:")
	fmt.Printf("\t number of messages:%v\n", num)
	fmt.Printf("\t data len = %v bytes\n", dataLen)

	opts := tarantool.Opts{
		User: user,
		Pass: passwd,
	}

	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		log.Fatalf("Connection refused: %s", err)
	}

	space := 512
	index := "id"

	data := make([]byte, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = 10
	}

	start := time.Now()

	for i := 0; i < num; i++ {
		resp, err := conn.Insert(space, []interface{}{i, data})
		if err != nil {
			log.Fatalf("Code:%v, err:%v", resp.Code, err)
		}
	}
	fmt.Printf("\t time needed:%v seconds\n", time.Since(start).Seconds())

	// test getting some data we have stored
	resp, err := conn.Select(space, index, 0, 1, tarantool.IterEq, []interface{}{0})
	if err != nil {
		log.Fatalf("Code:%v, err:%v", resp.Code, err)
	}
}
