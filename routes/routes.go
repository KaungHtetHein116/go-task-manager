package routes

import (
	"github.com/KaungHtetHein116/personal-task-manager/service"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, h *service.UserHandler) {
	e.POST("/user/register", h.Register)
	e.POST("/user/login", h.Login)
}
