package usecase

import (
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
)

type LabelUsecase interface {
	CreateLabel(label *entity.Label) error
	GetLabels(userID uint) ([]entity.Label, error)
	GetLabelByID(id uint) (*entity.Label, error)
	UpdateLabel(label *entity.Label) error
	DeleteLabel(id uint) error
}

type labelUsecase struct {
	repo repository.LabelRepository
}

func NewLabelUsecase(r repository.LabelRepository) LabelUsecase {
	return &labelUsecase{
		repo: r,
	}
}

func (u *labelUsecase) CreateLabel(label *entity.Label) error {
	return u.repo.CreateLabel(label)
}

func (u *labelUsecase) GetLabels(userID uint) ([]entity.Label, error) {
	return u.repo.GetLabels(userID)
}

func (u *labelUsecase) GetLabelByID(id uint) (*entity.Label, error) {
	return u.repo.GetLabelByID(id)
}

func (u *labelUsecase) UpdateLabel(label *entity.Label) error {
	return u.repo.UpdateLabel(label)
}

func (u *labelUsecase) DeleteLabel(id uint) error {
	return u.repo.DeleteLabel(id)
}
