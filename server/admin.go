package server

import (
	"SMS-panel/handlers"
	"SMS-panel/middlewares"

	"github.com/labstack/echo/v4"
)

func adminRoutes(e *echo.Echo) {
	e.POST("/admin/login", WithDBConnection(handlers.AdminLoginHandler))
	e.POST("/admin/register", WithDBConnection(handlers.AdminRegisterHandler))
	e.POST("/admin/add-config", WithDBConnection(handlers.AddConfigHandler), middlewares.IsAdmin)
	e.GET("/admin/sms-report", WithDBConnection(handlers.SmsReportHandler), middlewares.IsAdmin)
	e.POST("/admin/add-bad-word/:word", WithDBConnection(handlers.AddBadWordHandler), middlewares.IsAdmin)
	e.GET("/admin/search/:word", WithDBConnection(handlers.SmsSearchHandler), middlewares.IsAdmin)
	e.PATCH("/admin/deactivate/:id", WithDBConnection(handlers.DeactivateHandler), middlewares.IsAdmin)
	e.PATCH("/admin/activate/:id", WithDBConnection(handlers.ActivateHandler), middlewares.IsAdmin)
}
