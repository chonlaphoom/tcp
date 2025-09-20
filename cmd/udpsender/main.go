package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	n, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	conn, errorDialUDP := net.DialUDP("udp", nil, n)
	if errorDialUDP != nil {
		log.Fatal(errorDialUDP)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			os.Exit(1)
		}
		fmt.Println("Message sent:", message)
	}
}
