package server

import (
	"log"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

func init() {
	e = echo.New()
}

func StartServer() {
	accountRoutes(e)
	log.Fatal(e.Start(":8080"))
}
