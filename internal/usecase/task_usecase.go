package usecase

import (
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
)

type TaskUsecase interface {
	CreateTask(input *request.CreateTaskInput) error
	IsTaskExist(title string) (bool, error)
	GetTasks(userID uint) ([]entity.Task, error)
	GetTaskByID(taskID, userID uint) (*entity.Task, error)
	UpdateTask(taskID uint, input *request.UpdateTaskInput) error
	DeleteTask(taskID, userID uint) error
}

type taskUsecase struct {
	repo repository.TaskRepository
}

func NewTaskUsecase(repo repository.TaskRepository) TaskUsecase {
	return &taskUsecase{repo}
}

func (u *taskUsecase) CreateTask(input *request.CreateTaskInput) error {
	// Validate if project exists and belongs to user
	exists, err := u.repo.IsProjectExist(input.ProjectID, input.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return utils.ErrProjectNotFound
	}

	task := &entity.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		ProjectID:   input.ProjectID,
		UserID:      input.UserID,
		Priority:    input.Priority,
	}

	// Create or get existing labels and associate them with the task
	if len(input.Labels) > 0 {
		var labels []entity.Label
		for _, labelName := range input.Labels {
			label := entity.Label{
				Name:   labelName,
				UserID: input.UserID,
			}
			labels = append(labels, label)
		}
		task.Labels = labels
	}

	return u.repo.CreateTask(task)
}

func (u *taskUsecase) IsTaskExist(title string) (bool, error) {
	return u.repo.IsTaskExist(title)
}

func (u *taskUsecase) GetTasks(userID uint) ([]entity.Task, error) {
	return u.repo.GetTasks(userID)
}

func (u *taskUsecase) GetTaskByID(taskID, userID uint) (*entity.Task, error) {
	return u.repo.GetTaskByID(taskID, userID)
}

func (u *taskUsecase) UpdateTask(taskID uint, input *request.UpdateTaskInput) error {
	task, err := u.repo.GetTaskByID(taskID, input.UserID)
	if err != nil {
		return err
	}

	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Priority = input.Priority

	// Handle label updates
	if len(input.Labels) > 0 {
		var labels []entity.Label
		for _, labelName := range input.Labels {
			label := entity.Label{
				Name:   labelName,
				UserID: input.UserID,
			}
			labels = append(labels, label)
		}
		task.Labels = labels
	} else {
		task.Labels = []entity.Label{} // Clear labels if empty array is provided
	}

	return u.repo.UpdateTask(task)
}

func (u *taskUsecase) DeleteTask(taskID, userID uint) error {
	return u.repo.DeleteTask(taskID, userID)
}
