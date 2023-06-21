package server

import (
	"SMS-panel/handlers"

	"github.com/labstack/echo/v4"
)

func phonebookRoutes(e *echo.Echo, handler *handlers.PhonebookHandler) {

	// Register routes
	e.POST("/account/phone-books/", handler.CreatePhoneBook)
	e.GET("/account/:accountID/phone-books/", handler.GetAllPhoneBooks)
	e.GET("/account/:accountID/phone-books/:phoneBookID", handler.ReadPhoneBook)
	e.PUT("/account/:accountID/phone-books/:phoneBookID", handler.UpdatePhoneBook)
	e.DELETE("/account/:accountID/phone-books/:phoneBookID", handler.DeletePhoneBook)

	// Phone book number URLs
	e.POST("/account/phone-books/phone-book-numbers", handler.CreatePhoneBookNumber)
	e.GET("/account/phone-books/:phoneBookID/phone-book-numbers", handler.ListPhoneBookNumbers)
	e.GET("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.ReadPhoneBookNumber)
	e.PUT("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.UpdatePhoneBookNumber)
	e.DELETE("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.DeletePhoneBookNumber)
}
