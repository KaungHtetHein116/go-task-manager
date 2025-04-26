package repository

import (
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id uint, includeProjects bool) (*entity.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *entity.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return utils.HandleGormError(err, "user")
	}
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, utils.HandleGormError(err, "user")
	}
	return &user, nil
}

func (r *userRepo) GetUserByID(id uint, includeProjects bool) (*entity.User, error) {
	var user entity.User
	query := r.db
	if includeProjects {
		query = query.Preload("Projects")
	}
	if err := query.First(&user, id).Error; err != nil {
		return nil, utils.HandleGormError(err, "user")
	}
	return &user, nil
}
