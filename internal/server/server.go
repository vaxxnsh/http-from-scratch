package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	port   uint16
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")
	conn.Write(out)
	conn.Close()
}

func runServer(s *Server, listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return nil
		}

		if err != nil {
			return err
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	server := &Server{
		port:   port,
		closed: false,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))

	if err != nil {
		return nil, err
	}

	err = runServer(server, listener)

	if err != nil {
		return nil, err
	}

	return &Server{}, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
