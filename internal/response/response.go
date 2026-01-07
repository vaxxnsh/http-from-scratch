package response

import (
	"fmt"
	"io"

	"github.com/vaxxnsh/http-from-scratch/internal/headers"
)

type Response struct {
}

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func formatStatusLine(statuscode StatusCode, reasonPhrase string) []byte {
	return fmt.Appendf(nil, "HTTP/1.1 %d %s\r\n", statuscode, reasonPhrase)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	reasonPhrase := ""
	switch statusCode {
	case StatusOk:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		return fmt.Errorf("invalid or not supported status code")
	}
	_, err := w.Write(formatStatusLine(statusCode, reasonPhrase))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return *h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	b := []byte{}
	headers.ForEach(func(n, v string) {
		b = fmt.Appendf(b, "%s: %s\r\n", n, v)
	})
	b = fmt.Appendf(b, "\r\n")
	_, err := w.Write(b)
	return err
}
