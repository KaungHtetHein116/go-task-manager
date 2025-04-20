package transport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewApiErrorResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func NewApiSuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, echo.Map{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func NewApiCreateSuccessResponse(c echo.Context, message string, data interface{}) error {
	formattedMessage := message
	if message == "" {
		formattedMessage = http.StatusText(http.StatusCreated)
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"code":    http.StatusCreated,
		"message": formattedMessage,
		"data":    data,
	})
}
