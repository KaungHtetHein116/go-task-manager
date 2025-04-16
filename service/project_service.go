package service

import (
	"fmt"
	"net/http"
	"strconv"

	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
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

	project := &models.Project{Name: input.Name, UserID: userID, Description: &input.Description}
	if err := r.projectRepo.CreateProject(project); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to create in db",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Project created",
		"data": echo.Map{
			"ID":          project.ID,
			"name":        input.Name,
			"description": input.Description,
		},
	})
}

func (r *ProjectHandler) GetProjects(c echo.Context) error {
	var response []models.ProjectsResponse

	userID := c.Get("user_id").(uint)

	projects, err := r.projectRepo.GetUserProjects(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "db get error",
		})
	}

	for _, project := range projects {
		response = append(response, models.ProjectsResponse{
			ID:          project.ID,
			Name:        project.Name,
			Tasks:       project.Tasks,
			Description: project.Description,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successful",
		"data":    response,
	})
}

func (r *ProjectHandler) UpdateProject(c echo.Context) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)

	var input models.CreateProjectInput
	userID := c.Get("user_id").(uint)

	// check json valid
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid input",
		})
	}

	// check if project exsit
	exist := r.projectRepo.IsProjectExistByID(uint(projectID), userID)
	if !exist {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "project not found",
		})
	}

	var project models.Project
	// check field validation
	if err := utils.Validation.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationErrors(validationErrors)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "validation error",
				"errors":  formattedErrors,
			})
		}
	}

	// Update project fields
	project.Name = input.Name
	project.Description = &input.Description
	project.UserID = userID
	project.ID = uint(projectID)

	// Save the updated project
	if err := r.projectRepo.UpdateProject(&project); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "failed to update project",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "project updated successfully",
		"data": echo.Map{
			"id":          project.ID,
			"name":        project.Name,
			"description": project.Description,
		},
	})
}

func (r *ProjectHandler) DeleteProject(c echo.Context) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)
	userID := c.Get("user_id").(uint)

	// check if project exsit
	exist := r.projectRepo.IsProjectExistByID(uint(projectID), userID)
	if !exist {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "project not found",
		})
	}

	if err := r.projectRepo.DeleteProject(uint(projectID), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "database delete error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "project deleted",
	})
}

func (r *ProjectHandler) GetProjectByID(c echo.Context) error {
	idParam := c.Param("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid ID"})
	}

	userID := c.Get("user_id").(uint)
	cacheKey := fmt.Sprintf("project:%d", projectID)

	// ✅ Step 1: Try to get from cache
	var cachedProject models.Project
	found, _ := redisdb.Get(cacheKey, &cachedProject)

	if found {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "successful (from cache)",
			"data":    cachedProject,
		})
	}

	// 🐢 Step 2: Fetch from DB if not in cache
	project, err := r.projectRepo.GetProjectByID(uint(projectID), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "project not found"})
	}

	// 💾 Step 3: Cache the result for 10 mins
	_ = redisdb.Set(cacheKey, project)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successful (from DB)",
		"data":    project,
	})
}
