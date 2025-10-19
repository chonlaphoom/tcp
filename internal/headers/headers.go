package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	isContainsCRLF := bytes.Contains(data, []byte(CRLF))
	if !isContainsCRLF {
		// Not enough data to parse headers
		return 0, false, nil
	}

	if bytes.Index(data, []byte(CRLF)) == 0 {
		// Done parsing headers
		return 0, true, nil
	}

	part := bytes.Split(data, []byte(CRLF))[0]

	if !bytes.Contains(part, []byte(":")) {
		return 0, false, errors.New("invalid header format")
	}

	headerParts := bytes.SplitN(part, []byte(":"), 2)

	if len(headerParts) != 2 {
		return 0, false, errors.New("invalid header format")
	}
	if bytes.HasSuffix(headerParts[0], []byte(" ")) {
		return 0, false, errors.New("invalid header format: unexpected whitespace in header name")
	}

	headerName := string(bytes.TrimLeft(headerParts[0], " "))
	headerValue := string(bytes.TrimSpace(headerParts[1]))
	h.Set(headerName, headerValue)

	return len(part) + len(CRLF), false, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
