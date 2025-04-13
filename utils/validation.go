package utils

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var Validation = validator.New()

func FormatValidationErrors(ve validator.ValidationErrors) []ValidationErrorResponse {
	var errors []ValidationErrorResponse

	for _, fe := range ve {
		var message string

		switch fe.Tag() {
		case "required":
			message = fe.Field() + " is required"
		case "email":
			message = "Invalid email address"
		case "min":
			message = fe.Field() + " must be at least " + fe.Param() + " characters long"
		default:
			message = fe.Field() + " is invalid"
		}

		errors = append(errors, ValidationErrorResponse{
			Field:   fe.Field(),
			Message: message,
		})
	}

	return errors
}

func GenerateHashedPassword(p string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
