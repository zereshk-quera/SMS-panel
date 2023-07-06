package server

import (
	"net/http"

	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/middlewares"
	"SMS-panel/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func WithDBConnection(handlerFunc func(c echo.Context, db *gorm.DB) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		dbConn, err := database.GetConnection()
		if err != nil {
			return c.JSON(http.StatusBadGateway, models.Response{ResponseCode: 502, Message: "Can't Connect To Database"})
		}
		return handlerFunc(c, dbConn)
	}
}

func accountRoutes(e *echo.Echo) {
	e.POST("/accounts/login", WithDBConnection(handlers.LoginHandler))
	e.POST("/accounts/register", WithDBConnection(handlers.RegisterHandler))
	e.GET("/accounts/budget", handlers.BudgetAmountHandler, middlewares.IsLoggedIn)
	e.POST("/accounts/rent_number", handler.RentNumberHandler, middlewares.IsLoggedIn)
	e.GET("/accounts/sender_numbers", handler.GetAllSenderNumbersHandler, middlewares.IsLoggedIn)
}
