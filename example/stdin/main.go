package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Enter filename:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filename := scanner.Text()
	_, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if scanner.Err() != nil {
		// Handle error.
	}
}
