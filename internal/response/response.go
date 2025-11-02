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

func (w *Writer) ResetBuffer() error {
	return w.ResetBuffer()
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	lenInStr := fmt.Sprintf("%x\r\n", len(p))
	var buffer bytes.Buffer
	buffer.Write([]byte(lenInStr))
	buffer.Write(p)
	buffer.Write([]byte("\r\n"))

	err := w.WriteBody(buffer.Bytes())

	return buffer.Len(), err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	body := []byte("0\r\n\r\n")
	err := w.WriteBody(body)
	return len(body), err
}

func (w *Writer) WriteBodyTrailers(p []byte) (int, error) {
	var buffer bytes.Buffer
	buffer.Write(p)
	buffer.Write([]byte("\r\n"))
	w.WriteBody(buffer.Bytes())

	return buffer.Len(), nil
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

// function optional
// I don't use pointer here, since maps are already ref types
func NewResponseHeaders(headerCfgs ...ResponseHeaderCfg) headers.Headers {
	headers := headers.Headers{}
	for _, cfg := range headerCfgs {
		cfg(headers)
	}
	return headers
}

type ResponseHeaderCfg func(headers.Headers)

func NewContentType(cType string) ResponseHeaderCfg {
	return func(h headers.Headers) {
		h["Content-type"] = cType
	}
}
func NewContentLength(l int) ResponseHeaderCfg {
	return func(h headers.Headers) {
		h["Content-Length"] = fmt.Sprint(l)
	}
}
func NewTransferEncoding(ct string) ResponseHeaderCfg {
	return func(h headers.Headers) {
		h["Transfer-Encoding"] = ct
	}
}
func NewConnection(ct string) ResponseHeaderCfg {
	if ct == "" {
		ct = "close"
	}

	return func(h headers.Headers) {
		h["Connection"] = ct
	}
}
func NewTrailer(trailers []string) ResponseHeaderCfg {
	return func(h headers.Headers) {
		h[strings.ToLower("Trailer")] = strings.Join(trailers, ", ")
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
