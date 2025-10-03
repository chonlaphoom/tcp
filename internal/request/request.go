package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8
const crlf = "\r\n"

type requestState int

const (
	requestStateInitialized requestState = 1
	requestStateDone        requestState = 0
)

type Request struct {
	RequestLine RequestLine
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	request := &Request{
		state: requestStateInitialized,
	}

	readerToIndex := 0
	err := error(nil)
	for request.state != requestStateDone {
		// stretch the buffer if needed
		if readerToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// try to parse what we have
		numOfBytesRead, errRead := reader.Read(buf[readerToIndex:])

		if errRead != nil {
			if errors.Is(errRead, io.EOF) {
				_, errParse := request.parse(buf[:readerToIndex])
				if errParse != nil {
					return nil, errParse
				}
				request.state = requestStateDone
				break
			}
			return nil, errRead
		}

		readerToIndex += numOfBytesRead
	}

	if err != nil {
		return request, err
	}

	return &Request{
		RequestLine: request.RequestLine,
	}, nil
}

func parseRequestLine(buf []byte) (*RequestLine, int, error) {
	idx := bytes.Index(buf, []byte(crlf))
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
	firstLine := strings.Split(str, "\r\n")[0]
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
		r.state = requestStateDone
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
