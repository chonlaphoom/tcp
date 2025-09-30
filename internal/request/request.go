package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	parsedRequestLine, errParseRequestLin := parseRequestLine(buf)
	if errParseRequestLin != nil {
		return &Request{}, errParseRequestLin
	}

	return &Request{
		RequestLine: parsedRequestLine}, nil
}

func parseRequestLine(buf []byte) (RequestLine, error) {
	str := string(buf)
	firstLine := strings.Split(str, "\r\n")[0]
	parts := strings.Split(firstLine, " ")

	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("error parsing request line")
	}

	if parts[0] != "GET" && parts[0] != "POST" && parts[0] != "PUT" && parts[0] != "DELETE" && parts[0] != "HEAD" && parts[0] != "OPTIONS" && parts[0] != "PATCH" {
		return RequestLine{}, fmt.Errorf("invalid method in request line")
	}

	if parts[2] != "HTTP/1.1" {
		return RequestLine{}, fmt.Errorf("invalid HTTP version in request line")
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   strings.TrimPrefix(parts[2], "HTTP/"),
	}, nil
}
