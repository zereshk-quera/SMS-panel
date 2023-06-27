package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func adminRoutes(e *echo.Echo) {
	e.POST("/admin/login", WithDBConnection(handlers.AdminLoginHandler))
	e.POST("/admin/register", WithDBConnection(handlers.AdminRegisterHandler))
	e.POST("/admin/add-config", handlers.AddConfigHandler, middlewares.IsAdmin)
	e.GET("/admin/sms-report", handlers.SmsReportHandler, middlewares.IsAdmin)
	e.POST("/admin/add-bad-word/:word", handlers.AddBadWordHandler, middlewares.IsAdmin)
	e.GET("/admin/search/:word", handlers.SmsSearchHandler, middlewares.IsAdmin)
	e.PATCH("/admin/deactivate/:id", handlers.DeactivateHandler, middlewares.IsAdmin)
	e.PATCH("/admin/activate/:id", handlers.ActivateHandler, middlewares.IsAdmin)
}
