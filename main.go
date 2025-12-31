package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

const FILE_PATH = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(out)

		part := ""

		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err != nil {
				break
			}
			if newLineIdx := bytes.IndexByte(data, '\n'); newLineIdx != -1 {
				part += string(data[:newLineIdx])
				data = data[newLineIdx+1:]
				out <- part
				part = ""
			}
			part += string(data)
		}

		if len(part) != 0 {
			out <- part
		}
	}()

	return out
}

func main() {
	f, err := os.Open(FILE_PATH)

	if err != nil {
		log.Fatalf("error: couldn't open the file %s\n%s", FILE_PATH, err)
	}

	lines := getLinesChannel(f)

	for line := range lines {
		fmt.Printf("read: %v\n", line)
	}
}
