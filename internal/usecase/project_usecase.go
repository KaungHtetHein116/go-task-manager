package usecase

import (
	"errors"

	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"gorm.io/gorm"
)

type ProjectUsecase interface {
	CreateProject(userID uint, name string, description *string) (*entity.Project, error)
	GetProjects(userID uint) ([]entity.Project, error)
	UpdateProject(project *entity.Project) error
	DeleteProject(projectID, userID uint) error
	GetProjectByID(projectID, userID uint) (*entity.Project, error)
	IsProjectExistByID(projectID, userID uint) bool
}

type projectUsecase struct {
	repo repository.ProjectRepository
}

func NewProjectUsecase(repo repository.ProjectRepository) ProjectUsecase {
	return &projectUsecase{repo}
}

func (u *projectUsecase) CreateProject(userID uint, name string, description *string) (*entity.Project, error) {
	if u.repo.IsProjectExist(name, userID) {
		return nil, utils.ErrProjectNotFound
	}

	project := &entity.Project{Name: name, UserID: userID, Description: description}
	if err := u.repo.CreateProject(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (u *projectUsecase) GetProjects(userID uint) ([]entity.Project, error) {
	return u.repo.GetUserProjects(userID)
}

func (u *projectUsecase) UpdateProject(project *entity.Project) error {
	err := u.repo.UpdateProject(project)

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return utils.ErrDuplicateEntry
	}

	return err
}

func (u *projectUsecase) DeleteProject(projectID, userID uint) error {
	return u.repo.DeleteProject(projectID, userID)
}

func (u *projectUsecase) GetProjectByID(projectID, userID uint) (*entity.Project, error) {
	return u.repo.GetProjectByID(projectID, userID)
}

func (u *projectUsecase) IsProjectExistByID(projectID, userID uint) bool {
	return u.repo.IsProjectExistByID(projectID, userID)
}
