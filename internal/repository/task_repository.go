package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(task *entity.Task) error
	IsTaskExist(name string) (bool, error)
	IsProjectExist(projectID, userID uint) (bool, error)
	GetTasks(userID uint) ([]entity.Task, error)
	GetTaskByID(taskID, userID uint) (*entity.Task, error)
	UpdateTask(task *entity.Task) error
	DeleteTask(taskID, userID uint) error
}

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(task *entity.Task) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create or get existing labels
		if len(task.Labels) > 0 {
			for i := range task.Labels {
				var existingLabel entity.Label
				err := tx.Where("name = ? AND user_id = ?", task.Labels[i].Name, task.UserID).First(&existingLabel).Error
				if err == gorm.ErrRecordNotFound {
					// Create new label if it doesn't exist
					if err := tx.Create(&task.Labels[i]).Error; err != nil {
						return err
					}
				} else if err != nil {
					return err
				} else {
					// Use existing label
					task.Labels[i] = existingLabel
				}
			}
		}

		// Create the task
		if err := tx.Create(task).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *taskRepo) IsTaskExist(name string) (bool, error) {
	var task entity.Task
	err := r.db.Where("title = ?", name).First(&task).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *taskRepo) IsProjectExist(projectID, userID uint) (bool, error) {
	var project entity.Project
	err := r.db.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *taskRepo) GetTasks(userID uint) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.Preload("Labels").Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}

func (r *taskRepo) GetTaskByID(taskID, userID uint) (*entity.Task, error) {
	var task entity.Task
	err := r.db.Preload("Labels").Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) UpdateTask(task *entity.Task) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, remove existing label associations
		if err := tx.Model(task).Association("Labels").Clear(); err != nil {
			return err
		}

		// Create or get existing labels
		if len(task.Labels) > 0 {
			for i := range task.Labels {
				var existingLabel entity.Label
				err := tx.Where("name = ? AND user_id = ?", task.Labels[i].Name, task.UserID).First(&existingLabel).Error
				if err == gorm.ErrRecordNotFound {
					// Create new label if it doesn't exist
					if err := tx.Create(&task.Labels[i]).Error; err != nil {
						return err
					}
				} else if err != nil {
					return err
				} else {
					// Use existing label
					task.Labels[i] = existingLabel
				}
			}
		}

		// Update the task
		return tx.Save(task).Error
	})
}

func (r *taskRepo) DeleteTask(taskID, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", taskID, userID).Delete(&entity.Task{}).Error
}
