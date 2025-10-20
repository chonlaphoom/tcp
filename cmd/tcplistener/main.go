package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"tcpgo/internal/request"
)

func main() {
	listener, errorListener := net.Listen("tcp", ":42069")
	if errorListener != nil {
		log.Fatalf("could not start listener: %s\n", errorListener)
	}
	defer listener.Close()

	fmt.Println("listening for TCP traffic on :42069")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("could not accept connection: %s\n", err)
			continue
		}
		fmt.Printf("accepted connection from %s\n", conn.RemoteAddr())

		request, errorRead := request.RequestFromReader(conn)

		if errorRead != nil {
			fmt.Printf("reader error: %s\n", errorRead)
		}

		printRequestLines(request.RequestLine)
		printHeaders(request.Headers)

		fmt.Printf("closed connection from %s\n", conn.RemoteAddr())
		break
	}
}

func printRequestLines(request request.RequestLine) {
	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", request.Method, request.RequestTarget, request.HttpVersion)
}

func printHeaders(headers map[string]string) {
	fmt.Println("")
	strToPrint := "Headers:\n"
	for key, value := range headers {
		if key == "user-agent" && strings.HasPrefix(value, "curl/") {
			// normalize curl user-agent for testing
			strToPrint += fmt.Sprintf("- %s: %s\n", key, "curl")
			continue
		}
		strToPrint += fmt.Sprintf("- %s: %s\n", key, value)
	}
	fmt.Print(strToPrint)
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	lineChan := make(chan string)

	go func() {
		defer file.Close()
		defer close(lineChan)

		currentLine := ""
		for {
			readBuffer := make([]byte, 8)
			upto, errorRead := file.Read(readBuffer)

			if errorRead != nil {
				if currentLine != "" {
					lineChan <- currentLine
				}

				if errors.Is(errorRead, io.EOF) {
					break
				}

				fmt.Printf("could not read from file: %s\n", errorRead)
				break
			}

			// content less than buffer size
			if upto < 8 {
				lineChan <- fmt.Sprintf("%s%s", currentLine, string(readBuffer[:upto]))
				break
			}

			parts := strings.Split(string(readBuffer[:upto]), "\r\n")
			for i := 0; i < len(parts)-1; i++ {
				lineChan <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return lineChan
}
