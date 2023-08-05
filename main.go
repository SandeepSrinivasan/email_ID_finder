package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestData struct {
	FirstName  string `json:"FName"`
	LastName   string `json:"LName"`
	MiddleName string `json:"MName"`
	DomainName string `json:"DoName"`
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	fmt.Println("Received data:")
	fmt.Printf("FirstName %s\nMiddleName %s\nLastName %s\nDomainName %s\n", requestData.FirstName, requestData.MiddleName, requestData.LastName, requestData.DomainName)

	Telnet(requestData)

	// Send a JSON response
	response := map[string]string{
		"FirstName": requestData.FirstName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/email", handlePostRequest)

	port := "8080" // Choose a port number for your API
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
