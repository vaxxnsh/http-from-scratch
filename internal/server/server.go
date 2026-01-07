package server

import (
	"fmt"
	"io"
	"net"

	"github.com/vaxxnsh/http-from-scratch/internal/response"
)

type Server struct {
	port   uint16
	closed bool
}

func runConnection(_ *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	body := []byte("Hello World!\n")
	err := response.WriteStatusLine(conn, response.StatusOk)
	if err != nil {
		return
	}
	headers := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}
	conn.Write(body)
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
