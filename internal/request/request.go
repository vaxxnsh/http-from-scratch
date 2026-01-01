package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var MALFORMED_START_LINE = fmt.Errorf("malformed start-line")
var UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var SEPARATOR = "\r\n"

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func parseRequestLine(s string) (*RequestLine, string, error) {
	before, after, ok := strings.Cut(s, SEPARATOR)
	if !ok {
		return nil, s, nil
	}

	startLine := before
	restOfMsg := after

	parts := strings.Split(startLine, " ")

	if len(parts) != 3 {
		return nil, restOfMsg, MALFORMED_START_LINE
	}

	r := RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	if !r.ValidHTTP() {
		return nil, restOfMsg, UNSUPPORTED_HTTP_VERSION
	}

	r.HttpVersion = strings.Split(parts[2], "/")[1]

	return &r, restOfMsg, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadAll"), err)
	}

	str := string(data)

	rl, _, err := parseRequestLine(str)

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, nil
}
