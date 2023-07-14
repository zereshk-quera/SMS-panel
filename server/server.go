package server

import (
	"log"
	"time"

	database "SMS-panel/database"
	"SMS-panel/handlers"
	"SMS-panel/tasks"

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

	// task scheduler
	taskSchaduler := tasks.NewTaskScheduler()
	taskSchaduler.AddTask(tasks.RentNumberTask(db), 10*time.Second, 11, 51, 0)

	taskSchaduler.Run()

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Account
	// accountHandler := handlers.NewAccountHandler(db)
	// accountRoutes(e, accountHandler)
	accountRoutes(e)
	// Payment
	paymentRoutes(e)

	// Phonebook
	phonebookHandler := handlers.NewPhonebookHandler(db)
	phonebookRoutes(e, phonebookHandler)

	// SMS
	smsHandler := handlers.NewSmsPhoneBookHandler(db)
	smsRouter(e, smsHandler)

//admin
	adminRoutes(e)

	log.Fatal(e.Start(":8080"))
}
