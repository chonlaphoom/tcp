package response

import (
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

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Type":   "text/plain",
		"Content-Length": fmt.Sprint(contentLen),
		"Connection":     "close",
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
