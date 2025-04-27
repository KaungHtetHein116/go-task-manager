package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"gorm.io/gorm"
)

type LabelRepository interface {
	CreateLabel(label *entity.Label) error
	GetLabels(userID uint) ([]entity.Label, error)
	GetLabelByID(id uint) (*entity.Label, error)
	UpdateLabel(label *entity.Label) error
	DeleteLabel(id uint) error
}

type labelRepository struct {
	db *gorm.DB
}

func NewLabelRepository(db *gorm.DB) LabelRepository {
	return &labelRepository{
		db: db,
	}
}

func (r *labelRepository) CreateLabel(label *entity.Label) error {
	return r.db.Create(label).Error
}

func (r *labelRepository) GetLabels(userID uint) ([]entity.Label, error) {
	var labels []entity.Label
	err := r.db.Where("user_id = ?", userID).Find(&labels).Error
	return labels, err
}

func (r *labelRepository) GetLabelByID(id uint) (*entity.Label, error) {
	var label entity.Label
	err := r.db.First(&label, id).Error
	return &label, err
}

func (r *labelRepository) UpdateLabel(label *entity.Label) error {
	return r.db.Model(label).Updates(map[string]interface{}{
		"name": label.Name,
	}).Error
}

func (r *labelRepository) DeleteLabel(id uint) error {
	return r.db.Delete(&entity.Label{}, id).Error
}
