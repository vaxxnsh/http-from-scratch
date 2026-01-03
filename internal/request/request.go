package request

import (
	"bytes"
	"fmt"
	"io"

	"github.com/vaxxnsh/http-from-scratch/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateError   parserState = "error"
	StateHeaders parserState = "headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
}

var MALFORMED_START_LINE = fmt.Errorf("malformed start-line")
var UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	before, _, ok := bytes.Cut(b, SEPARATOR)
	if !ok {
		return nil, 0, nil
	}

	startLine := before
	read := len(SEPARATOR) + len(before)

	parts := bytes.Split(startLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, MALFORMED_START_LINE
	}

	r := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2]),
	}

	if !r.ValidHTTP() {
		return nil, 0, UNSUPPORTED_HTTP_VERSION
	}

	versionParts := bytes.Split(parts[2], []byte("/"))

	if len(versionParts) < 2 {
		return nil, 0, MALFORMED_START_LINE
	}
	r.HttpVersion = string(versionParts[1])

	return &r, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		fmt.Printf("state: %s\ndata: %s\n", r.state, string(currentData))
		switch r.state {
		case StateError:
			return 0, REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				return 0, err
			}
			if done {
				r.state = StateDone
			}
			if n == 0 {
				break outer
			}
			read += n
		case StateDone:
			break outer
		default:
			panic("somehow we have programmed poorly")
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func (r *Request) error() bool {
	return r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.done() && !request.error() {
		n, err := reader.Read(buf[bufLen:])
		if n > 0 {
			bufLen += n
			readN, err := request.parse(buf[:bufLen])
			if err != nil {
				return nil, err
			}
			copy(buf, buf[readN:bufLen])
			bufLen -= readN
		}

		if err != nil {
			if err == io.EOF {
				if !request.done() {
					return nil, io.ErrUnexpectedEOF
				}
				break
			}
			return nil, err
		}
	}

	return request, nil
}
