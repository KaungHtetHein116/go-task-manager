package handler

import (
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: usecase}
}

func (h *UserHandler) Register(c echo.Context, input *request.RegisterUserInput) error {
	err := h.userUsecase.Register(input.Name, input.Email, input.Password)
	if err != nil {
		if err.Error() == "email already exists" {
			return transport.NewApiErrorResponse(c, http.StatusConflict, err.Error(), nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "cannot create the user", nil)
	}

	return transport.NewApiCreateSuccessResponse(c, "user created successfully", nil)
}

func (h *UserHandler) Login(c echo.Context, input *request.LoginUserInput) error {
	token, user, err := h.userUsecase.Login(input.Email, input.Password)
	if err != nil {
		if err.Error() == "invalid credentials" || err.Error() == "invalid password" {
			return transport.NewApiErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, "unknown error", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "login successful", echo.Map{
		"name":  user.Username,
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.userUsecase.GetProfile(userID)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, "user not found", nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, "successful", user)
}
