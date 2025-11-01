package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"tcpgo/internal/request"
	"tcpgo/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	Msg  string
	Code response.StatusCode
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.Code)
	response.WriteHeaders(w, response.GetDefaultHeaders(len(h.Msg)))
	w.Write([]byte(h.Msg))
	log.Printf("handler error: %s: %v\n", h.Msg, h.Code)
}

func Serve(port int, handler Handler) (*Server, error) {
	portStr := fmt.Sprintf(":%s", strconv.Itoa(port))
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: l,
		handler:  handler,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("could not accept connection: %s\n", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerError := &HandlerError{
			Msg:  fmt.Sprintf("could not parse request: %v", err),
			Code: response.StatusBadRequest,
		}
		handlerError.Write(conn)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	hdrErr := s.handler(buffer, req)
	if hdrErr != nil {
		hdrErr.Write(conn)
		return
	}

	// successful handling
	response.WriteStatusLine(conn, response.StatusOK)
	bufferedData := buffer.Bytes()

	responseHeader := response.GetDefaultHeaders(len(bufferedData))
	response.WriteHeaders(conn, responseHeader)

	var concatenatedBuffer bytes.Buffer
	concatenatedBuffer.Write(bufferedData)
	conn.Write(concatenatedBuffer.Bytes())
}
