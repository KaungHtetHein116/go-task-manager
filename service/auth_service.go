package service

import (
	"errors"
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/models"
	"github.com/KaungHtetHein116/personal-task-manager/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

var v = validator.New()

func (h *UserHandler) Register(c echo.Context) error {
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

	// Check if email already exists
	existingUser, err := h.userRepo.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "email already exists"})
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: utils.GenerateHashedPassword(input.Password),
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot create the user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

func (h *UserHandler) Login(c echo.Context) error {
	var input models.UserLoginInput
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
	}

	// Get user by email
	user, err := h.userRepo.GetUserByEmail(input.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid credentials"})
		}

		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "unkown error"})
	}

	// Compare password
	if !utils.ComparePasswords(user.Password, input.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Invalid password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID, user.Username, user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "token generation error",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"data": echo.Map{
			"name":  user.Username,
			"token": token,
		},
	})
}
