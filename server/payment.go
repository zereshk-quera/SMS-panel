package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func paymentRoutes(e *echo.Echo) {
	e.POST("/accounts/payment/request", WithDBConnection(handlers.PaymentRequestHandler), middlewares.IsLoggedIn)
	e.GET("/accounts/payment/verify", handlers.PaymentVerifyHandler)
}
