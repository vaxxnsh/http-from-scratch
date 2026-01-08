package server

import (
	"fmt"
	"io"
	"net"

	"github.com/vaxxnsh/http-from-scratch/internal/request"
	"github.com/vaxxnsh/http-from-scratch/internal/response"
)

type Server struct {
	port    uint16
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Messsage   string
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		headers := response.GetDefaultHeaders(0)
		responseWriter.WriteHeaders(&headers)
		return
	}

	s.handler(responseWriter, r)
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

func Serve(port uint16, handler Handler) (*Server, error) {
	server := &Server{
		port:    port,
		closed:  false,
		handler: handler,
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
