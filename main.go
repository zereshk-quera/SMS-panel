package main

import (
	"SMS-panel/handlers"

	echo "github.com/labstack/echo/v4"
)

func main() {
	// Create a new Echo instance
	e := echo.New()

	// Register routes
	e.POST("/phone-books/phone-book-numbers", handlers.CreatePhoneBookNumber)
	e.GET("/phone-books/:phoneBookID/phone-book-numbers", handlers.ListPhoneBookNumbers)
	e.GET("/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.ReadPhoneBookNumber)
	e.PUT("/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.UpdatePhoneBookNumber)
	e.DELETE("/phone-books/phone-book-numbers/:phoneBookNumberID", handlers.DeletePhoneBookNumber)

	// Start the server
	e.Start(":8080")
}
