package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var count int
	flag.IntVar(&count, "n", 5, "number of hello world messages")
	flag.Parse()
	for i := 0; i < count; i++ {
		fmt.Println("Hello world")
		time.Sleep(1 * time.Second)
	}
}
