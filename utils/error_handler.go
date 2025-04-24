package utils

import (
	"errors"
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// Handle validation errors
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		validationErrors, _ := FormatValidationErrors(err)

		transport.NewApiErrorResponse(c,
			http.StatusBadRequest, ErrValidationFailed,
			validationErrors)

		return
	}

	// Handle GORM record not found errors
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, ErrRecordNotFound) {
		transport.NewApiErrorResponse(c,
			http.StatusNotFound, ErrErrorRecordNotFound,
			nil)
		return
	}

	// Handle HTTP errors
	if he, ok := err.(*echo.HTTPError); ok {
		transport.NewApiErrorResponse(c,
			he.Code, he.Message.(string),
			nil)
		return
	}

	// Handle other errors

	transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
}
