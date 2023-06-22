package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	database "SMS-panel/database"
	"SMS-panel/models"
	"SMS-panel/utils"
)

type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number" example:"1234567890"`
	Message     string `json:"message" example:"Hello, World!"`
	Username    string `json:"username" example:"johndoe"`
}

type SendSMSResponse struct {
	Message string `json:"message" example:"SMS sent successfully"`
}

type ErrorResponseSingle struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SendSingleSMSHandler sends a single SMS message and saves the result in the SMSMessage table.
// @Summary Send Single SMS
// @Description Sends a single SMS message and saves the result in the SMSMessage table
// @Tags SMS
// @Accept json
// @Produce json
// @Param Cookie header string true "account_token" default("account_token")
// @Param sendSMSRequest body SendSMSRequest true "Request body for sending an SMS message"
// @Success 200 {object} SendSMSResponse
// @Failure 400 {object} ErrorResponseSingle
// @Failure 403 {object} ErrorResponseSingle
// @Failure 500 {object} ErrorResponseSingle
// @Router /sms/single-sms [post]
func SendSingleSMSHandler(c echo.Context) error {
	account := c.Get("account").(models.Account)

	reqBody := new(SendSMSRequest)
	if err := c.Bind(reqBody); err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	if !utils.ValidatePhone(reqBody.PhoneNumber) {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusBadRequest,
			Message: "Invalid phone number",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	db, err := database.GetConnection()
	if err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	tx := db.Begin()

	var singleSMSCost int
	if err := tx.Table("configuration").Where("name = ?", "single sms").Select("value").Scan(&singleSMSCost).Error; err != nil {
		tx.Rollback()
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve single SMS cost",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	if account.Budget < int64(singleSMSCost) {
		tx.Rollback()
		errResponse := ErrorResponseSingle{
			Code:    http.StatusForbidden,
			Message: "Insufficient budget",
		}
		return c.JSON(http.StatusForbidden, errResponse)
	}

	account.Budget -= int64(singleSMSCost)

	var phoneNumber models.PhoneBookNumber
	var message string

	if reqBody.Username != "" {
		if err := tx.
			Joins("JOIN phone_books ON phone_books.id = phone_book_numbers.phone_book_id").
			Where("phone_books.account_id = ? AND phone_book_numbers.username = ?", account.ID, reqBody.Username).
			First(&phoneNumber).Error; err != nil {
			phoneNumber = models.PhoneBookNumber{}
		}
		message = CreateSMSTemplate(reqBody.Message, phoneNumber)

		if phoneNumber.ID == 0 {
			tx.Rollback()
			errResponse := ErrorResponseSingle{
				Code:    http.StatusNotFound,
				Message: "Username not found",
			}
			return c.JSON(http.StatusNotFound, errResponse)
		}
	} else {
		if err := tx.
			Joins("JOIN phone_books ON phone_books.id = phone_book_numbers.phone_book_id").
			Where("phone_books.account_id = ? AND phone_book_numbers.phone = ?", account.ID, reqBody.PhoneNumber).
			First(&phoneNumber).Error; err != nil {
			phoneNumber = models.PhoneBookNumber{}
		}
		message = CreateSMSTemplate(reqBody.Message, phoneNumber)

		if phoneNumber.ID == 0 {
			tx.Rollback()
			errResponse := ErrorResponseSingle{
				Code:    http.StatusNotFound,
				Message: "Phone number not found",
			}
			return c.JSON(http.StatusNotFound, errResponse)
		}
	}

	deliveryReport, err := MockSendMessage(&Message{
		Text:        message,
		Source:      account.Username,
		Destination: reqBody.PhoneNumber,
	})
	if err != nil {
		tx.Rollback()
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to send SMS",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	sms := models.SMSMessage{
		Sender:         account.Username,
		Recipient:      reqBody.PhoneNumber,
		Message:        message,
		Schedule:       nil,
		DeliveryReport: deliveryReport,
		CreatedAt:      time.Now(),
		AccountID:      account.ID,
	}

	if err := tx.Create(&sms).Error; err != nil {
		tx.Rollback()
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save SMS message",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	if err := tx.Model(&account).Update("budget", account.Budget).Error; err != nil {
		tx.Rollback()
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update account's budget",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	tx.Commit()

	response := SendSMSResponse{
		Message: "SMS sent successfully",
	}
	return c.JSON(http.StatusOK, response)
}