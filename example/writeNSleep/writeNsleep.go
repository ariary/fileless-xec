package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	filename := "writeNsleep.txt"

	f, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	fmt.Println("Writing in", filename)
	for i := 0; i < 10; i++ {
		_, err2 := f.WriteString("Hello fileless world\n")
		time.Sleep(3 * time.Second)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	fmt.Println("done")
}
