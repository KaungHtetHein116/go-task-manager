package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/models"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/response"
	redisdb "github.com/KaungHtetHein116/personal-task-manager/internal/redis"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	projectRepo repository.ProjectRepository
}

func NewProjectHandler(repo repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: repo}
}

func (r *ProjectHandler) CreateProject(c echo.Context, input *request.CreateProjectInput) error {
	userID := c.Get("user_id").(uint)

	if r.projectRepo.IsProjectExist(input.Name, userID) {
		return transport.NewApiErrorResponse(c, http.StatusConflict, "project name already exists", nil)
	}

	project := &models.Project{Name: input.Name, UserID: userID, Description: &input.Description}
	if err := r.projectRepo.CreateProject(project); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to create in db", nil)
	}

	return transport.NewApiCreateSuccessResponse(c, "Project created", echo.Map{
		"ID":          project.ID,
		"name":        input.Name,
		"description": input.Description,
	})
}

func (r *ProjectHandler) GetProjects(c echo.Context) error {
	var resp []response.ProjectsResponse

	userID := c.Get("user_id").(uint)

	projects, err := r.projectRepo.GetUserProjects(userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "db get error", nil)
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

func (r *ProjectHandler) UpdateProject(c echo.Context, input *request.CreateProjectInput) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)

	userID := c.Get("user_id").(uint)

	// check if project exsit
	exist := r.projectRepo.IsProjectExistByID(uint(projectID), userID)
	if !exist {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	var project models.Project

	// Update project fields
	project.Name = input.Name
	project.Description = &input.Description
	project.UserID = userID
	project.ID = uint(projectID)

	// Save the updated project
	if err := r.projectRepo.UpdateProject(&project); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "failed to update project", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "project updated successfully", echo.Map{
		"id":          project.ID,
		"name":        project.Name,
		"description": project.Description,
	})
}

func (r *ProjectHandler) DeleteProject(c echo.Context) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)
	userID := c.Get("user_id").(uint)

	// check if project exsit
	exist := r.projectRepo.IsProjectExistByID(uint(projectID), userID)
	if !exist {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	if err := r.projectRepo.DeleteProject(uint(projectID), userID); err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "database delete error", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "project deleted", nil)
}

func (r *ProjectHandler) GetProjectByID(c echo.Context) error {
	idParam := c.Param("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "invalid ID", nil)
	}

	userID := c.Get("user_id").(uint)
	cacheKey := fmt.Sprintf("project:%d", projectID)

	// ✅ Step 1: Try to get from cache
	var cachedProject models.Project
	found, _ := redisdb.Get(cacheKey, &cachedProject)

	if found {
		return transport.NewApiSuccessResponse(c, http.StatusOK, "successful (from cache)", cachedProject)
	}

	// 🐢 Step 2: Fetch from DB if not in cache
	project, err := r.projectRepo.GetProjectByID(uint(projectID), userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}

	// 💾 Step 3: Cache the result for 10 mins
	_ = redisdb.Set(cacheKey, project)

	return transport.NewApiSuccessResponse(c, http.StatusOK, "successful (from DB)", project)
}
