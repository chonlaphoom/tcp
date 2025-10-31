package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"tcpgo/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	portStr := fmt.Sprintf(":%s", strconv.Itoa(port))
	l, err := net.Listen("tcp", portStr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: l,
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
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
}
