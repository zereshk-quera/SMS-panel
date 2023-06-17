package handlers

import (
	"log"
	"strings"
	"time"

	"SMS-panel/models"
)

type Message struct {
	Text        string `json:"text"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func SendMessageHandler(message *Message) (string, error) {
	// Log the message details
	log.Printf("Message sent - Text: %s, Source: %s, Destination: %s", message.Text, message.Source, message.Destination)

	deliveryReport := "Message sent successfully"
	return deliveryReport, nil
}

func CreateSMSTemplate(template string, phoneNumber models.PhoneBookNumber) string {
	template = strings.ReplaceAll(template, "%name", phoneNumber.Name)
	template = strings.ReplaceAll(template, "%prefix", phoneNumber.Prefix)

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	template = strings.ReplaceAll(template, "%date", currentTime)
	template = strings.ReplaceAll(template, "%username", phoneNumber.Username)

	return template
}
