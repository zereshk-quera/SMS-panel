package server

import (
	"SMS-panel/handlers"

	"github.com/labstack/echo/v4"
)

func accountRoutes(e *echo.Echo) {

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
