package server

import (
	"SMS-panel/handlers"

	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func smsRouter(e *echo.Echo, handler *handlers.SmsPhoneBookHandler) {
	e.POST("/sms/single-sms", handlers.SendSingleSMSHandler, middlewares.IsLoggedIn)
	e.POST("/sms/periodic-sms", handlers.PeriodicSendSMSHandler, middlewares.IsLoggedIn)
	e.POST("/sms/phonebooks", handler.SendMessageToPhoneBooksHandler, middlewares.IsLoggedIn)
}