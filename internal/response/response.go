package response

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"tcpgo/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	Buffer      *bytes.Buffer
	WriterState int // 0 good, 1 bad
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	err := WriteStatusLine(w.Buffer, statusCode)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	err := WriteHeaders(w.Buffer, headers)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(body []byte) error {
	_, err := w.Buffer.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var HTTPVersion = "1.1"
	var Protocol = "HTTP"
	switch statusCode {
	case StatusOK:
		statusLine := fmt.Sprint(Protocol, "/", HTTPVersion, " ", StatusOK, " ", "OK", "\r\n")
		_, err := w.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	case StatusBadRequest:
		statusLine := fmt.Sprint(Protocol, "/", HTTPVersion, " ", StatusBadRequest, " ", "Bad Request", "\r\n")
		_, err := w.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	case StatusInternalServerError:
		statusLine := fmt.Sprint(Protocol, "/", HTTPVersion, " ", StatusInternalServerError, " ", "Internal Server Error", "\r\n")
		_, err := w.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	default:
		statusLine := fmt.Sprint(Protocol, "/", HTTPVersion, " ", statusCode, " ", "\r\n")
		_, err := w.Write([]byte(statusLine))
		if err != nil {
			return err
		}
		return nil
	}
}

func NewResponseHeaders(ctnlen int, ctntype string, connectiontype string) headers.Headers {
	return headers.Headers{
		"Content-Type":   ctntype,
		"Content-Length": fmt.Sprint(ctnlen),
		"Connection":     connectiontype,
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	headerList := make([]string, 0, len(headers))
	for key, value := range headers {
		headerList = append(headerList, fmt.Sprintf("%s: %s", key, value))
	}

	_, err := w.Write([]byte(strings.Join(headerList, "\r\n") + "\r\n\r\n"))
	return err
}
