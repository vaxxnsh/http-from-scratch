package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, &headers)
		return
	}
	writer := bytes.NewBuffer([]byte{})
	handleError := s.handler(writer, r)
	if handleError != nil {
		response.WriteStatusLine(conn, handleError.StatusCode)
		writer.Write([]byte(handleError.Messsage))
	}
	body := writer.Bytes()
	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, &headers)
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
