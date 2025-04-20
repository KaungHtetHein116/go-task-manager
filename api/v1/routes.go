package v1

import (
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/handler"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	constants "github.com/KaungHtetHein116/personal-task-manager/pkg/constant"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUserRoute(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	userRoutes := e.Group(constants.USER_API_PREFIX)
	userRoutes.POST("/register", userHandler.Register)
	userRoutes.POST("/login", userHandler.Login)
	userRoutes.GET("/me", userHandler.GetProfile)
}

func RegisterProjectRoute(e *echo.Echo, db *gorm.DB) {
	projectRepo := repository.NewProjectRepository(db)
	projectHandler := handler.NewProjectHandler(projectRepo)

	projectRoutes := e.Group(constants.PROJECT_API_PREFIX)
	projectRoutes.POST("", projectHandler.CreateProject)
	projectRoutes.GET("", projectHandler.GetProjects)
	projectRoutes.GET("/:id", projectHandler.GetProjectByID)
	projectRoutes.PATCH("/:id", projectHandler.UpdateProject)
	projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
}
