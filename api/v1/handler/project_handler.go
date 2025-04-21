package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/response"
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
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
		if err.Error() == "project name already exists" {
			return transport.NewApiErrorResponse(c, http.StatusConflict, err.Error(), nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to create project", nil)
	}

	return transport.NewApiCreateSuccessResponse(c, "Project created", echo.Map{
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
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to get projects", nil)
	}

	for _, project := range projects {
		resp = append(resp, response.ProjectsResponse{
			ID:          project.ID,
			Name:        project.Name,
			Tasks:       project.Tasks,
			Description: project.Description,
		})
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "successful", resp)
}

func (h *ProjectHandler) UpdateProject(c echo.Context, input *request.CreateProjectInput) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)

	userID := c.Get("user_id").(uint)

	if !h.usecase.IsProjectExistByID(uint(projectID), userID) {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	project := &entity.Project{
		Name:        input.Name,
		Description: &input.Description,
		UserID:      userID,
	}

	if err := h.usecase.UpdateProject(project); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to update project", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "project updated successfully", echo.Map{
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
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	if err := h.usecase.DeleteProject(uint(projectID), userID); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to delete project", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "project deleted", nil)
}

func (h *ProjectHandler) GetProjectByID(c echo.Context) error {
	idParam := c.Param("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "invalid ID", nil)
	}

	userID := c.Get("user_id").(uint)
	cacheKey := fmt.Sprintf("project:%d", projectID)

	var cachedProject entity.Project
	found, _ := redisdb.Get(cacheKey, &cachedProject)

	if found {
		return transport.NewApiSuccessResponse(c, http.StatusOK, "successful (from cache)", cachedProject)
	}

	project, err := h.usecase.GetProjectByID(uint(projectID), userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	_ = redisdb.Set(cacheKey, project)

	return transport.NewApiSuccessResponse(c, http.StatusOK, "successful (from DB)", project)
}
