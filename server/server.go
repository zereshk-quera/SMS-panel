package server

import (
	"log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer() {
	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	accountRoutes(e)
	phonebookRoutes(e)
	log.Fatal(e.Start(":8080"))
}
