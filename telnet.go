// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"os"
// 	"strings"
// 	"time"
// )

// func Telnet(data RequestData) {

// 	// fmt.Println("enter the domain name:")
// 	// //domain := "fampay.in"
// 	// domain := data.FirstName

// 	// fmt.Println("enter the First name:")
// 	// //	FirstName := "Sandeep"
// 	// FirstName := readInput()

// 	// fmt.Println("enter the middle name:")
// 	// //	FirstName := "Sandeep"
// 	// MiddleName := readInput()

// 	// fmt.Println("enter the last name:")
// 	// // LastName := "Srinivasan"
// 	// LastName := readInput()

// 	// var data RequestData

// 	UserEmailaddress := toLowerCase(data.FirstName + "." + data.LastName + "@" + data.DomainName)
// 	FirstNameEmailaddress := toLowerCase(data.FirstName + "@" + data.DomainName)
// 	MiddleNameEmailaddress := toLowerCase(data.MiddleName + "@" + data.DomainName)
// 	LastNameEmailaddress := toLowerCase(data.LastName + "@" + data.DomainName)
// 	InitalNameEmailaddress := toLowerCase(data.FirstName + "." + string(data.LastName[0]) + "@" + data.DomainName)

// 	fmt.Println("Connected String:", UserEmailaddress)

// 	mxRecords, err := net.LookupMX(data.DomainName)
// 	if err != nil {
// 		fmt.Println("Error performing MX record lookup:", err)
// 		os.Exit(1)
// 	}

// 	if len(mxRecords) == 0 {
// 		fmt.Printf("No MX records found for %s\n", data.DomainName)
// 	} else {
// 		// Get the host of the first MX record
// 		mailExchange := mxRecords[0].Host
// 		fmt.Printf("First MX record for %s\n", mailExchange)

// 		// Combine the server IP and port to create the address
// 		serverAddress := mailExchange + ":" + "25"

// 		// Connect to the server
// 		conn, err := net.Dial("tcp", serverAddress)
// 		if err != nil {
// 			fmt.Println("Error connecting to the server:", err)
// 			return
// 		}
// 		defer conn.Close()

// 		fmt.Println("Connected to the server.")

// 		// Start a loop to read and write data from/to the server
// 		done := make(chan struct{})
// 		go readFromServer(conn, done)

// 		// List of commands to send to the server
// 		commands := []string{
// 			"EHLO " + data.DomainName,
// 			"MAIL FROM: <example@example.com>",
// 			"RCPT TO: <" + UserEmailaddress + ">",
// 			"RCPT TO: <" + FirstNameEmailaddress + ">",
// 			"RCPT TO: <" + InitalNameEmailaddress + ">",
// 			"RCPT TO: <" + MiddleNameEmailaddress + ">",
// 			"RCPT TO: <" + LastNameEmailaddress + ">",
// 		}

// 		for _, command := range commands {
// 			// Print the selected command
// 			fmt.Printf("Selected command: %s\n", command)

// 			// Check if the user wants to exit
// 			if command == "exit" {
// 				break
// 			}

// 			// Send the command to the server
// 			_, err := fmt.Fprintf(conn, "%s\n", command)
// 			if err != nil {
// 				fmt.Println("Error sending command:", err)
// 				break
// 			}

// 			// Wait for a short time to allow the server to respond
// 			// (adjust the duration as needed)
// 			time.Sleep(time.Millisecond * 1800)
// 		}

// 		fmt.Println("Connection closed.")
// 	}
// }

// func readFromServer(conn net.Conn, done chan struct{}) {
// 	reader := bufio.NewReader(conn)

// 	// Initialize a flag to indicate if the response is received
// 	emailStatusReceived := false

// 	for {
// 		message, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println("Error reading from server:", err)
// 			break
// 		}

// 		emailStatusReceived = true

// 		if emailStatusReceived {
// 			// We have received the email status response, now handle it
// 			if strings.Contains(message, "550-5.1.1 ") {
// 				fmt.Println("Email doesn't exist")
// 			} else if strings.Contains(message, "250 2.1.5 OK") {
// 				fmt.Println("Email exists")
// 			}
// 		}
// 	}

// 	// Signal that the goroutine has completed its task
// 	done <- struct{}{}
// }

// func readInput() string {
// 	reader := bufio.NewReader(os.Stdin)
// 	input, _ := reader.ReadString('\n')
// 	return input[:len(input)-1] // Remove the newline character
// }

// func toLowerCase(str string) string {
// 	return strings.ToLower(str)
// }

// func processAPIInput(data RequestData) {
// 	fmt.Println("First Name:", data.FirstName)
// 	fmt.Println("Middle Name:", data.FirstName)
// 	fmt.Println("Last Name:", data.LastName)
// 	fmt.Println("Domain Name:", data.DomainName)
// }
