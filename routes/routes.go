package routes

import (
	"os"

	"github.com/KaungHtetHein116/personal-task-manager/middleware"
	"github.com/KaungHtetHein116/personal-task-manager/service"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, userHandler *service.UserHandler, projectHandler *service.ProjectHandler) {
	e.POST("/user/register", userHandler.Register)
	e.POST("/user/login", userHandler.Login)

	jwtSecret := os.Getenv("JWT_SECRET")
	protected := e.Group("/protected")
	protected.Use(middleware.JWTMiddleware(jwtSecret))
	protected.GET("/user/me", userHandler.GetProfile)

	protected.POST("/projects", projectHandler.CreateProject)
	protected.GET("/projects/:id", projectHandler.GetProjectByID)
	protected.GET("/projects", projectHandler.GetProjects)
	protected.PATCH("/projects/:id", projectHandler.UpdateProject)
	protected.DELETE("/projects/:id", projectHandler.DeleteProject)
}
