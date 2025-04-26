package handler

import (
	"net/http"
	"strconv"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/response"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	"github.com/KaungHtetHein116/personal-task-manager/pkg/constants"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TaskHandler struct {
	usecase usecase.TaskUsecase
}

func NewTaskHandler(usecase usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{usecase}
}

func (h *TaskHandler) CreateTask(c echo.Context, input *request.CreateTaskInput) error {
	userID := c.Get("user_id").(uint)
	input.UserID = userID

	found, err := h.usecase.IsTaskExist(input.Title)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	if found {
		return transport.NewApiErrorResponse(c, http.StatusConflict, constants.ErrDuplicatedData, nil)
	}

	if err = h.usecase.CreateTask(input); err != nil {
		if err == utils.ErrProjectNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrProjectNotFound, nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiCreateSuccessResponse(c, constants.SuccTaskCreated, nil)
}

func (h *TaskHandler) GetTasks(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	tasks, err := h.usecase.GetTasks(userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tasks", nil)
	}

	var taskResponses []response.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, response.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			ProjectID:   task.ProjectID,
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Tasks retrieved successfully", taskResponses)
}

func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid task ID", nil)
	}

	task, err := h.usecase.GetTaskByID(uint(taskID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Task not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve task", nil)
	}

	taskResponse := response.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		ProjectID:   task.ProjectID,
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Task retrieved successfully", taskResponse)
}

func (h *TaskHandler) UpdateTask(c echo.Context, input *request.UpdateTaskInput) error {
	userID := c.Get("user_id").(uint)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid task ID", nil)
	}
	input.UserID = userID

	err = h.usecase.UpdateTask(uint(taskID), input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Task not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to update task", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Task updated successfully", nil)
}

func (h *TaskHandler) DeleteTask(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, "Invalid task ID", nil)
	}

	err = h.usecase.DeleteTask(uint(taskID), userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, "Task not found", nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "Failed to delete task", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}
