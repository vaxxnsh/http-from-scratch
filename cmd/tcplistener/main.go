package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

const SERVER_PORT = ":42069"

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		part := ""

		for {
			buf := make([]byte, 8)
			n, err := f.Read(buf)
			if err != nil {
				break
			}

			data := buf[:n]

			for {
				idx := bytes.IndexByte(data, '\n')
				if idx == -1 {
					part += string(data)
					break
				}

				part += string(bytes.TrimSuffix(data[:idx], []byte("\r")))
				out <- part
				part = ""

				data = data[idx+1:]
			}
		}

		if len(part) > 0 {
			out <- part
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", SERVER_PORT)

	if err != nil {
		log.Fatalf("error:  %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error:  %s", err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("'%s'\n", line)
		}
	}
}
