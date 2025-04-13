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

// Register handles user registration by validating input, checking for duplicate emails,
// and creating a new user with a hashed password. Returns HTTP 201 on success.
// Returns HTTP 400 for invalid input, HTTP 409 for duplicate email, or HTTP 500 for server errors.
func (h *UserHandler) Register(c echo.Context) error {
	var input models.UserRegisterInput
	if err := c.Bind(&input); err != nil {
		// Return 400 for malformed request body
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	if err := utils.Validation.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			formattedErrors := utils.FormatValidationErrors(validationErrors)
			// Return 400 with detailed validation errors
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "validation error",
				"errors":  formattedErrors,
			})
		}
		// Return 400 for other validation failures
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	// Check if email already exists
	existingUser, err := h.userRepo.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		// Return 409 when email is already registered
		return c.JSON(http.StatusConflict, echo.Map{"message": "email already exists"})
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: utils.GenerateHashedPassword(input.Password),
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		// Return 500 for database operation failures
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "cannot create the user"})
	}

	// Return 201 for successful user creation
	return c.JSON(http.StatusCreated, echo.Map{"message": "user created successfully"})
}

// Login authenticates a user with email and password, returning a JWT token on success.
// Returns HTTP 200 with token on success, HTTP 400 for invalid input,
// HTTP 401 for invalid credentials, or HTTP 500 for server errors.
func (h *UserHandler) Login(c echo.Context) error {
	var input models.UserLoginInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	if err := utils.Validation.Struct(input); err != nil {
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
			// Return 401 when user not found to avoid user enumeration
			return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid credentials"})
		}

		// Return 500 for unexpected database errors
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "unkown error"})
	}

	// Compare password
	if !utils.ComparePasswords(user.Password, input.Password) {
		// Return 401 for incorrect password
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "Invalid password",
		})
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "token generation error",
		})
	}

	// Return 200 with JWT token and user data on successful login
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"data": echo.Map{
			"name":  user.Username,
			"token": token,
		},
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.userRepo.GetUserByID(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "successful",
		"data":    user,
	})
}
