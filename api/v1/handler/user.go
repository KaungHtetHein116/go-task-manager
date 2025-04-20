package handler

import (
	"errors"
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/models"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/internal/repository"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
}

// Register handles user registration by validating input, checking for duplicate emails,
// and creating a new user with a hashed password. Returns HTTP 201 on success.
// Returns HTTP 400 for invalid input, HTTP 409 for duplicate email, or HTTP 500 for server errors.
func (h *UserHandler) Register(c echo.Context, input *request.UserRegisterInput) error {
	// Check if email already exists
	existingUser, err := h.userRepo.GetUserByEmail(input.Email)
	if err == nil && existingUser != nil {
		// Return 409 when email is already registered
		return transport.NewApiErrorResponse(c, http.StatusConflict, "email already exists", nil)
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: utils.GenerateHashedPassword(input.Password),
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		// Return 500 for database operation failures
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "cannot create the user", nil)
	}

	// Return 201 for successful user creation
	return transport.NewApiCreateSuccessResponse(c, "user created successfully", nil)
}

// Login authenticates a user with email and password, returning a JWT token on success.
// Returns HTTP 200 with token on success, HTTP 400 for invalid input,
// HTTP 401 for invalid credentials, or HTTP 500 for server errors.
func (h *UserHandler) Login(c echo.Context, input *request.UserLoginInput) error {
	// Get user by email
	user, err := h.userRepo.GetUserByEmail(input.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return 401 when user not found to avoid user enumeration
			return transport.NewApiErrorResponse(c, http.StatusUnauthorized, "invalid credentials", nil)
		}

		// Return 500 for unexpected database errors
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "unknown error", nil)
	}

	// Compare password
	if !utils.ComparePasswords(user.Password, input.Password) {
		// Return 401 for incorrect password
		return transport.NewApiErrorResponse(c, http.StatusUnauthorized, "invalid password", nil)
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "token generation error", nil)
	}

	// Return 200 with JWT token and user data on successful login
	return transport.NewApiSuccessResponse(c, http.StatusOK, "login successful", echo.Map{
		"name":  user.Username,
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.userRepo.GetUserByID(userID)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "user not found", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "successful", user)
}
