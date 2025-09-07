package main

import (
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("messages.txt")
	defer file.Close()

	readBuffer := make([]byte, 8)
	for {
		_, errorRead := file.Read(readBuffer)
		if errorRead != nil {
			break
		}

		msg := fmt.Sprint("read: ", string(readBuffer))
		fmt.Println(msg)
	}

	os.Exit(0)
}
