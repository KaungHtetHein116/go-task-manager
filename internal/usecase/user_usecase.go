package usecase

import (
	"errors"

	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/jinzhu/gorm"
)

type UserUsecase interface {
	Register(name, email, password string) error
	Login(email, password string) (string, *entity.User, error)
	GetProfile(userID uint) (*entity.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo}
}

func (u *userUsecase) Register(name, email, password string) error {
	// Check if email already exists
	existingUser, err := u.repo.GetUserByEmail(email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	user := &entity.User{
		Username: name,
		Email:    email,
		Password: utils.GenerateHashedPassword(password),
	}

	return u.repo.CreateUser(user)
}

func (u *userUsecase) Login(email, password string) (string, *entity.User, error) {
	// Get user by email
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid credentials")
		}
		return "", nil, err
	}

	// Compare password
	if !utils.ComparePasswords(user.Password, password) {
		return "", nil, errors.New("invalid password")
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		return "", nil, errors.New("token generation error")
	}

	return token, user, nil
}

func (u *userUsecase) GetProfile(userID uint) (*entity.User, error) {
	return u.repo.GetUserByID(userID)
}
