package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Get the server IP and port from user input
	fmt.Print("Enter the server IP address: ")
	serverIP := readInput()

	// Combine the server IP and port to create the address
	serverAddress := serverIP + ":" + "25"

	// Connect to the server
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the server.")

	// Start a loop to read and write data from/to the server
	go readFromServer(conn)

	for {
		// Read input from the user
		fmt.Print("Enter your message (type 'exit' to quit): ")
		message := readInput()

		// Check if the user wants to exit
		if message == "exit" {
			break
		}

		// Send the user's message to the server
		_, err := fmt.Fprintf(conn, "%s\n", message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	fmt.Println("Connection closed.")
}

func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}
		fmt.Println("Server:", message)
	}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return input[:len(input)-1] // Remove the newline character
}
