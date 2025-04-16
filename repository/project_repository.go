package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetUserProjects(userID uint) ([]models.Project, error)
	IsProjectExist(name string, userID uint) bool
	IsProjectExistByID(id uint, userID uint) bool
	GetProjectByID(id uint, userID uint) (*models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id uint, userID uint) error
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
	var project models.Project
	err := r.db.Model(&models.Project{}).
		Where("name = ? AND user_id = ?", name, userID).
		First(&project).Error

	return err == nil
}

func (r *projectRepo) IsProjectExistByID(id uint, userID uint) bool {
	var project models.Project
	err := r.db.Model(&models.Project{}).
		Where("ID = ? AND user_id = ?", id, userID).
		First(&project).Error

	return err == nil
}

func (r *projectRepo) GetProjectByID(id uint, userID uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Model(&models.Project{}).
		Preload("User").
		Where("ID = ? AND user_id = ?", id, userID).
		First(&project).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepo) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepo) DeleteProject(id uint, userID uint) error {
	return r.db.Where("user_id = ? AND ID = ?", userID, id).Delete(&models.Project{}).Error
}
