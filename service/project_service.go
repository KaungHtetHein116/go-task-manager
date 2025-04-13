package service

import (
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/models"
	"github.com/KaungHtetHein116/personal-task-manager/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: repo}
}

func (r *ProjectHandler) CreateProject(c echo.Context) error {
	var input *models.CreateProjectInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid input",
		})
	}

	// validate
	if err := utils.Validation.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationErrors(validationErrors)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "validation error",
				"errors":  formattedErrors,
			})
		}
	}

	userID := c.Get("user_id").(uint)

	if r.projectRepo.IsProjectExist(input.Name, userID) {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "project name already exists",
		})
	}

	project := &models.Project{Name: input.Name, UserID: userID}
	if err := r.projectRepo.CreateProject(project); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create in db",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Project created",
		"data": echo.Map{
			"name": input.Name,
		},
	})
}
