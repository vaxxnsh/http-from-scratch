package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/vaxxnsh/http-from-scratch/internal/request"
	"github.com/vaxxnsh/http-from-scratch/internal/response"
	"github.com/vaxxnsh/http-from-scratch/internal/server"
)

const port = 42069

var ERROR_WHILE_GETTING_CHUNKED_BODY = fmt.Errorf("error while getting chunking body")

func getHtmlBodyForCode(statusCode response.StatusCode) []byte {
	switch statusCode {
	case response.StatusOk:
		return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
	case response.StatusBadRequest:
		return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
	case response.StatusInternalServerError:
		return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
	default:
		return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>unknow statusCode given but anything goes here.</p>
  </body>
</html>`)
	}
}

func respondWithhtml(w *response.Writer, statusCode response.StatusCode, htmlBody []byte) {
	w.WriteStatusLine(statusCode)
	h := response.GetDefaultHeaders(len(htmlBody))
	w.WriteHeaders(&h)
	w.WriteBody(htmlBody)
}

func sendCunkedResponse(w *response.Writer, numResponses int) error {
	res, err := http.Get(fmt.Sprintf("https://httpbin.org/stream/%d", numResponses))
	if err != nil {
		return ERROR_WHILE_GETTING_CHUNKED_BODY
	}
	h := response.GetChunkedHeaders()
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(&h)
	for {
		data := make([]byte, 32)
		n, err := res.Body.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		_, err = w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
		if err != nil {
			return err
		}
		_, err = w.WriteBody(data[:n])
		if err != nil {
			return err
		}
		_, err = w.WriteBody([]byte("\r\n"))
		if err != nil {
			return err
		}
	}
	_, err = w.WriteBody([]byte("0\r\n\r\n"))
	return err
}

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		target := req.RequestLine.RequestTarget
		switch target {
		case "/yourproblem":
			respondWithhtml(w, response.StatusBadRequest, getHtmlBodyForCode(response.StatusBadRequest))

		case "/myproblem":
			respondWithhtml(w, response.StatusBadRequest, getHtmlBodyForCode(response.StatusInternalServerError))
		default:
			binTarget := "/httpbin/stream/"
			videoTarget := "/prachi"
			if strings.HasPrefix(target, binTarget) {
				numResp, err := strconv.Atoi(target[len(binTarget):])
				if err != nil {
					fmt.Println(target[len(binTarget):])
					respondWithhtml(w, response.StatusBadRequest, getHtmlBodyForCode(response.StatusBadRequest))
					return
				}
				err = sendCunkedResponse(w, numResp)
				if err != nil {
					if errors.Is(err, ERROR_WHILE_GETTING_CHUNKED_BODY) {
						respondWithhtml(w, response.StatusInternalServerError, getHtmlBodyForCode(response.StatusInternalServerError))
					} else {
						fmt.Printf("error while reading body: %s\n", err)
					}
				}
			} else if strings.HasPrefix(target, videoTarget) {
				w.WriteStatusLine(response.StatusOk)
				f, err := os.ReadFile("./assets/umm...mp4")
				if err != nil {
					respondWithhtml(w, response.StatusInternalServerError, getHtmlBodyForCode(response.StatusInternalServerError))
					return
				}
				h := response.GetDefaultHeaders(len(f[:]))
				h.Replace("content-type", "video/mp4")
				w.WriteHeaders(&h)
				w.WriteBody(f[:])
			} else {
				respondWithhtml(w, response.StatusOk, getHtmlBodyForCode(response.StatusOk))
			}
		}
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
