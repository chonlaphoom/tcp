package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	const inputFilePath = "messages.txt"

	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	readBuffer := make([]byte, 8, 8)
	for {
		_, errorRead := file.Read(readBuffer)
		if errorRead != nil {
			if errors.Is(errorRead, io.EOF) {
				break
			}
			fmt.Printf("could not read from %s: %s\n", inputFilePath, errorRead)
			break
		}

		msg := fmt.Sprint("read: ", string(readBuffer))
		fmt.Println(msg)
	}

	os.Exit(0)
}
