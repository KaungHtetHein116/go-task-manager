package v1

import (
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/handler"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	constants "github.com/KaungHtetHein116/personal-task-manager/pkg/constant"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUserRoute(e *echo.Echo, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	userRoutes := e.Group(constants.USER_API_PREFIX)
	userRoutes.POST("/register", utils.BindAndValidateDecorator(userHandler.Register))
	userRoutes.POST("/login", utils.BindAndValidateDecorator(userHandler.Login))
	userRoutes.GET("/me", userHandler.GetProfile)
}

func RegisterProjectRoute(e *echo.Echo, db *gorm.DB) {
	projectRepo := repository.NewProjectRepository(db)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	projectRoutes := e.Group(constants.PROJECT_API_PREFIX)
	projectRoutes.POST("", utils.BindAndValidateDecorator(projectHandler.CreateProject))
	projectRoutes.GET("", projectHandler.GetProjects)
	projectRoutes.GET("/:id", projectHandler.GetProjectByID)
	projectRoutes.PATCH("/:id", utils.BindAndValidateDecorator(projectHandler.UpdateProject))
	projectRoutes.DELETE("/:id", projectHandler.DeleteProject)
}
