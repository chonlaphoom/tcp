package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"tcpgo/internal/headers"
)

const BUFFER_SIZE = 8
const CRLF = "\r\n"

const (
	requestStateInitialized = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	state          int
	bodyLengthRead int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BUFFER_SIZE)
	request := &Request{
		state:       requestStateInitialized,
		RequestLine: RequestLine{},
		Headers:     headers.Headers{},
	}

	readerToIndex := 0
	for request.state != requestStateDone {
		// stretch the buffer if needed
		if readerToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numOfBytesRead, errRead := reader.Read(buf[readerToIndex:])

		if errRead != nil {
			if errors.Is(errRead, io.EOF) {
				request.state = requestStateDone
				break
			}
			return nil, errRead
		}

		readerToIndex += numOfBytesRead

		numBytesParsed, errParse := request.parse(buf[:readerToIndex])
		if errParse != nil {
			return nil, errParse
		}

		copy(buf, buf[numBytesParsed:])
		readerToIndex -= numBytesParsed
	}

	return &Request{
		RequestLine: request.RequestLine,
		Headers:     request.Headers,
		Body:        request.Body,
		state:       request.state,
	}, nil
}

func parseRequestLine(buf []byte) (*RequestLine, int, error) {
	idx := bytes.Index(buf, []byte(CRLF))
	if idx == -1 {
		// not enough data to parse the request line
		return nil, 0, nil
	}
	requestLineText := string(buf[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	// +2 to account for the \r\n
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	firstLine := strings.Split(str, CRLF)[0]
	parts := strings.Split(firstLine, " ")

	if len(parts) != 3 {
		return &RequestLine{}, fmt.Errorf("error parsing request line")
	}

	if parts[0] != "GET" && parts[0] != "POST" && parts[0] != "PUT" && parts[0] != "DELETE" && parts[0] != "HEAD" && parts[0] != "OPTIONS" && parts[0] != "PATCH" {
		return &RequestLine{}, fmt.Errorf("invalid method in request line")
	}

	if parts[2] != "HTTP/1.1" {
		return &RequestLine{}, fmt.Errorf("invalid HTTP version in request line")
	}

	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// need more data
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		if r.Headers == nil {
			r.Headers = headers.NewHeaders()
		}
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 && !done {
			// need more data
			return 0, nil
		}
		if done {
			r.state = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		contentLengthStr, exists := r.Headers.Get("content-length")
		if !exists {
			// no body
			r.state = requestStateDone
			return 0, nil
		}
		contentLengthInt, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length header")
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > contentLengthInt {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.bodyLengthRead == contentLengthInt {
			r.state = requestStateDone
		}
		return len(data), nil
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
