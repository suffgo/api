package infrastructure

import (
	"net/http"
	"strconv"
	u "suffgo/internal/option/application/useCases"

	d "suffgo/internal/option/domain"
	v "suffgo/internal/option/domain/valueObjects"

	"github.com/labstack/echo/v4"

	sv "suffgo/internal/shared/domain/valueObjects"

	e "suffgo/internal/option/domain/errors"
	se "suffgo/internal/shared/domain/errors"
)

type OptionEchoHandler struct {
	CreateOptionUsecase     *u.CreateUsecase
	DeleteOptionUsecase     *u.DeleteUsecase
	GetAllOptionUsecase     *u.GetAllUsecase
	GetOtionByIDUsecase     *u.GetByIDUsecase
	GetOptionByValueUsecase *u.GetByValueUsecase
}

func NewOptionEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	getByValue *u.GetByValueUsecase,
) *OptionEchoHandler {
	return &OptionEchoHandler{
		CreateOptionUsecase:     createUC,
		DeleteOptionUsecase:     deleteUC,
		GetAllOptionUsecase:     getAllUC,
		GetOtionByIDUsecase:     getByIDUC,
		GetOptionByValueUsecase: getByValue,
	}
}

func (h *OptionEchoHandler) CreateOption(c echo.Context) error {
	var req d.OptionCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	value, err := v.NewValue(req.Value)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	option := d.NewOption(
		nil,
		*value,
	)

	err = h.CreateOptionUsecase.Execute(*option)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *OptionEchoHandler) DeleteOption(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteOptionUsecase.Execute(*id)

	if err != nil {
		if err.Error() == "option not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"succes": "Option deleted succesfully"})
}

func (h *OptionEchoHandler) GetAllOptions(c echo.Context) error {
	options, err := h.GetAllOptionUsecase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var OptionsDTO []d.OptionDTO
	for _, option := range options {
		OptionDTO := &d.OptionDTO{
			ID:    option.ID().Id,
			Value: option.Value().Value,
		}
		OptionsDTO = append(OptionsDTO, *OptionDTO)
	}

	return c.JSON(http.StatusOK, OptionsDTO)
}

func (h *OptionEchoHandler) GetOptionByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	option, err := h.GetOtionByIDUsecase.Execute(*id)

	if err != nil {
		if err.Error() == "option not found" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	optionDTO := &d.OptionDTO{
		ID:    option.ID().Id,
		Value: option.Value().Value,
	}

	return c.JSON(http.StatusOK, optionDTO)
}

func (h *OptionEchoHandler) GetOptionByValue(c echo.Context) error {
	valueParam := c.Param("value")

	value, err := v.NewValue(valueParam)
	if err != nil {
		invalidErr := &e.InvalidValueError{Value: valueParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}
	option, err := h.GetOptionByValueUsecase.Execute(*value)

	if err != nil {
		if err.Error() == "option not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	optionDTO := &d.OptionDTO{
		ID:    option.ID().Id,
		Value: option.Value().Value,
	}
	return c.JSON(http.StatusOK, optionDTO)

}
