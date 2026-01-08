package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vaxxnsh/http-from-scratch/internal/request"
	"github.com/vaxxnsh/http-from-scratch/internal/response"
	"github.com/vaxxnsh/http-from-scratch/internal/server"
)

const port = 42069

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

func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			respondWithhtml(w, response.StatusBadRequest, getHtmlBodyForCode(response.StatusBadRequest))

		case "/myproblem":
			respondWithhtml(w, response.StatusBadRequest, getHtmlBodyForCode(response.StatusInternalServerError))
		default:
			respondWithhtml(w, response.StatusOk, getHtmlBodyForCode(response.StatusOk))
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
