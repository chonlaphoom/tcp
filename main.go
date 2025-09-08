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
	defer file.Close()

	var currentLine string
	for {
		readBuffer := make([]byte, 8, 8)
		upto, errorRead := file.Read(readBuffer)
		if errorRead != nil {
			if errors.Is(errorRead, io.EOF) {
				if currentLine != "" {
					fmt.Print("read: ", currentLine)
					currentLine = ""
				}
				break
			}
			fmt.Printf("could not read from %s: %s\n", inputFilePath, errorRead)
			break
		}

		parts := strings.Split(string(readBuffer[:upto]), "\n")
		for i := 0; i < len(parts)-1; i++ {
			fmt.Printf("read: %s%s\n", currentLine, parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}

	os.Exit(0)
}
