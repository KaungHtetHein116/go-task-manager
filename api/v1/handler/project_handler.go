package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/response"
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	"github.com/KaungHtetHein116/personal-task-manager/pkg/constants"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	usecase usecase.ProjectUsecase
}

func NewProjectHandler(u usecase.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{usecase: u}
}

func (h *ProjectHandler) CreateProject(c echo.Context, input *request.CreateProjectInput) error {
	userID := c.Get("user_id").(uint)

	project, err := h.usecase.CreateProject(userID, input.Name, &input.Description)
	if err != nil {
		if errors.Is(err, utils.ErrProjectNotFound) {
			return transport.NewApiErrorResponse(c, http.StatusConflict, constants.ErrProjectAlreadyExist, nil)
		}
		return err
	}

	return transport.NewApiCreateSuccessResponse(c, constants.SuccProjectCreated, echo.Map{
		"ID":          project.ID,
		"name":        project.Name,
		"description": project.Description,
	})
}

func (h *ProjectHandler) GetProjects(c echo.Context) error {
	var resp []response.ProjectsResponse

	userID := c.Get("user_id").(uint)

	projects, err := h.usecase.GetProjects(userID)
	if err != nil {
		return err
	}

	for _, project := range projects {
		resp = append(resp, response.ProjectsResponse{
			ID:          project.ID,
			Name:        project.Name,
			Tasks:       project.Tasks,
			Description: project.Description,
		})
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.Successful, resp)
}

func (h *ProjectHandler) UpdateProject(c echo.Context, input *request.CreateProjectInput) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)

	userID := c.Get("user_id").(uint)

	if !h.usecase.IsProjectExistByID(uint(projectID), userID) {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrProjectNotFound, nil)
	}

	project := &entity.Project{
		Name:        input.Name,
		Description: &input.Description,
		UserID:      userID,
	}

	if err := h.usecase.UpdateProject(project); err != nil {
		if errors.Is(err, utils.ErrDuplicateEntry) {
			return err
		}
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.SuccProjectUpdated, echo.Map{
		"id":          project.ID,
		"name":        project.Name,
		"description": project.Description,
	})
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)
	userID := c.Get("user_id").(uint)

	if !h.usecase.IsProjectExistByID(uint(projectID), userID) {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrProjectNotFound, nil)
	}

	if err := h.usecase.DeleteProject(uint(projectID), userID); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, constants.ErrFailedDeleteProject, nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.ErrProjectDeleted, nil)
}

func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	idParam := c.Param("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, constants.ErrInvalidProjectID, nil)
	}

	userID := c.Get("user_id").(uint)
	cacheKey := fmt.Sprintf("project:%d", projectID)

	var cachedProject entity.Project
	found, _ := redisdb.Get(cacheKey, &cachedProject)

	if found {
		return transport.NewApiSuccessResponse(c, http.StatusOK, constants.Successful+" (from cache)", cachedProject)
	}

	project, err := h.usecase.GetProjectByID(uint(projectID), userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrProjectNotFound, nil)
	}

	_ = redisdb.Set(cacheKey, project)

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.Successful+" (from DB)", project)
}
