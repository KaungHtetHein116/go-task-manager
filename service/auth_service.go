package service

import (
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/models"
	"github.com/KaungHtetHein116/personal-task-manager/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

func (h *UserHandler) Register(c echo.Context) error {
	v := validator.New()
	var input models.UserRegisterInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	if err := v.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationErrors(validationErrors)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "validation error",
				"errors":  formattedErrors,
			})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot create the user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

func (h *UserHandler) Login(c echo.Context) error {
	return nil
}
