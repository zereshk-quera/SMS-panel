package server

import (
	"log"

	database "SMS-panel/database"
	"SMS-panel/handlers"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer() {
	db, err := database.GetConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Account
	accountRoutes(e)

	// Phonebook
	phonebookHandler := handlers.NewPhonebookHandler(db)
	phonebookRoutes(e, phonebookHandler)

	// SMS
	smsHandler := handlers.NewSmsPhoneBookHandler(db)
	smsRouter(e, smsHandler)

	log.Fatal(e.Start(":8080"))
}
