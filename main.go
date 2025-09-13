package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	const inputFilePath = "messages.txt"

	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	c := getLinesChannel(file)

	for line := range c {
		fmt.Print(line)
	}

	os.Exit(0)
}

func getLinesChannel(file io.ReadCloser) <-chan string {
	c := make(chan string)

	go func() {
		defer file.Close()
		defer close(c)

		currentLine := ""
		for {
			readBuffer := make([]byte, 8)
			upto, errorRead := file.Read(readBuffer)
			if errorRead != nil {
				if errors.Is(errorRead, io.EOF) {
					if currentLine != "" {
						fmt.Print("read: ", currentLine)
						c <- fmt.Sprint("read: ", currentLine)
					}
					break
				}

				fmt.Printf("could not read from file: %s\n", errorRead)
				break
			}

			parts := strings.Split(string(readBuffer[:upto]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				c <- fmt.Sprintf("read: %s%s\n", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return c
}
