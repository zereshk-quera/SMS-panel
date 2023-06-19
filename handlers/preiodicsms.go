package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	database "SMS-panel/database"
	"SMS-panel/models"

	"github.com/labstack/echo/v4"
)

type SendSMSRequestPeriodic struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Schedule  string `json:"schedule"`
	Interval  int    `json:"interval"`
}

func PeriodicSendSMSHandler(c echo.Context) error {
	account := c.Get("account").(models.Account)
	// Parse the request body to get the sender, recipient, message, schedule time, and interval
	var request SendSMSRequestPeriodic
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	// Parse the schedule time
	scheduleTime, err := parseTime(request.Schedule)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid schedule time format")
	}

	// Create a new SMSMessage object
	sms := &models.SMSMessage{
		Sender:    request.Sender,
		Recipient: request.Recipient,
		Message:   request.Message,
		AccountID: account.ID,
		Schedule:  &scheduleTime,
	}

	// Calculate the interval duration based on the interval seconds
	intervalDuration := time.Duration(request.Interval) * time.Second

	// Call the function to send the SMS with repetition
	go sendSMSWithRepetition(sms, intervalDuration)

	return c.String(http.StatusOK, "SMS scheduled successfully")
}

func parseTime(schedule string) (time.Time, error) {
	// Split the schedule string into hours and minutes
	parts := strings.Split(schedule, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid schedule time format")
	}

	// Get the hours and minutes
	hour := parts[0]
	minute := parts[1]

	// Get the current date and time in UTC
	now := time.Now().UTC()

	// Create a new time object using the current date, parsed hours, and minutes in UTC
	scheduleTime, err := time.Parse("2006-01-02 15:04", now.Format("2006-01-02")+" "+hour+":"+minute)
	if err != nil {
		return time.Time{}, err
	}

	// Add 24 hours to the parsed time if it's in the past and the parsed time is before the current time
	if scheduleTime.Before(now) {
		// Increment the date by 1 day
		scheduleTime = scheduleTime.AddDate(0, 0, 1)
	}

	return scheduleTime, nil
}

func sendSMSWithRepetition(sms *models.SMSMessage, interval time.Duration) {
	// Get the database connection
	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %s", err.Error())
		return
	}

	for {
		// Get the current time
		now := time.Now()

		// Calculate the remaining time until the next SMS
		remainingTime := time.Until(*sms.Schedule)
		if remainingTime > 0 {
			// Sleep until the next SMS schedule
			time.Sleep(remainingTime)
		}

		// Call the function to send the actual SMS using your desired implementation
		deliveryReport, err := MockSendMessage(&Message{
			Text:        sms.Message,
			Source:      sms.Sender,
			Destination: sms.Recipient,
		})
		if err != nil {
			log.Printf("Failed to send SMS: %s", err.Error())
		} else {
			// Update the delivery report in the SMS message
			sms.DeliveryReport = deliveryReport
			err = db.Save(sms).Error
			if err != nil {
				log.Printf("Failed to update SMS delivery report: %s", err.Error())
			}
		}

		// Calculate the next schedule time based on the interval
		nextSchedule := now.Add(interval)

		// If the next schedule time is in the past, adjust it to the future based on the interval
		for nextSchedule.Before(time.Now()) {
			nextSchedule = nextSchedule.Add(interval)
		}

		// Set the new schedule time
		sms.Schedule = &nextSchedule
	}
}
