package handlers

import (
	"SMS-panel/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SendSMessageToPhoneBooksBody struct {
	Account    models.Account `json:"-"`
	PhoneBooks []string          `json:"phoneBooks" binding:"required"`
	Message    string         `json:"message" binding:"required"`
}

type SmsPhoneBookHandler struct {
	db *gorm.DB
}

func NewSmsPhoneBookHandler(db *gorm.DB) *SmsPhoneBookHandler {
	return &SmsPhoneBookHandler{db: db}
}

// @Summary Send sms to phone books numbers
// @Description Send sms to phone books numbers
// @Tags SMS
// @Accept json
// @Produce json
// @Param body body SendSMessageToPhoneBooksBody true "Phone books sms details."
// @Success 200 {object} SendSMSResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sms/phonebooks [post]
func (sp *SmsPhoneBookHandler) SendMessageToPhoneBooksHandler(c echo.Context) error {
	account := c.Get("account").(models.Account)
	body := SendSMessageToPhoneBooksBody{}
	body.Account = account
	ctx := c.Request().Context()

	if err := c.Bind(&body); err != nil {
		errResponse := ErrorResponseSingle{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	err := SendMessageToPhoneBooks(ctx, body, sp.db)
	if err != nil {
		switch e := err.(type) {
		case PhoneBooksNotFoundError:
			errorResponse := ErrorResponseSingle{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}
			return c.JSON(http.StatusNotFound, errorResponse)
		case AcountDoesNotHaveBudgetError:
			errorResponse := ErrorResponseSingle{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
			return c.JSON(http.StatusBadRequest, errorResponse)
		default:
			log.Println(e)
		}
	}

	response := SendSMSResponse{
		Message: "Done",
	}
	return c.JSON(http.StatusOK, response)
}
