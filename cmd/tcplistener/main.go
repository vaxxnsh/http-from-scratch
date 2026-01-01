package main

import (
	"fmt"
	"log"
	"net"

	"github.com/vaxxnsh/http-from-scratch/internal/request"
)

const SERVER_PORT = ":42069"

func main() {
	listener, err := net.Listen("tcp", SERVER_PORT)

	if err != nil {
		log.Fatalf("error:  %s\n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error:  %s\n", err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		fmt.Printf("Request Line: \n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
