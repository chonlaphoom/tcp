package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tcpgo/internal/request"
	"tcpgo/internal/response"
	"tcpgo/internal/server"
)

const port = 42069

func handler(w response.Writer, req *request.Request) {
	contentType := "text/html"
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		w.WriteStatusLine(response.StatusBadRequest)
		body := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		w.WriteHeaders(response.NewResponseHeaders(response.NewContentLength(len(body)), response.NewContentType(contentType), response.NewConnection("")))
		w.WriteBody([]byte(body))
		return
	}

	if after, ok := strings.CutPrefix(target, "/httpbin/stream/"); ok {
		count := after
		res, err := http.Get("https://httpbin.org/stream/" + count)
		if err != nil {
			w.WriteStatusLine(response.StatusInternalServerError)
			w.WriteHeaders(response.NewResponseHeaders(response.NewContentLength(0), response.NewContentType(contentType), response.NewConnection("")))
			return
		}

		w.WriteStatusLine(response.StatusOK)
		w.WriteHeaders(response.NewResponseHeaders(response.NewContentType("text/plain"), response.NewTransferEncoding("chunked")))
		//
		buffer := make([]byte, 1024)
		defer res.Body.Close()
		for {
			n, err := res.Body.Read(buffer)
			if n > 0 {
				w.WriteChunkedBody(buffer[:n])
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				w.ResetBuffer()
				w.WriteStatusLine(response.StatusInternalServerError)
				w.WriteHeaders(response.NewResponseHeaders(response.NewContentLength(0), response.NewContentType(contentType), response.NewConnection(""), response.NewConnection("")))
				return
			}
		}

		w.WriteChunkedBodyDone()
		log.Printf("%s", w.Buffer.Bytes())
		return
	}

	if target == "/myproblem" {
		w.WriteStatusLine(response.StatusInternalServerError)
		body := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		w.WriteHeaders(response.NewResponseHeaders(response.NewContentLength(len(body)), response.NewContentType(contentType), response.NewConnection("")))
		w.WriteBody([]byte(body))
		return
	}

	w.WriteStatusLine(response.StatusOK)
	body := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	w.WriteHeaders(response.NewResponseHeaders(response.NewContentLength(len(body)), response.NewContentType(contentType), response.NewConnection("")))
	w.WriteBody([]byte(body))
}

func main() {
	server, err := server.Serve(port, handler)
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
