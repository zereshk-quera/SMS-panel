package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func adminRoutes(e *echo.Echo) {
	e.POST("/admin/login", handlers.AdminLoginHandler)
	e.POST("/admin/register", handlers.AdminRegisterHandler)
	e.POST("/admin/deactivate", handlers.DeactivateHandler, middlewares.IsAdmin)
	e.POST("/admin/activate", handlers.ActivateHandler, middlewares.IsAdmin)
	e.POST("/admin/add-config", handlers.AddConfigHandler, middlewares.IsAdmin)
	e.GET("/admin/sms-report", handlers.SmsReportHandler, middlewares.IsAdmin)
	e.POST("/admin/search", handlers.SmsSearchHandler, middlewares.IsAdmin)
	e.POST("/admin/add-bad-word", handlers.AddBadWordHandler, middlewares.IsAdmin)
}
