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
	"gorm.io/gorm"
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

	scheduleTime, err := parseTime(request.Schedule)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid schedule time format")
	}

	if request.Phone == "" && request.Username == "" && request.PhoneBookID == "" {
		return c.String(http.StatusBadRequest, "Recipient not provided")
	}

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

	var phoneBookNumbers []models.PhoneBookNumber
	if err := phoneNumberQuery.Preload("PhoneBook").Find(&phoneBookNumbers).Error; err != nil {
		return c.String(http.StatusBadRequest, "Recipient does not exist in the phone book")
	}

	smsCount := len(phoneBookNumbers)
	reduceErr := reduceAccountBudget(db, account, smsCount)

	if reduceErr != nil {
		return c.String(http.StatusBadRequest, "Insufficient budget")
	}

	for _, phoneBookNumber := range phoneBookNumbers {
		templateMessage := CreateSMSTemplate(request.Message, phoneBookNumber)
		sms := &models.SMSMessage{
			Sender:    account.Username,
			Recipient: phoneBookNumber.Phone,
			Message:   templateMessage,
			AccountID: account.ID,
			Schedule:  &scheduleTime,
		}

		intervalDuration := time.Duration(request.Interval) * time.Second

		go sendSMSWithRepetition(db, sms, intervalDuration, account)
	}

	return c.String(http.StatusOK, "SMS scheduled successfully")
}

func sendSMSWithRepetition(db *gorm.DB, sms *models.SMSMessage, interval time.Duration, account models.Account) {
	for {
		now := time.Now()

		remainingTime := time.Until(*sms.Schedule)
		if remainingTime > 0 {
			timer := time.NewTimer(remainingTime)
			<-timer.C
		}

		reduceErr := reduceAccountBudget(db, account, 1)
		if reduceErr != nil {
			log.Printf("Failed to reduce account budget: %s", reduceErr.Error())
			break
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
		}

		db.Model(sms).Where("id = ?", sms.ID).Updates(models.SMSMessage{
			DeliveryReport: sms.DeliveryReport,
		})

		time.Sleep(time.Second)

		nextSchedule := now.Add(interval)

		for nextSchedule.Before(time.Now()) {
			nextSchedule = nextSchedule.Add(interval)
		}

		sms.Schedule = &nextSchedule
	}
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

func reduceAccountBudget(db *gorm.DB, account models.Account, smsCount int) error {
	var smsCost int

	if smsCount == 1 {
		if err := db.Table("configuration").
			Where("name = ?", "single sms").
			Pluck("value", &smsCost).Error; err != nil {
			log.Printf("Failed to retrieve single SMS cost: %s", err.Error())
			return err
		}
	} else {
		if err := db.Table("configuration").
			Where("name IN (?)", []string{"single sms", "group sms"}).
			Pluck("value", &smsCost).Error; err != nil {
			log.Printf("Failed to retrieve SMS costs: %s", err.Error())
			return err
		}
	}

	totalCost := smsCost * smsCount

	if account.Budget < int64(totalCost) {
		return fmt.Errorf("insufficient budget")
	}

	account.Budget -= int64(totalCost)

	if err := db.Save(&account).Error; err != nil {
		log.Printf("Failed to update account's budget: %s", err.Error())
		return err
	}

	return nil
}
