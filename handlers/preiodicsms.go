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
	Username    string `json:"username"`
	Phone       string `json:"phone"`
	Message     string `json:"message"`
	Schedule    string `json:"schedule"`
	Interval    int    `json:"interval"`
	PhoneBookID string `json:"phone_book_id"`
}

func PeriodicSendSMSHandler(c echo.Context) error {
	account := c.Get("account").(models.Account)
	var request SendSMSRequestPeriodic
	if err := c.Bind(&request); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request payload")
	}

	// Parse the schedule time
	scheduleTime, err := parseTime(request.Schedule)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid schedule time format")
	}

	if request.Phone == "" && request.Username == "" && request.PhoneBookID == "" {
		return c.String(http.StatusBadRequest, "Recipient not provided")
	}

	var phoneBookNumbers []models.PhoneBookNumber
	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %s", err.Error())
		return fmt.Errorf("database issue")
	}
	phoneNumberQuery := db.Joins("JOIN phone_books ON phone_books.id = phone_book_numbers.phone_book_id").
		Where("phone_books.account_id = ?", account.ID)

	if request.Phone != "" {
		phoneNumberQuery = phoneNumberQuery.Where("phone_book_numbers.phone = ?", request.Phone)
	} else if request.Username != "" {
		phoneNumberQuery = phoneNumberQuery.Where("phone_book_numbers.username = ?", request.Username)
	} else if request.PhoneBookID != "" {
		phoneNumberQuery = phoneNumberQuery.Where("phone_book_numbers.phone_book_id = ?", request.PhoneBookID)
	}

	if err := phoneNumberQuery.Find(&phoneBookNumbers).Error; err != nil {
		return c.String(http.StatusBadRequest, "Recipient does not exist in the phone book")
	}

	smsCount := len(phoneBookNumbers)
	reduceAccountBudget(account, smsCount)

	for _, phoneBookNumber := range phoneBookNumbers {
		sms := &models.SMSMessage{
			Sender:    account.Username,
			Recipient: phoneBookNumber.Phone,
			Message:   request.Message,
			AccountID: account.ID,
			Schedule:  &scheduleTime,
		}

		intervalDuration := time.Duration(request.Interval) * time.Second

		go sendSMSWithRepetition(sms, intervalDuration)
	}

	return c.String(http.StatusOK, "SMS scheduled successfully")
}

func parseTime(schedule string) (time.Time, error) {
	parts := strings.Split(schedule, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid schedule time format")
	}

	hour := parts[0]
	minute := parts[1]

	now := time.Now().UTC()

	scheduleTime, err := time.Parse("2006-01-02 15:04", now.Format("2006-01-02")+" "+hour+":"+minute)
	if err != nil {
		return time.Time{}, err
	}

	if scheduleTime.Before(now) {
		scheduleTime = scheduleTime.AddDate(0, 0, 1)
	}

	return scheduleTime, nil
}

func sendSMSWithRepetition(sms *models.SMSMessage, interval time.Duration) {
	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %s", err.Error())
		return
	}

	for {
		now := time.Now()

		remainingTime := time.Until(*sms.Schedule)
		if remainingTime > 0 {
			time.Sleep(remainingTime)
		}

		deliveryReport, err := MockSendMessage(&Message{
			Text:        sms.Message,
			Source:      sms.Sender,
			Destination: sms.Recipient,
		})
		if err != nil {
			log.Printf("Failed to send SMS: %s", err.Error())
		} else {
			sms.DeliveryReport = deliveryReport
			err = db.Save(sms).Error
			if err != nil {
				log.Printf("Failed to update SMS delivery report: %s", err.Error())
			}
		}

		nextSchedule := now.Add(interval)

		for nextSchedule.Before(time.Now()) {
			nextSchedule = nextSchedule.Add(interval)
		}

		sms.Schedule = &nextSchedule
	}
}

func reduceAccountBudget(account models.Account, smsCount int) {
	db, err := database.GetConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %s", err.Error())
		return
	}

	var smsCost int
	if smsCount == 1 {
		if err := db.Table("configuration").Where("name = ?", "single sms").Select("value").Scan(&smsCost).Error; err != nil {
			log.Printf("Failed to retrieve single SMS cost: %s", err.Error())
			return
		}
	} else {
		if err := db.Table("configuration").Where("name = ?", "group sms").Select("value").Scan(&smsCost).Error; err != nil {
			log.Printf("Failed to retrieve group SMS cost: %s", err.Error())
			return
		}
	}

	totalCost := smsCost * smsCount

	account.Budget -= int64(totalCost)

	if err := db.Save(&account).Error; err != nil {
		log.Printf("Failed to update account's budget: %s", err.Error())
	}
}
