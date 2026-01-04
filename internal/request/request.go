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
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	state       parserState
	RequestLine RequestLine
	Headers     *headers.Headers
	Body        []byte
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

	read := len(before) + len(SEPARATOR)
	parts := bytes.Split(before, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, MALFORMED_START_LINE
	}

	rl := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(parts[2]),
	}

	if !rl.ValidHTTP() {
		return nil, 0, UNSUPPORTED_HTTP_VERSION
	}

	versionParts := bytes.Split(parts[2], []byte("/"))
	if len(versionParts) < 2 {
		return nil, 0, MALFORMED_START_LINE
	}

	rl.HttpVersion = string(versionParts[1])
	return &rl, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		currentData := data[read:]
		switch r.state {

		case StateError:
			fmt.Printf("got error in state - %s\n", r.state)
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
				r.state = StateError
				return 0, err
			}

			read += n

			if done {
				fmt.Printf("went here state %s\n", r.state)
				if r.hasBody() {
					r.state = StateBody
					break
				}
				r.state = StateDone
				break
			}

			if n == 0 {
				break outer
			}

		case StateBody:
			contentLen, err := r.Headers.GetInt("content-length", 0)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if len(currentData) == 0 {
				break outer
			}

			remaining := min(contentLen-len(r.Body), len(currentData))
			r.Body = append(r.Body, currentData[:remaining]...)
			read += remaining

			if len(r.Body) == contentLen {
				r.state = StateDone
			}

		case StateDone:
			break outer

		default:
			panic("invalid parser state")
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

func (r *Request) hasBody() bool {
	contentLen, _ := r.Headers.GetInt("content-length", 0)
	return contentLen > 0
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
