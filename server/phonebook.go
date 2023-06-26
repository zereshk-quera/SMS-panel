package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func phonebookRoutes(e *echo.Echo, handler *handlers.PhonebookHandler) {

	// Register routes
	e.POST("/account/phone-books/", handler.CreatePhoneBook, middlewares.IsLoggedIn)
	e.GET("/account/:accountID/phone-books/", handler.GetAllPhoneBooks, middlewares.IsLoggedIn)
	e.GET("/account/:accountID/phone-books/:phoneBookID", handler.ReadPhoneBook, middlewares.IsLoggedIn)
	e.PUT("/account/:accountID/phone-books/:phoneBookID", handler.UpdatePhoneBook, middlewares.IsLoggedIn)
	e.DELETE("/account/:accountID/phone-books/:phoneBookID", handler.DeletePhoneBook, middlewares.IsLoggedIn)

	// Phone book number URLs
	e.POST("/account/phone-books/phone-book-numbers", handler.CreatePhoneBookNumber, middlewares.IsLoggedIn)
	e.GET("/account/phone-books/:phoneBookID/phone-book-numbers", handler.ListPhoneBookNumbers, middlewares.IsLoggedIn)
	e.GET("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.ReadPhoneBookNumber, middlewares.IsLoggedIn)
	e.PUT("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.UpdatePhoneBookNumber, middlewares.IsLoggedIn)
	e.DELETE("/account/phone-books/phone-book-numbers/:phoneBookNumberID", handler.DeletePhoneBookNumber, middlewares.IsLoggedIn)
}
