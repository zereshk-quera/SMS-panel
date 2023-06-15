package server

import (
	"SMS-panel/handlers"

	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func smsRouter(e *echo.Echo) {
	e.POST("sms/single-sms", handlers.SendSingleSMSHandler, middlewares.IsLoggedIn)
}
