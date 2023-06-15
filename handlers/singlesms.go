package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	database "SMS-panel/database"
	"SMS-panel/models"
)

/*
// SendSMSRequest represents the request body for sending an SMS message
type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number" example:"1234567890"`
	Message     string `json:"message" example:"Hello, World!"`
}

// SendSMSResponse represents the response for sending an SMS message
type SendSMSResponse struct {
	Message string `json:"message" example:"SMS sent successfully"`
}

// ErrorResponseSingle represents the structure of an error response
type ErrorResponseSingle struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
*/
// SendSMSRequest represents the request body for sending an SMS message.
type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number" example:"1234567890"`
	Message     string `json:"message" example:"Hello, World!"`
}

// SendSMSResponse represents the response for sending an SMS message.
type SendSMSResponse struct {
	Message string `json:"message" example:"SMS sent successfully"`
}

// ErrorResponseSingle represents the structure of an error response.
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
	// Retrieve the logged-in account from the context
	account := c.Get("account").(models.Account)

	// Read the request body
	reqBody := new(SendSMSRequest)
	if err := c.Bind(reqBody); err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	// Check if the user has sufficient budget
	if account.Budget < 50 {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusForbidden,
			Message: "Insufficient budget",
		}
		return c.JSON(http.StatusForbidden, errResponse)
	}

	// Reduce the budget by 50
	account.Budget -= 50

	// Call the mock API to send the SMS message
	deliveryReport, err := SendMessageHandler(&Message{
		Text:        reqBody.Message,
		Source:      account.Username,
		Destination: reqBody.PhoneNumber,
	})
	if err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to send SMS",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	// Create the SMS message object
	sms := models.SMSMessage{
		Sender:         account.Username,
		Recipient:      reqBody.PhoneNumber,
		Message:        reqBody.Message,
		Schedule:       nil,
		DeliveryReport: deliveryReport,
		CreatedAt:      time.Now(),
		AccountID:      account.ID,
	}

	// Save the SMS message in the database
	db, err := database.GetConnection()
	if err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	if err := db.Create(&sms).Error; err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to save SMS message",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	// Update the account's budget in the database
	if err := db.Model(&account).Update("budget", account.Budget).Error; err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update account's budget",
		}
		return c.JSON(http.StatusInternalServerError, errResponse)
	}

	// Return success response
	response := SendSMSResponse{
		Message: "SMS sent successfully",
	}
	return c.JSON(http.StatusOK, response)
}
