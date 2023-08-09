package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type RequestData struct {
	FirstName  string `json:"FName"`
	LastName   string `json:"LName"`
	MiddleName string `json:"MName"`
	DomainName string `json:"DoName"`
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Check if the required fields are empty
	if requestData.FirstName == "" || requestData.LastName == "" || requestData.DomainName == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	fmt.Println("Received data:")
	fmt.Printf("FirstName %s\nMiddleName %s\nLastName %s\nDomainName %s\n", requestData.FirstName, requestData.MiddleName, requestData.LastName, requestData.DomainName)

	emailID := Telnet(requestData)

	if emailID != "" {
		// Prepare the response with the email ID
		response := map[string]string{
			"EmailID": emailID,
		}

		// Set the response headers and encode the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		// Return an error response if email ID doesn't exist
		http.Error(w, "Email ID not found", http.StatusNotFound)
	}
}

func Telnet(data RequestData) string {
	UserEmailaddress := toLowerCase(data.FirstName + "." + data.LastName + "@" + data.DomainName)
	FirstNameEmailaddress := toLowerCase(data.FirstName + "@" + data.DomainName)
	MiddleNameEmailaddress := toLowerCase(data.MiddleName + "@" + data.DomainName)
	LastNameEmailaddress := toLowerCase(data.LastName + "@" + data.DomainName)
	InitalNameEmailaddress := toLowerCase(data.FirstName + "." + string(data.LastName[0]) + "@" + data.DomainName)

	fmt.Println("Connected String:", UserEmailaddress)

	mxRecords, err := net.LookupMX(data.DomainName)
	if err != nil {
		fmt.Println("Error performing MX record lookup:", err)
		os.Exit(1)
	}

	if len(mxRecords) == 0 {
		fmt.Printf("No MX records found for %s\n", data.DomainName)
		return "" // Exit the function since there are no MX records
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
			return ""
		}
		defer conn.Close()

		fmt.Println("Connected to the server.")

		// Start a loop to read and write data from/to the server
		done := make(chan struct{})
		go readFromServer(conn, done, UserEmailaddress)

		// List of commands to send to the server
		commands := []string{
			"EHLO " + data.DomainName,
			"MAIL FROM: <example@example.com>",
			"RCPT TO: <" + UserEmailaddress + ">",
			"RCPT TO: <" + FirstNameEmailaddress + ">",
			"RCPT TO: <" + InitalNameEmailaddress + ">",
			"RCPT TO: <" + MiddleNameEmailaddress + ">",
			"RCPT TO: <" + LastNameEmailaddress + ">",
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
			time.Sleep(time.Millisecond * 1800)
		}

		fmt.Println("Connection closed.")
	}

	return UserEmailaddress
}

func readFromServer(conn net.Conn, done chan struct{}, userEmail string) string {
	reader := bufio.NewReader(conn)

	// Initialize a flag to indicate if the response is received
	emailStatusReceived := false

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}

		emailStatusReceived = true

		if emailStatusReceived {
			// We have received the email status response, now handle it
			if strings.Contains(message, "550-5.1.1 ") {
				fmt.Println("Email doesn't exist")
				return "" // Return an empty string if email doesn't exist
			} else if strings.Contains(message, "250 2.1.5 OK") {
				return userEmail // Return the email ID if it exists
			}
		}
	}

	// Signal that the goroutine has completed its task
	done <- struct{}{}
	return "" // Return an empty string if email status is not received
}

func toLowerCase(str string) string {
	return strings.ToLower(str)
}

func main() {
	http.HandleFunc("/email", handlePostRequest)

	port := "8080" // Choose a port number for your API
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
