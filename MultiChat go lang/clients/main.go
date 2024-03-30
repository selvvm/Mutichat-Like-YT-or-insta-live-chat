package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	// Server address
	serverAddr := "http://localhost:8080"

	// Messages to send
	messages := []string{"Hello from Client 1!", "Greetings from Client 2!", "Hey there from Client 3!"}

	// Iterate over messages and send them using different clients
	for i, msg := range messages {
		clientNum := i + 1
		clientMsg := []byte(msg)
		url := fmt.Sprintf("%s/client%d", serverAddr, clientNum)
		sendRequest(url, clientMsg)
	}
}

// Function to send HTTP request
func sendRequest(url string, body []byte) {
	// Send POST request with the message
	resp, err := http.Post(url, "text/plain", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Print response status
	fmt.Printf("Client sent message to %s. Response Status: %s\n", url, resp.Status)
}
