package main

import (
	"fmt"
	"log"
	"os"
)

const FILE_PATH = "messages.txt"

func main() {
	f, err := os.Open(FILE_PATH)

	if err != nil {
		log.Fatalf("error: couldn't open the file %s\n%s", FILE_PATH, err)
	}

	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			break
		}
		fmt.Printf("read: %v\n", string(data[:n]))
	}
}
