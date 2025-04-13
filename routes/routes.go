package routes

import (
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/middleware"
	"github.com/KaungHtetHein116/personal-task-manager/service"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, h *service.UserHandler) {
	e.POST("/user/register", h.Register)
	e.POST("/user/login", h.Login)

	jwtSecret := os.Getenv("JWT_SECRET")
	protected := e.Group("/protected")
	protected.Use(middleware.JWTMiddleware(jwtSecret))
	protected.GET("/user/me", h.GetProfile)
}
