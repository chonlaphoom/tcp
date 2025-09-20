package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"net"
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

		linesChannel := getLinesChannel(conn)

		for line := range linesChannel {
			fmt.Println(line)
		}

		fmt.Printf("closed connection from %s\n", conn.RemoteAddr())
	}
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
				if errors.Is(errorRead, io.EOF) {
					if currentLine != "" {
						lineChan <- fmt.Sprint(currentLine)
					}
					break
				}

				fmt.Printf("could not read from file: %s\n", errorRead)
				break
			}

			parts := strings.Split(string(readBuffer[:upto]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				lineChan <- fmt.Sprintf("%s%s\n", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return lineChan
}
