package handler

import (
	"errors"
	"net/http"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	"github.com/KaungHtetHein116/personal-task-manager/pkg/constants"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: usecase}
}

func (h *UserHandler) Register(c echo.Context, input *request.RegisterUserInput) error {
	err := h.userUsecase.Register(input)
	if err != nil {
		// Check for specific business logic errors
		switch {
		case errors.Is(err, utils.ErrDuplicateEntry):
			return transport.NewApiErrorResponse(c, http.StatusConflict, constants.ErrEmailAlreadyRegistered, nil)
		case errors.Is(err, utils.ErrInvalidData):
			return transport.NewApiErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			// Let the CustomHTTPErrorHandler handle unexpected errors
			return err
		}
	}

	return transport.NewApiCreateSuccessResponse(c, constants.SuccUserRegistered, nil)
}

func (h *UserHandler) Login(c echo.Context, input *request.LoginUserInput) error {
	token, user, err := h.userUsecase.Login(input)

	if err != nil {

		if errors.Is(err, utils.ErrInvalidData) {
			return transport.NewApiErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		}
		return err
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.SuccLoginSuccessful, echo.Map{
		"name":  user.Username,
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	user, err := h.userUsecase.GetProfile(userID)

	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrUserNotFound, nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.Successful, user)
}
