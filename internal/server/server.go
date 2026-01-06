package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	port uint16
}

func runConnection(s *Server, conn io.ReadCloser) {

}

func runServer(s *Server, listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go runConnection(s, conn)
	}

	return nil
}

func Serve(port uint16) (*Server, error) {
	server := &Server{
		port: port,
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
	return nil
}
