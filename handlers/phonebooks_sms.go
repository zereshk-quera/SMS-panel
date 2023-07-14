package handlers

import (
	"log"
	"net/http"

	"SMS-panel/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SendSMessageToPhoneBooksBody struct {
	Account      models.Account `json:"-"`
	SenderNumber string         `json:"senderNumbers" binding:"required"`
	PhoneBooks   []string       `json:"phoneBooks" binding:"required"`
	Message      string         `json:"message" binding:"required"`
}

type SmsPhoneBookHandler struct {
	db *gorm.DB
}

func NewSmsPhoneBookHandler(db *gorm.DB) *SmsPhoneBookHandler {
	return &SmsPhoneBookHandler{db: db}
}

// @Summary Send sms to phone books numbers
// @Description Send sms to phone books numbers
// @Tags messages
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Token"
// @Param body body SendSMessageToPhoneBooksBody true "Phone books sms details."
// @Success 200 {object} SendSMSResponse
// @Failure 204 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /sms/phonebooks [post]
func (sp *SmsPhoneBookHandler) SendMessageToPhoneBooksHandler(c echo.Context) error {
	account := c.Get("account").(models.Account)
	body := SendSMessageToPhoneBooksBody{}
	body.Account = account
	ctx := c.Request().Context()

	if err := c.Bind(&body); err != nil {
		errResponse := ErrorResponse{
			Message: "Invalid request payload",
		}
		return c.JSON(http.StatusBadRequest, errResponse)
	}

	err := SendMessageToPhoneBooks(ctx, body, sp.db)
	if err != nil {
		errorResponse := ErrorResponse{
			Message: err.Error(),
		}
		switch e := err.(type) {
		case PhoneBooksNotFoundError:
			return c.JSON(http.StatusNotFound, errorResponse)
		case AcountDoesNotHaveBudgetError:
			return c.JSON(http.StatusBadRequest, errorResponse)
		case PhoneBooksNumbersAreEmptyError:
			return c.JSON(http.StatusNoContent, errorResponse)
		case SenderNumberNotFoundError:
			return c.JSON(http.StatusNotFound, errorResponse)
		default:
			log.Println(e)
			errorResponse := ErrorResponse{
				Message: "Internal server error!",
			}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}
	}

	response := SendSMSResponse{
		Message: "Done",
	}
	return c.JSON(http.StatusOK, response)
}
