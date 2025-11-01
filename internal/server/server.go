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

type Handler func(w response.Writer, req *request.Request)
type Server struct {
	listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

type HandlerError struct {
	Msg  string
	Code response.StatusCode
}

func (h *HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, h.Code)
	headers := response.NewResponseHeaders(len(h.Msg), "text/plain", "close")
	response.WriteHeaders(w, headers)
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
	res := response.Writer{
		Buffer:      buffer,
		WriterState: 0,
	}
	s.handler(res, req)
	conn.Write(res.Buffer.Bytes())
}
