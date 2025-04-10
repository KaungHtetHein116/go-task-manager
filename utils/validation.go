package utils

import "github.com/go-playground/validator/v10"

type ValidationErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

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
