package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetUserProjects(userID uint) ([]models.Project, error)
	IsProjectExist(name string, userID uint) bool
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepo) GetUserProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("user_id = ?", userID).Find(&projects).Error
	return projects, err
}

func (r *projectRepo) IsProjectExist(name string, userID uint) bool {
	var count int64
	r.db.Model(&models.Project{}).Where("name = ? AND user_id = ?", name, userID).Count(&count)

	return count > 0
}
