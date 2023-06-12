package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func accountRoutes(e *echo.Echo) {
	e.POST("/accounts/login", handlers.LoginHandler)
	e.POST("/accounts/register", handlers.RegisterHandler)
	e.GET("/accounts/budget", handlers.BudgetAmountHandler, middlewares.IsLoggedIn)
}
