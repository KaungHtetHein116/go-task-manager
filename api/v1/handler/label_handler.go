package handler

import (
	"net/http"
	"strconv"

	"github.com/KaungHtetHein116/personal-task-manager/api/transport"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/request"
	"github.com/KaungHtetHein116/personal-task-manager/api/v1/response"
	"github.com/KaungHtetHein116/personal-task-manager/internal/entity"
	"github.com/KaungHtetHein116/personal-task-manager/internal/usecase"
	"github.com/KaungHtetHein116/personal-task-manager/pkg/constants"
	"github.com/KaungHtetHein116/personal-task-manager/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LabelHandler struct {
	usecase usecase.LabelUsecase
}

func NewLabelHandler(u usecase.LabelUsecase) *LabelHandler {
	return &LabelHandler{
		usecase: u,
	}
}

func (h *LabelHandler) CreateLabel(c echo.Context, input *request.CreateLabelInput) error {
	userID := c.Get("user_id").(uint)
	label := &entity.Label{
		Name:   input.Name,
		UserID: userID,
	}

	err := h.usecase.CreateLabel(label)
	if err != nil {
		if err == utils.ErrDuplicateEntry {
			return transport.NewApiErrorResponse(c, http.StatusConflict, constants.ErrDuplicatedData, nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiCreateSuccessResponse(c, constants.MSG_SUCCESS, response.NewLabelResponse(label))
}

func (h *LabelHandler) GetLabels(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	labels, err := h.usecase.GetLabels(userID)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.MSG_SUCCESS, response.NewLabelResponses(labels))
}

func (h *LabelHandler) GetLabelByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, constants.ErrInvalidProjectID, nil)
	}

	label, err := h.usecase.GetLabelByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrErrorRecordNotFound, nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.Successful, response.NewLabelResponse(label))
}

func (h *LabelHandler) UpdateLabel(c echo.Context, input *request.UpdateLabelInput) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, constants.ErrInvalidProjectID, nil)
	}

	label := &entity.Label{
		Name: input.Name,
	}
	label.ID = uint(id)

	err = h.usecase.UpdateLabel(label)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrErrorRecordNotFound, nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.MSG_SUCCESS, response.NewLabelResponse(label))
}

func (h *LabelHandler) DeleteLabel(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return transport.NewApiErrorResponse(c, http.StatusBadRequest, constants.ErrInvalidProjectID, nil)
	}

	err = h.usecase.DeleteLabel(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return transport.NewApiErrorResponse(c, http.StatusNotFound, constants.ErrErrorRecordNotFound, nil)
		}
		return transport.NewApiErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return transport.NewApiSuccessResponse(c, http.StatusOK, constants.MSG_SUCCESS, nil)
}
