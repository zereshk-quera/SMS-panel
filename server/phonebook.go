package server

import (
	"SMS-panel/handlers"

	_ "SMS-panel/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @Title SMS-PANEL
// @version 1.0
// @description Quera SMS-PANEL server
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http

func phonebookRoutes(e *echo.Echo) {
	// Register Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register routes
	e.POST("/account/phone-books/", handlers.CreatePhoneBook)
	e.GET("/account/:accountID/phone-books/", handlers.GetAllPhoneBooks)
	e.GET("/account/:accountID/phone-books/:phoneBookID", handlers.ReadPhoneBook)
	e.PUT("/account/:accountID/phone-books/:phoneBookID", handlers.UpdatePhoneBook)
	e.DELETE("/account/:accountID/phone-books/:phoneBookID", handlers.DeletePhoneBook)

	// Phone book number URLs
	e.POST("/account/phone-books/phone-book-numbers", handlers.CreatePhoneBookNumber)
	e.GET("/account/phone-books/:phoneBookID/phone-book-numbers", handlers.ListPhoneBookNumbers)
	e.GET("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.ReadPhoneBookNumber)
	e.PUT("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.UpdatePhoneBookNumber)
	e.DELETE("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.DeletePhoneBookNumber)
}
