package headers

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vaxxnsh/http-from-scratch/internal/utils"
)

type Headers struct {
	headers map[string]string
}

var rn = []byte("\r\n")
var MALFORMED_FIELD_LINE = fmt.Errorf("malformed field-line")
var MALFORMED_FIELD_NAME = fmt.Errorf("malformed field-name")
var INVALID_TOKEN_IN_FIELD_NAME = fmt.Errorf("invalid token in field name")

func NewHeaders() *Headers {
	return &Headers{
		headers: map[string]string{},
	}
}

func isToken(str string) bool {
	isValid := true

	for _, ch := range str {
		if !utils.IsAlphabet(ch) && !utils.IsNum(ch) && !utils.IsValidSpecial(ch) {
			fmt.Printf("invalid token found : %s\n", string(ch))
			isValid = false
			break
		}
	}

	return isValid
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", MALFORMED_FIELD_LINE
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", MALFORMED_FIELD_NAME
	}

	return string(name), string(value), nil
}

func (h *Headers) Get(name string) (string, bool) {
	v, ok := h.headers[strings.ToLower(name)]

	return v, ok
}
func (h *Headers) Set(name, value string) error {
	if !isToken(name) {
		return INVALID_TOKEN_IN_FIELD_NAME
	}

	name = strings.ToLower(name)

	prev, ok := h.Get(name)

	if ok {
		h.headers[name] = fmt.Sprintf("%s, %s", prev, value)
	} else {
		h.headers[name] = value
	}

	return nil
}

func (h *Headers) Parse(data []byte) (int, bool, error) {
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
		err = h.Set(name, value)

		if err != nil {
			return 0, false, err
		}
	}
	return read, done, nil
}
