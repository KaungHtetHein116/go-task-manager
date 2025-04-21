package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *entity.Project) error
	GetUserProjects(userID uint) ([]entity.Project, error)
	IsProjectExist(name string, userID uint) bool
	IsProjectExistByID(id uint, userID uint) bool
	GetProjectByID(id uint, userID uint) (*entity.Project, error)
	UpdateProject(project *entity.Project) error
	DeleteProject(id uint, userID uint) error
}

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepo{db: db}
}

func (r *projectRepo) CreateProject(project *entity.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepo) GetUserProjects(userID uint) ([]entity.Project, error) {
	var projects []entity.Project
	err := r.db.Where("user_id = ?", userID).Find(&projects).Error
	return projects, err
}

func (r *projectRepo) IsProjectExist(name string, userID uint) bool {
	var project entity.Project
	err := r.db.Model(&entity.Project{}).
		Where("name = ? AND user_id = ?", name, userID).
		First(&project).Error

	return err == nil
}

func (r *projectRepo) IsProjectExistByID(id uint, userID uint) bool {
	var project entity.Project
	err := r.db.Model(&entity.Project{}).
		Where("ID = ? AND user_id = ?", id, userID).
		First(&project).Error

	return err == nil
}

func (r *projectRepo) GetProjectByID(id uint, userID uint) (*entity.Project, error) {
	var project entity.Project
	err := r.db.Model(&entity.Project{}).
		Preload("User").
		Where("ID = ? AND user_id = ?", id, userID).
		First(&project).Error

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepo) UpdateProject(project *entity.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepo) DeleteProject(id uint, userID uint) error {
	return r.db.Where("user_id = ? AND ID = ?", userID, id).Delete(&entity.Project{}).Error
}
