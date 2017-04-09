package main

import (
	"fmt"
	"time"

	"github.com/darkhelmet/env"
)

// greeting returns a pleasant, semi-useful greeting.
func greeting() string {
	ss = "Hello world, the time is: " + time.Now().String()
	ss += env.String("PATH")
}

func main() {
	fmt.Println(greeting())
}
