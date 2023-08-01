package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <domain_name>")
		os.Exit(1)
	}

	domain := os.Args[1]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		fmt.Println("Error performing MX record lookup:", err)
		os.Exit(1)
	}

	if len(mxRecords) == 0 {
		fmt.Printf("No MX records found for %s\n", domain)
	} else {
		// Get the host of the first MX record
		mailExchange := mxRecords[0].Host
		fmt.Printf("First MX record for %s\n", mailExchange)

		// Combine the server IP and port to create the address
		serverAddress := mailExchange + ":" + "25"

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

		// List of commands to send to the server
		commands := []string{
			"EHLO " + domain,
			"MAIL FROM: <example@example.com>",
			"RCPT TO: <1@in.in>",
			// Add more commands here as needed
		}

		for _, command := range commands {
			// Print the selected command
			fmt.Printf("Selected command: %s\n", command)

			// Check if the user wants to exit
			if command == "exit" {
				break
			}

			// Send the command to the server
			_, err := fmt.Fprintf(conn, "%s\n", command)
			if err != nil {
				fmt.Println("Error sending command:", err)
				break
			}

			// Wait for a short time to allow the server to respond
			// (adjust the duration as needed)
			time.Sleep(time.Millisecond * 800)
		}

		fmt.Println("Connection closed.")
	}
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
