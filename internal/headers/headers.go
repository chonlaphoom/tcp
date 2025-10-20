package headers

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if isContainsCRLF := bytes.Contains(data, []byte(CRLF)); !isContainsCRLF {
		// Not enough data to parse headers
		return 0, false, nil
	}

	if bytes.Index(data, []byte(CRLF)) == 0 {
		// Done parsing headers
		return 0, true, nil
	}

	if isEndsWithCRLF := bytes.HasSuffix(data, []byte(CRLF)); !isEndsWithCRLF {
		// Not enough data to parse headers
		return 0, false, nil
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
	if !containsOnlyAllowCharacters(headerName) {
		return 0, false, errors.New("invalid characters in header")
	}

	headerValue := string(bytes.TrimSpace(headerParts[1]))
	h.Set(strings.ToLower(headerName), headerValue)

	return len(part) + len(CRLF), false, nil
}

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Set(key, value string) {
	if h[key] != "" {
		value = h[key] + ", " + value
	}
	h[key] = value
}

func containsOnlyAllowCharacters(data string) bool {
	allowedChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$%&'*+-.^_`|~"
	pattern := "^[" + regexp.QuoteMeta(allowedChars) + "]+$"
	matched, _ := regexp.MatchString(pattern, data)
	return matched
}
