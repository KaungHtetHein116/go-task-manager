package v1

import (
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/handler"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	constants "github.com/KaungHtetHein116/personal-task-manager/pkg/constants"
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

func RegisterTaskRoute(e *echo.Echo, db *gorm.DB) {
	taskRepo := repository.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUsecase(taskRepo)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	taskRoutes := e.Group(constants.TASK_API_PREFIX)
	taskRoutes.POST("", utils.BindAndValidateDecorator(taskHandler.CreateTask))
	taskRoutes.GET("", taskHandler.GetTasks)
	taskRoutes.GET("/:id", taskHandler.GetTaskByID)
	taskRoutes.PATCH("/:id", utils.BindAndValidateDecorator(taskHandler.UpdateTask))
	taskRoutes.DELETE("/:id", taskHandler.DeleteTask)
}

func RegisterLabelRoute(e *echo.Echo, db *gorm.DB) {
	labelRepo := repository.NewLabelRepository(db)
	labelUsecase := usecase.NewLabelUsecase(labelRepo)
	labelHandler := handler.NewLabelHandler(labelUsecase)

	labelRoutes := e.Group(constants.LABEL_API_PREFIX)
	labelRoutes.POST("", utils.BindAndValidateDecorator(labelHandler.CreateLabel))
	labelRoutes.GET("", labelHandler.GetLabels)
	labelRoutes.GET("/:id", labelHandler.GetLabelByID)
	labelRoutes.PATCH("/:id", utils.BindAndValidateDecorator(labelHandler.UpdateLabel))
	labelRoutes.DELETE("/:id", labelHandler.DeleteLabel)
}
