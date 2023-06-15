package handlers

import (
	"log"
)

type Message struct {
	Text        string `json:"text"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// SendMessageHandler simulates sending an SMS message and returns the delivery report
func SendMessageHandler(message *Message) (string, error) {
	// Log the message details
	log.Printf("Message sent - Text: %s, Source: %s, Destination: %s", message.Text, message.Source, message.Destination)

	// Simulate a successful response
	deliveryReport := "Message sent successfully"
	return deliveryReport, nil
}
