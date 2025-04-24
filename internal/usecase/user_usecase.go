package usecase

import (
	"errors"

	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
)

type UserUsecase interface {
	Register(input *request.RegisterUserInput) error
	Login(input *request.LoginUserInput) (string, *entity.User, error)
	GetProfile(userID uint) (*entity.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{repo}
}

func (u *userUsecase) Register(input *request.RegisterUserInput) error {
	existingUser, err := u.repo.GetUserByEmail(input.Email)
	if err != nil && !errors.Is(err, utils.ErrRecordNotFound) {
		return err
	}
	if existingUser != nil {
		return utils.ErrDuplicateEntry
	}

	user := &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: utils.GenerateHashedPassword(input.Password),
	}

	return u.repo.CreateUser(user)
}

func (u *userUsecase) Login(input *request.LoginUserInput) (string, *entity.User, error) {
	user, err := u.repo.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, utils.ErrRecordNotFound) {
			return "", nil, utils.ErrInvalidData
		}
		return "", nil, err
	}

	if !utils.ComparePasswords(user.Password, input.Password) {
		return "", nil, utils.ErrInvalidData
	}

	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *userUsecase) GetProfile(userID uint) (*entity.User, error) {
	user, err := u.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
