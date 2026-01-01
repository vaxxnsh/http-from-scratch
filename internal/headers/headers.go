package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")
var MALFORMED_Field_LINE = fmt.Errorf("malformed field-line")
var MALFORMED_Field_NAME = fmt.Errorf("malformed field-name")

func NewHeaders() Headers {
	return map[string]string{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", MALFORMED_Field_LINE
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", MALFORMED_Field_NAME
	}

	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		read += idx + len(rn)
		if err != nil {
			return 0, false, err
		}
		h[name] = value
	}
	return read, done, nil
}
