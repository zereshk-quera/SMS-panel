package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func accountRoutes(e *echo.Echo, handler *handlers.AccountHandler) {
	e.POST("/accounts/login", handler.LoginHandler)
	e.POST("/accounts/register", handler.RegisterHandler)
	e.POST("/accounts/rent_number", handler.RentNumberHandler, middlewares.IsLoggedIn)
	e.GET("/accounts/budget", handler.BudgetAmountHandler, middlewares.IsLoggedIn)
	e.GET("/accounts/sender_numbers", handler.GetAllSenderNumbersHandler, middlewares.IsLoggedIn)
}
