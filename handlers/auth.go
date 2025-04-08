package handlers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/KaungHtetHein116/personal-task-manager/database"
	"github.com/KaungHtetHein116/personal-task-manager/models"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type RegisterInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func validateRegisterRequest(req *RegisterInput) []ValidationError {
	var errors []ValidationError

	if err := validate.Struct(req); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var ve ValidationError
			ve.Field = err.Field()
			switch err.Tag() {
			case "required":
				ve.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				ve.Message = fmt.Sprintf("%s must be a valid email", err.Field())
			case "min":
				ve.Message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "max":
				ve.Message = fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
			}
			errors = append(errors, ve)
		}
	}

	// Optional: stricter regex for email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString((*req).Email) {
		errors = append(errors, ValidationError{
			Field:   "Email",
			Message: "Invalid email format",
		})
	}

	return errors
}

func emailExists(email string) bool {
	var user models.User
	return database.DB.Where("email = ?", email).First(&user).Error == nil
}

func RegisterUser(c echo.Context) error {
	var input RegisterInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}

	validationErrors := validateRegisterRequest(&input)
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validationErrors})
	}

	if emailExists(input.Email) {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "email already exists"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := &models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "couldn't create user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created"})
}

func validateLoginInput(input *LoginInput) []ValidationError {
	var errors []ValidationError

	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var ve ValidationError
			ve.Field = err.Field()
			switch err.Tag() {
			case "required":
				ve.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				ve.Message = fmt.Sprintf("%s must be a valid email", err.Field())
			case "min":
				ve.Message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			}
			errors = append(errors, ve)
		}
	}

	return errors
}

func LoginUser(c echo.Context) error {
	var input LoginInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid input"})
	}
	validationErrors := validateLoginInput(&input)
	if len(validationErrors) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": validationErrors})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "invalid email or password"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "login successful", "data": echo.Map{
		"email": input.Email,
		"token": "generated_token",
	}})
}
