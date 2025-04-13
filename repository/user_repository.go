package repository

import (
	"errors"

	"github.com/KaungHtetHein116/personal-task-manager/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		// Check for unique constraint violation
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("user with this email already exists")
		}
		// Check for invalid data
		if errors.Is(err, gorm.ErrInvalidData) {
			return errors.New("invalid user data provided")
		}
		// Return a generic database error for other cases
		return errors.New("failed to create user: " + err.Error())
	}
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
