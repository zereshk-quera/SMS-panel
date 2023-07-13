package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"SMS-panel/models"
	"SMS-panel/utils"

	"github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SendSMSRequestPeriodic struct {
	Username     string `json:"username"`
	SenderNumber string `json:"senderNumbers" binding:"required"`
	Phone        string `json:"phone"`
	Message      string `json:"message"`
	Schedule     string `json:"schedule"`
	Interval     string `json:"interval"`
	PhoneBookID  string `json:"phone_book_id"`
}

// PeriodicSendSMSHandler sends periodic SMS messages
// @Summary Send periodic SMS messages
// @Description Send periodic SMS messages with specified schedule and interval
// @Tags messages
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Token"
// @Param sendSMSRequestPeriodic body SendSMSRequestPeriodic true "SMS message details"
// @Success 200 {string} string "SMS scheduled successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 400 {string} string "Invalid schedule time format"
// @Failure 400 {string} string "Recipient not provided"
// @Failure 400 {string} string "Recipient does not exist in the phone book"
// @Failure 400 {string} string "Insufficient budget"
// @Failure 500 {string} string "Internal server error"
// @Router /sms/periodic [post]
func PeriodicSendSMSHandler(c echo.Context, db *gorm.DB) error {
	account := c.Get("account").(models.Account)
	ctx := c.Request().Context()
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

	// Check if sender number is available
	senderNumberExisted := !utils.IsSenderNumberExist(
		ctx, db, request.SenderNumber, account.UserID,
	)
	if senderNumberExisted {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusNotFound,
			Message: "Sender number not found!",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
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

	scheduler := gocron.NewScheduler(time.UTC)

	switch request.Interval {
	case "hourly":
		_, err := scheduler.Every(1).Hour().StartAt(scheduleTime).Do(func() {
			log.Println("schdule time is 2 ", scheduleTime)
			log.Println("now is in cron jobs ", time.Now().UTC())
			sendSMS(db, request.SenderNumber, phoneBookNumbers, account, request, scheduleTime)
		})
		if err != nil {
			log.Printf("Failed to schedule hourly task: %s", err.Error())
			return c.String(http.StatusInternalServerError, "Failed to schedule SMS")
		}

	case "daily":
		_, err := scheduler.Every(1).Day().At(scheduleTime.Format("15:04")).Do(func() {
			sendSMS(db, request.SenderNumber, phoneBookNumbers, account, request, scheduleTime)
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to schedule SMS")
		}
	}

	scheduler.StartAsync()

	return c.String(http.StatusOK, "SMS scheduled successfully")
}

func sendSMS(db *gorm.DB, senderNumber string, phoneBookNumbers []models.PhoneBookNumber, account models.Account, request SendSMSRequestPeriodic, scheduleTime time.Time) {
	for _, phoneBookNumber := range phoneBookNumbers {
		templateMessage := CreateSMSTemplate(request.Message, phoneBookNumber)
		sms := &models.SMSMessage{
			Sender:    senderNumber,
			Recipient: phoneBookNumber.Phone,
			Message:   templateMessage,
			AccountID: account.ID,
			Schedule:  &scheduleTime,
		}
		reduceErr := reduceAccountBudget(db, account, 1)
		log.Println("Budget reduced")

		if reduceErr != nil {
			log.Printf("Failed to reduce account budget: %s", reduceErr.Error())
			break
		}

		err := db.Create(&sms).Error
		if err != nil {
			log.Printf("Failed to create SMSMessage: %s", err.Error())
			continue
		}

		deliveryReport, err := MockSendMessage(&Message{
			Text:        sms.Message,
			Source:      sms.Sender,
			Destination: sms.Recipient,
		})
		if err != nil {
			sms.DeliveryReport = deliveryReport
			log.Printf("Failed to send SMS: %s", err.Error())
		} else {
			sms.DeliveryReport = deliveryReport
		}

		err = db.Model(&sms).Updates(models.SMSMessage{
			DeliveryReport: sms.DeliveryReport,
		}).Error
		if err != nil {
			log.Printf("Failed to update SMSMessage: %s", err.Error())
		}

		time.Sleep(time.Second)
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
